<?php

declare(strict_types=1);

namespace App\Service\Login;

use App\Entity\User;
use App\Repository\UserRepository;
use App\Service\RateLimiterService;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\PasswordHasher\Hasher\UserPasswordHasherInterface;

/**
 * Login Validation Service
 * Centralizes all login validation logic (rate limiting, user checks, password verification)
 */
class LoginValidationService
{
    public function __construct(
        private RateLimiterService $rateLimiterService,
        private UserRepository $userRepository,
        private UserPasswordHasherInterface $passwordHasher
    ) {}

    /**
     * Validate IP rate limiting
     *
     * @param string $ipAddress Client IP address
     * @return JsonResponse|null Returns error response if rate limited, null if ok
     */
    public function validateIpRateLimit(string $ipAddress): ?JsonResponse
    {
        if ($this->rateLimiterService->isIpRateLimited($ipAddress)) {
            $lockoutInfo = $this->rateLimiterService->getLockoutTimeRemaining($ipAddress);

            return new JsonResponse([
                'error' => 'Too many failed attempts. IP address is temporarily blocked.',
                'lockout_remaining' => $lockoutInfo['ip_lockout_remaining'],
                'retry_after' => $lockoutInfo['ip_lockout_remaining']
            ], 429);
        }

        return null;
    }

    /**
     * Validate user-specific rate limiting
     *
     * @param int $userId User ID
     * @param string $ipAddress Client IP address
     * @return JsonResponse|null Returns error response if locked, null if ok
     */
    public function validateUserRateLimit(int $userId, string $ipAddress): ?JsonResponse
    {
        if ($this->rateLimiterService->isUserLocked($userId)) {
            $lockoutInfo = $this->rateLimiterService->getLockoutTimeRemaining($ipAddress, $userId);

            return new JsonResponse([
                'error' => 'Account temporarily locked due to multiple failed login attempts.',
                'lockout_remaining' => $lockoutInfo['user_lockout_remaining'],
                'retry_after' => $lockoutInfo['user_lockout_remaining']
            ], 423);
        }

        return null;
    }

    /**
     * Find user by email
     *
     * @param string $email User email
     * @return User|null Returns user if found, null otherwise
     */
    public function findUserByEmail(string $email): ?User
    {
        return $this->userRepository->findByEmail($email);
    }

    /**
     * Validate user exists
     *
     * @param string $email User email
     * @return array{user: User|null, error: JsonResponse|null}
     */
    public function validateUserExists(string $email): array
    {
        $user = $this->findUserByEmail($email);

        if (!$user) {
            return [
                'user' => null,
                'error' => new JsonResponse(['error' => 'Invalid credentials'], 401)
            ];
        }

        return ['user' => $user, 'error' => null];
    }

    /**
     * Validate user is active
     *
     * @param User $user User entity
     * @return JsonResponse|null Returns error response if inactive, null if ok
     */
    public function validateUserActive(User $user): ?JsonResponse
    {
        if (!$user->isActive()) {
            return new JsonResponse(['error' => 'Account is deactivated'], 401);
        }

        return null;
    }

    /**
     * Validate password
     *
     * @param User $user User entity
     * @param string $password Plain text password
     * @param string $ipAddress Client IP address
     * @return array{valid: bool, error: JsonResponse|null}
     */
    public function validatePassword(User $user, string $password, string $ipAddress): array
    {
        if (!$this->passwordHasher->isPasswordValid($user, $password)) {
            $remainingInfo = $this->rateLimiterService->getRemainingAttempts($ipAddress, $user->getId());

            return [
                'valid' => false,
                'error' => new JsonResponse([
                    'error' => 'Invalid credentials',
                    'remaining_attempts' => [
                        'ip' => $remainingInfo['ip_remaining'],
                        'user' => $remainingInfo['user_remaining']
                    ]
                ], 401)
            ];
        }

        return ['valid' => true, 'error' => null];
    }

    /**
     * Validate email is verified (for OAuth)
     *
     * @param array $oauthUserData OAuth user data
     * @return JsonResponse|null Returns error response if not verified, null if ok
     */
    public function validateEmailVerified(array $oauthUserData): ?JsonResponse
    {
        if (!$oauthUserData['email'] || !$oauthUserData['email_verified']) {
            return new JsonResponse([
                'error' => 'Email not verified with OAuth provider'
            ], 401);
        }

        return null;
    }

    /**
     * Validate all login requirements for email login
     *
     * @param string $email User email
     * @param string $password Plain text password
     * @param string $ipAddress Client IP address
     * @return array{user: User|null, error: JsonResponse|null}
     */
    public function validateEmailLogin(string $email, string $password, string $ipAddress): array
    {
        // Check IP rate limiting
        $ipError = $this->validateIpRateLimit($ipAddress);
        if ($ipError) {
            return ['user' => null, 'error' => $ipError];
        }

        // Find user
        $result = $this->validateUserExists($email);
        if ($result['error']) {
            return $result;
        }

        $user = $result['user'];

        // Check user-specific rate limiting
        $userError = $this->validateUserRateLimit($user->getId(), $ipAddress);
        if ($userError) {
            return ['user' => null, 'error' => $userError];
        }

        // Check if user is active
        $activeError = $this->validateUserActive($user);
        if ($activeError) {
            return ['user' => null, 'error' => $activeError];
        }

        // Verify password
        $passwordResult = $this->validatePassword($user, $password, $ipAddress);
        if (!$passwordResult['valid']) {
            return ['user' => null, 'error' => $passwordResult['error']];
        }

        return ['user' => $user, 'error' => null];
    }

    /**
     * Validate all login requirements for OAuth login
     *
     * @param array $oauthUserData OAuth user data with email and email_verified
     * @param string $ipAddress Client IP address
     * @return array{user: User|null, error: JsonResponse|null}
     */
    public function validateOAuthLogin(array $oauthUserData, string $ipAddress): array
    {
        // Check IP rate limiting
        $ipError = $this->validateIpRateLimit($ipAddress);
        if ($ipError) {
            return ['user' => null, 'error' => $ipError];
        }

        // Validate email is verified
        $emailError = $this->validateEmailVerified($oauthUserData);
        if ($emailError) {
            return ['user' => null, 'error' => $emailError];
        }

        // Find user (may be null for auto-registration)
        $user = $this->findUserByEmail($oauthUserData['email']);

        if ($user) {
            // Check if user is active
            $activeError = $this->validateUserActive($user);
            if ($activeError) {
                return ['user' => null, 'error' => $activeError];
            }
        }

        return ['user' => $user, 'error' => null];
    }
}
