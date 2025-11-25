<?php

declare(strict_types=1);

namespace App\Service\Register;

use App\Entity\User;
use App\Interface\RegisterInterface;
use App\Repository\UserRepository;
use App\Security\JWTAuthenticator;
use App\Service\GoogleOAuthService;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\PasswordHasher\Hasher\UserPasswordHasherInterface;

/**
 * Google OAuth Registration Strategy
 * Handles user registration with Google OAuth
 * Note: Google registration is same as login - auto-creates user if doesn't exist
 */
class GoogleRegisterStrategy implements RegisterInterface
{
    public function __construct(
        private GoogleOAuthService $googleOAuthService,
        private RegistrationValidationService $validationService,
        private UserRepository $userRepository,
        private UserPasswordHasherInterface $passwordHasher,
        private JWTAuthenticator $authenticator,
        private EntityManagerInterface $entityManager
    ) {}

    public function register(Request $request, array $data): JsonResponse
    {
        $ipAddress = $request->getClientIp() ?? '127.0.0.1';

        if (!isset($data['credential'])) {
            return new JsonResponse(['error' => 'Missing Google credential'], 400);
        }

        // Verify Google ID token
        $googleUser = $this->googleOAuthService->verifyIdToken($data['credential']);
        if (!$googleUser) {
            return new JsonResponse(['error' => 'Invalid Google credential'], 401);
        }

        // Validate OAuth registration
        $validationError = $this->validationService->validateOAuthRegistration($googleUser, $ipAddress);
        if ($validationError) {
            return $validationError;
        }

        // Parse name from Google
        $nameParts = explode(' ', $googleUser['name'] ?? '');
        $firstName = $nameParts[0] ?? 'Google';
        $lastName = isset($nameParts[1]) ? implode(' ', array_slice($nameParts, 1)) : 'User';

        // Create new user with Google data
        $user = new User();
        $user->setEmail($googleUser['email']);
        $user->setFirstName($firstName);
        $user->setLastName($lastName);
        $user->setRoles(['ROLE_USER']);
        $user->setIsActive(true);
        $user->setCreatedAt(new \DateTimeImmutable());
        $user->setUpdatedAt(new \DateTimeImmutable());

        // Set random password (user won't use it - OAuth only)
        $randomPassword = bin2hex(random_bytes(32));
        $hashedPassword = $this->passwordHasher->hashPassword($user, $randomPassword);
        $user->setPassword($hashedPassword);

        // Save to database
        $this->entityManager->persist($user);
        $this->entityManager->flush();

        // Generate JWT token
        $token = $this->authenticator->generateToken(
            $user->getId(),
            $user->getEmail(),
            'ROLE_USER'
        );

        return new JsonResponse([
            'success' => true,
            'message' => 'User registered successfully with Google',
            'token' => $token,
            'user' => [
                'id' => $user->getId(),
                'email' => $user->getEmail(),
                'firstName' => $user->getFirstName(),
                'lastName' => $user->getLastName(),
                'roles' => $user->getRoles()
            ]
        ], 201);
    }
}
