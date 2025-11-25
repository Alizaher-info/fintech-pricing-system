<?php

declare(strict_types=1);

namespace App\Service\Login;

use App\Entity\User;
use App\Event\UserLoginFailedEvent;
use App\Event\UserLoginSuccessEvent;
use App\Interface\LoginInterface;
use App\Security\JWTAuthenticator;
use App\Service\GoogleOAuthService;
use App\Service\Login\LoginValidationService;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\EventDispatcher\EventDispatcherInterface;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\PasswordHasher\Hasher\UserPasswordHasherInterface;

/**
 * Google OAuth Login Strategy
 * Handles login with Google OAuth
 */
class GoogleLoginStrategy implements LoginInterface
{
    public function __construct(
        private LoginValidationService $validationService,
        private GoogleOAuthService $googleOAuthService,
        private JWTAuthenticator $authenticator,
        private EventDispatcherInterface $eventDispatcher,
        private UserPasswordHasherInterface $passwordHasher,
        private EntityManagerInterface $entityManager
    ) {}

    /**
     * Login with Google OAuth
     */
    public function login(Request $request, array $data): JsonResponse
    {
        // Validate required fields
        if (!isset($data['credential'])) {
            return new JsonResponse([
                'error' => 'Missing Google credential (ID token)'
            ], 400);
        }

        $credential = $data['credential'];
        $ipAddress = $request->getClientIp() ?? 'unknown';

        try {
            // Verify Google credential and get user info
            $googleUser = $this->googleOAuthService->verifyIdToken($credential);

            // Validate OAuth login (rate limiting, email verified)
            $validation = $this->validationService->validateOAuthLogin($googleUser, $ipAddress);

            if ($validation['error']) {
                $this->handleLoginFailure(
                    $googleUser['email'] ?? 'unknown',
                    'OAuth validation failed',
                    $ipAddress,
                    $request->headers->get('User-Agent'),
                    null
                );
                return $validation['error'];
            }

            $user = $validation['user'];

            // Auto-create user if not exists
            if (!$user) {
                $user = $this->createUserFromGoogle($googleUser);
            }

            // Login success
            return $this->handleSuccessfulLogin(
                $user,
                $ipAddress,
                $request->headers->get('User-Agent')
            );

        } catch (\Exception $e) {
            $this->handleLoginFailure(
                $googleUser['email'] ?? 'unknown',
                'Google OAuth verification failed: ' . $e->getMessage(),
                $ipAddress,
                $request->headers->get('User-Agent'),
                null
            );
            return new JsonResponse([
                'error' => 'Google login failed: ' . $e->getMessage()
            ], 401);
        }
    }

    /**
     * Create user from Google OAuth data
     */
    private function createUserFromGoogle(array $googleUser): User
    {
        $user = new User();
        $user->setEmail($googleUser['email']);

        // Parse name from Google (format: "First Last")
        $nameParts = explode(' ', $googleUser['name'] ?? '', 2);
        $user->setFirstName($nameParts[0] ?? 'Google');
        $user->setLastName($nameParts[1] ?? 'User');

        // Generate random password (user won't use it, Google OAuth only)
        $randomPassword = bin2hex(random_bytes(32));
        $hashedPassword = $this->passwordHasher->hashPassword($user, $randomPassword);
        $user->setPassword($hashedPassword);

        $user->setRoles(['ROLE_USER']);
        $user->setIsActive(true);
        $user->setCreatedAt(new \DateTimeImmutable());
        $user->setUpdatedAt(new \DateTimeImmutable());

        // Save the new user
        $this->entityManager->persist($user);
        $this->entityManager->flush();

        error_log("✅ Auto-created Google OAuth user: {$googleUser['email']}");

        return $user;
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
