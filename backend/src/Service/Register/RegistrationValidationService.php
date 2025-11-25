<?php

declare(strict_types=1);

namespace App\Service\Register;

use App\Service\Login\LoginValidationService;
use App\Repository\UserRepository;
use Symfony\Component\HttpFoundation\JsonResponse;

/**
 * Registration Validation Service
 * Centralizes all registration validation logic
 */
class RegistrationValidationService
{
    public function __construct(
        private UserRepository $userRepository,
        private LoginValidationService $loginValidationService
    ) {}

    /**
     * Validate required fields for email registration
     *
     * @param array $data Request data
     * @return JsonResponse|null Returns error response if invalid, null if ok
     */
    public function validateRequiredFields(array $data): ?JsonResponse
    {
        $required = ['email', 'password', 'firstName', 'lastName'];
        $missing = [];

        foreach ($required as $field) {
            if (!isset($data[$field]) || trim($data[$field]) === '') {
                $missing[] = $field;
            }
        }

        if (!empty($missing)) {
            return new JsonResponse([
                'error' => 'Missing required fields: ' . implode(', ', $missing)
            ], 400);
        }

        return null;
    }

    /**
     * Validate email format
     *
     * @param string $email Email address
     * @return JsonResponse|null Returns error response if invalid, null if ok
     */
    public function validateEmailFormat(string $email): ?JsonResponse
    {
        if (!filter_var($email, FILTER_VALIDATE_EMAIL)) {
            return new JsonResponse([
                'error' => 'Invalid email format'
            ], 400);
        }

        return null;
    }

    /**
     * Validate password strength
     *
     * @param string $password Password
     * @return JsonResponse|null Returns error response if weak, null if ok
     */
    public function validatePasswordStrength(string $password): ?JsonResponse
    {
        if (strlen($password) < 8) {
            return new JsonResponse([
                'error' => 'Password must be at least 8 characters long'
            ], 400);
        }

        // Optional: Add more password complexity rules here
        // - At least one uppercase letter
        // - At least one number
        // - At least one special character

        return null;
    }

    /**
     * Check if user already exists
     *
     * @param string $email Email address
     * @return JsonResponse|null Returns error response if exists, null if ok
     */
    public function validateUserDoesNotExist(string $email): ?JsonResponse
    {
        $existingUser = $this->loginValidationService->findUserByEmail($email);

        if ($existingUser) {
            return new JsonResponse([
                'error' => 'User with this email already exists'
            ], 409);
        }

        return null;
    }

    /**
     * Validate IP rate limiting (reuse from login)
     *
     * @param string $ipAddress Client IP address
     * @return JsonResponse|null Returns error response if rate limited, null if ok
     */
    public function validateIpRateLimit(string $ipAddress): ?JsonResponse
    {
        return $this->loginValidationService->validateIpRateLimit($ipAddress);
    }

    /**
     * Validate all registration requirements for email registration
     *
     * @param array $data Request data
     * @param string $ipAddress Client IP address
     * @return JsonResponse|null Returns error response if validation fails, null if ok
     */
    public function validateEmailRegistration(array $data, string $ipAddress): ?JsonResponse
    {
        // Check IP rate limiting
        $ipError = $this->validateIpRateLimit($ipAddress);
        if ($ipError) {
            return $ipError;
        }

        // Check required fields
        $fieldsError = $this->validateRequiredFields($data);
        if ($fieldsError) {
            return $fieldsError;
        }

        // Validate email format
        $emailError = $this->validateEmailFormat($data['email']);
        if ($emailError) {
            return $emailError;
        }

        // Validate password strength
        $passwordError = $this->validatePasswordStrength($data['password']);
        if ($passwordError) {
            return $passwordError;
        }

        // Check if user already exists
        $existsError = $this->validateUserDoesNotExist($data['email']);
        if ($existsError) {
            return $existsError;
        }

        return null;
    }

    /**
     * Validate OAuth registration
     *
     * @param array $oauthUserData OAuth user data
     * @param string $ipAddress Client IP address
     * @return JsonResponse|null Returns error response if validation fails, null if ok
     */
    public function validateOAuthRegistration(array $oauthUserData, string $ipAddress): ?JsonResponse
    {
        // Check IP rate limiting
        $ipError = $this->validateIpRateLimit($ipAddress);
        if ($ipError) {
            return $ipError;
        }

        // Validate email is verified (reuse from login validation)
        $emailError = $this->loginValidationService->validateEmailVerified($oauthUserData);
        if ($emailError) {
            return $emailError;
        }

        // Check if user already exists
        $existsError = $this->validateUserDoesNotExist($oauthUserData['email']);
        if ($existsError) {
            return $existsError;
        }

        return null;
    }
}
