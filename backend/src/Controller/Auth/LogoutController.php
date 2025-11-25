<?php

declare(strict_types=1);

namespace App\Controller\Auth;

use App\Security\JWTAuthenticator;
use App\Service\Auth\TokenBlacklistService;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;
use Exception;

class LogoutController extends AbstractController
{
    private JWTAuthenticator $authenticator;
    private TokenBlacklistService $blacklistService;

    public function __construct(
        JWTAuthenticator $authenticator,
        TokenBlacklistService $blacklistService
    ) {
        $this->authenticator = $authenticator;
        $this->blacklistService = $blacklistService;
    }

    #[Route('/api/logout', name: 'api_logout', methods: ['POST'])]
    public function logout(Request $request): JsonResponse
    {
        try {
            // Validate token and get user
            $user = $this->authenticator->requireAuth($request);

            // Extract token from request
            $token = $this->authenticator->extractToken($request);

            if ($token) {
                // Blacklist this specific token
                $this->blacklistService->blacklistToken($token);
            }

            return new JsonResponse([
                'success' => true,
                'message' => 'Logged out successfully'
            ]);

        } catch (Exception $e) {
            return new JsonResponse([
                'error' => 'Logout failed: ' . $e->getMessage()
            ], 500);
        }
    }

    #[Route('/api/logout/all', name: 'api_logout_all', methods: ['POST'])]
    public function logoutAll(Request $request): JsonResponse
    {
        try {
            // Validate token and get user
            $user = $this->authenticator->requireAuth($request);

            // Blacklist ALL tokens for this user
            // Any token issued before this moment will be invalid
            $this->blacklistService->blacklistAllUserTokens($user['user_id']);

            return new JsonResponse([
                'success' => true,
                'message' => 'Logged out from all devices successfully'
            ]);

        } catch (Exception $e) {
            return new JsonResponse([
                'error' => 'Logout all failed: ' . $e->getMessage()
            ], 500);
        }
    }
}
