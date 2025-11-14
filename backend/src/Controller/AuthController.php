<?php

declare(strict_types=1);

namespace App\Controller;

use App\Entity\User;
use App\Event\UserLoginSuccessEvent;
use App\Event\UserLoginFailedEvent;
use App\Repository\UserRepository;
use App\Security\JWTAuthenticator;
use App\Service\RateLimiterService;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\EventDispatcher\EventDispatcherInterface;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;
use Symfony\Component\PasswordHasher\Hasher\UserPasswordHasherInterface;
use Exception;

class AuthController extends AbstractController
{
    private JWTAuthenticator $authenticator;
    private UserRepository $userRepository;
    private EventDispatcherInterface $eventDispatcher;
    private UserPasswordHasherInterface $passwordHasher;
    private RateLimiterService $rateLimiterService;

    public function __construct(
        JWTAuthenticator $authenticator,
        UserRepository $userRepository,
        EventDispatcherInterface $eventDispatcher,
        UserPasswordHasherInterface $passwordHasher,
        RateLimiterService $rateLimiterService
    ) {
        $this->authenticator = $authenticator;
        $this->userRepository = $userRepository;
        $this->eventDispatcher = $eventDispatcher;
        $this->passwordHasher = $passwordHasher;
        $this->rateLimiterService = $rateLimiterService;
    }

    #[Route('/api/auth/login', name: 'api_auth_login', methods: ['POST'])]
    public function login(Request $request): JsonResponse
    {
        try {
            $data = json_decode($request->getContent(), true);

            if (!$data || !isset($data['email']) || !isset($data['password'])) {
                return new JsonResponse([
                    'error' => 'Missing email or password'
                ], 400);
            }

            $email = $data['email'];
            $password = $data['password'];
            $ipAddress = $request->getClientIp() ?? 'unknown';

            // Check IP rate limiting FIRST (fastest check)
            if ($this->rateLimiterService->isIpRateLimited($ipAddress)) {
                $lockoutInfo = $this->rateLimiterService->getLockoutTimeRemaining($ipAddress);
                
                return new JsonResponse([
                    'error' => 'Too many failed attempts. IP address is temporarily blocked.',
                    'lockout_remaining' => $lockoutInfo['ip_lockout_remaining'],
                    'retry_after' => $lockoutInfo['ip_lockout_remaining']
                ], 429);
            }

            // Find user in database
            $user = $this->userRepository->findByEmail($email);

            if (!$user) {
                // Dispatch login failed event - user not found
                $failedEvent = new UserLoginFailedEvent(
                    $email,
                    'User not found',
                    $ipAddress,
                    $request->headers->get('User-Agent'),
                    null // no user object since user doesn't exist
                );
                $this->eventDispatcher->dispatch($failedEvent);
                
                return new JsonResponse([
                    'error' => 'Invalid credentials'
                ], 401);
            }

            // Check user-specific rate limiting
            if ($this->rateLimiterService->isUserLocked($user->getId())) {
                $lockoutInfo = $this->rateLimiterService->getLockoutTimeRemaining($ipAddress, $user->getId());
                
                return new JsonResponse([
                    'error' => 'Account temporarily locked due to multiple failed login attempts.',
                    'lockout_remaining' => $lockoutInfo['user_lockout_remaining'],
                    'retry_after' => $lockoutInfo['user_lockout_remaining']
                ], 423); // 423 Locked
            }

            // Check if user is active
            if (!$user->isActive()) {
                return new JsonResponse([
                    'error' => 'Account is deactivated'
                ], 401);
            }

            // Verify password
            if (!$this->passwordHasher->isPasswordValid($user, $password)) {
                // Dispatch login failed event (will trigger rate limiting)
                $failedEvent = new UserLoginFailedEvent(
                    $email,
                    'Invalid password',
                    $ipAddress,
                    $request->headers->get('User-Agent'),
                    $user
                );
                $this->eventDispatcher->dispatch($failedEvent);
                
                // Get remaining attempts for response
                $remainingInfo = $this->rateLimiterService->getRemainingAttempts($ipAddress, $user->getId());
                
                $response = [
                    'error' => 'Invalid credentials',
                    'remaining_attempts' => [
                        'ip' => $remainingInfo['ip_remaining'],
                        'user' => $remainingInfo['user_remaining']
                    ]
                ];
                
                return new JsonResponse($response, 401);
            }

            // Dispatch login success event (will clear rate limiting)
            $successEvent = new UserLoginSuccessEvent(
                $user,
                $ipAddress,
                $request->headers->get('User-Agent')
            );
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

        } catch (Exception $e) {
            return new JsonResponse([
                'error' => 'Login failed: ' . $e->getMessage()
            ], 500);
        }
    }

    #[Route('/api/auth/validate', name: 'api_auth_validate', methods: ['GET'])]
    public function validate(Request $request): JsonResponse
    {
        try {
            $user = $this->authenticator->requireAuth($request);

            return new JsonResponse([
                'success' => true,
                'user' => $user,
                'message' => 'Token is valid'
            ]);

        } catch (Exception $e) {
            return new JsonResponse([
                'success' => false,
                'error' => $e->getMessage()
            ], 401);
        }
    }

    #[Route('/api/auth/register', name: 'api_auth_register', methods: ['POST'])]
    public function register(Request $request): JsonResponse
    {
        try {
            $data = json_decode($request->getContent(), true);

            if (!$data || !isset($data['email']) || !isset($data['password']) ||
                !isset($data['firstName']) || !isset($data['lastName'])) {
                return new JsonResponse([
                    'error' => 'Missing required fields'
                ], 400);
            }

            // Check if user already exists
            $existingUser = $this->userRepository->findByEmail($data['email']);
            if ($existingUser) {
                return new JsonResponse([
                    'error' => 'User with this email already exists'
                ], 409);
            }

            // Create new user
            $user = new User();
            $user->setEmail($data['email']);
            $user->setFirstName($data['firstName']);
            $user->setLastName($data['lastName']);
            $user->setRoles(['ROLE_USER']);
            $user->setIsActive(true);
            $user->setCreatedAt(new \DateTimeImmutable());
            $user->setUpdatedAt(new \DateTimeImmutable());

            // Hash password
            $hashedPassword = $this->passwordHasher->hashPassword($user, $data['password']);
            $user->setPassword($hashedPassword);

            // Save to database
            $this->userRepository->save($user, true);

            // Generate JWT token
            $token = $this->authenticator->generateToken(
                $user->getId(),
                $user->getEmail(),
                'ROLE_USER'
            );

            return new JsonResponse([
                'success' => true,
                'message' => 'User registered successfully',
                'token' => $token,
                'user' => [
                    'id' => $user->getId(),
                    'email' => $user->getEmail(),
                    'firstName' => $user->getFirstName(),
                    'lastName' => $user->getLastName(),
                    'roles' => $user->getRoles()
                ]
            ], 201);

        } catch (Exception $e) {
            return new JsonResponse([
                'error' => 'Registration failed: ' . $e->getMessage()
            ], 500);
        }
    }
}
