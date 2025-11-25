<?php

declare(strict_types=1);

namespace App\Service\Login;

use App\Entity\User;
use App\Event\UserLoginFailedEvent;
use App\Event\UserLoginSuccessEvent;
use App\Interface\LoginInterface;
use App\Repository\UserRepository;
use App\Security\JWTAuthenticator;
use App\Service\Login\LoginValidationService;
use Symfony\Component\EventDispatcher\EventDispatcherInterface;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;

/**
 * Email Login Strategy
 * Handles login with email and password
 */
class EmailLoginStrategy implements LoginInterface
{
    public function __construct(
        private LoginValidationService $validationService,
        private JWTAuthenticator $authenticator,
        private EventDispatcherInterface $eventDispatcher
    ) {}

    /**
     * Login with email and password
     */
    public function login(Request $request, array $data): JsonResponse
    {
        // Validate required fields
        if (!isset($data['email']) || !isset($data['password'])) {
            return new JsonResponse([
                'error' => 'Missing email or password'
            ], 400);
        }

        $email = $data['email'];
        $password = $data['password'];
        $ipAddress = $request->getClientIp() ?? 'unknown';

        // Validate email login (rate limiting, user exists, active, password)
        $validation = $this->validationService->validateEmailLogin($email, $password, $ipAddress);

        if ($validation['error']) {
            // Dispatch failure event
            $this->handleLoginFailure(
                $email,
                'Login validation failed',
                $ipAddress,
                $request->headers->get('User-Agent'),
                $validation['user']
            );
            return $validation['error'];
        }

        // Login success
        return $this->handleSuccessfulLogin(
            $validation['user'],
            $ipAddress,
            $request->headers->get('User-Agent')
        );
    }

    /**
     * Handle login failure - dispatch failure event
     */
    private function handleLoginFailure(
        string $email,
        string $reason,
        string $ipAddress,
        ?string $userAgent,
        ?User $user
    ): void {
        $failedEvent = new UserLoginFailedEvent(
            $email,
            $reason,
            $ipAddress,
            $userAgent,
            $user
        );
        $this->eventDispatcher->dispatch($failedEvent);
    }

    /**
     * Handle successful login - dispatch event and generate JWT
     */
    private function handleSuccessfulLogin(User $user, string $ipAddress, ?string $userAgent): JsonResponse
    {
        // Dispatch success event
        $successEvent = new UserLoginSuccessEvent($user, $ipAddress, $userAgent);
        $this->eventDispatcher->dispatch($successEvent);

        // Generate JWT token
        $token = $this->authenticator->generateToken(
            $user->getId(),
            $user->getEmail(),
            implode(',', $user->getRoles())
        );

        return new JsonResponse([
            'success' => true,
            'token' => $token,
            'user' => [
                'id' => $user->getId(),
                'email' => $user->getEmail(),
                'firstName' => $user->getFirstName(),
                'lastName' => $user->getLastName(),
                'roles' => $user->getRoles()
            ]
        ]);
    }
}
