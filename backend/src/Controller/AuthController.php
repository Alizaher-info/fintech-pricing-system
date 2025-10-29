<?php

declare(strict_types=1);

namespace App\Controller;

use App\Security\JWTAuthenticator;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;
use Exception;

class AuthController extends AbstractController
{
    private JWTAuthenticator $authenticator;

    public function __construct(JWTAuthenticator $authenticator)
    {
        $this->authenticator = $authenticator;
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

            // TODO: In echter Anwendung - Database lookup und password verification
            // Für Demo verwenden wir hardcoded values
            if ($email === 'admin@example.com' && $password === 'admin123') {
                $token = $this->authenticator->generateToken(1, $email, 'admin');

                return new JsonResponse([
                    'success' => true,
                    'token' => $token,
                    'user' => [
                        'id' => 1,
                        'email' => $email,
                        'role' => 'admin'
                    ]
                ]);
            }

            if ($email === 'user@example.com' && $password === 'user123') {
                $token = $this->authenticator->generateToken(2, $email, 'user');

                return new JsonResponse([
                    'success' => true,
                    'token' => $token,
                    'user' => [
                        'id' => 2,
                        'email' => $email,
                        'role' => 'user'
                    ]
                ]);
            }

            return new JsonResponse([
                'error' => 'Invalid credentials'
            ], 401);

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
}
