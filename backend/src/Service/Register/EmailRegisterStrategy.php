<?php

declare(strict_types=1);

namespace App\Service\Register;

use App\Entity\User;
use App\Interface\RegisterInterface;
use App\Repository\UserRepository;
use App\Security\JWTAuthenticator;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\PasswordHasher\Hasher\UserPasswordHasherInterface;

/**
 * Email Registration Strategy
 * Handles user registration with email and password
 */
class EmailRegisterStrategy implements RegisterInterface
{
    public function __construct(
        private RegistrationValidationService $validationService,
        private UserRepository $userRepository,
        private UserPasswordHasherInterface $passwordHasher,
        private JWTAuthenticator $authenticator,
        private EntityManagerInterface $entityManager
    ) {}

    public function register(Request $request, array $data): JsonResponse
    {
        $ipAddress = $request->getClientIp() ?? '127.0.0.1';

        // Validate all registration requirements
        $validationError = $this->validationService->validateEmailRegistration($data, $ipAddress);
        if ($validationError) {
            return $validationError;
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
        $this->entityManager->persist($user);
        $this->entityManager->flush();

        // Generate JWT token (auto-login after registration)
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
    }
}
