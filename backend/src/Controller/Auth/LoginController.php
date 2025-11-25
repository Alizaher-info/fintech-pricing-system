<?php

declare(strict_types=1);

namespace App\Controller\Auth;

use App\Security\JWTAuthenticator;
use App\Service\Login\LoginStrategyFactory;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;
use Exception;

/**
 * Login Controller
 * Routes login requests to appropriate strategy via LoginStrategyFactory
 * Clean, simple controller that delegates all logic to services
 */
class LoginController extends AbstractController
{
    public function __construct(
        private JWTAuthenticator $authenticator,
        private LoginStrategyFactory $loginStrategyFactory
    ) {}

    /**
     * Login endpoint - routes to strategy based on path parameter
     *
     * @param string $strategy Login strategy (email, google, microsoft, etc.)
     */
    #[Route('/api/login/{strategy}', name: 'api_login', methods: ['POST'])]
    public function login(Request $request, string $strategy): JsonResponse
    {
        try {
            $data = json_decode($request->getContent(), true);

            // Delegate to LoginStrategyFactory
            return $this->loginStrategyFactory->login($strategy, $request, $data);

        } catch (Exception $e) {
            return new JsonResponse([
                'error' => 'Login failed: ' . $e->getMessage()
            ], 500);
        }
    }

    /**
     * Validate JWT token
     */
    #[Route('/api/validate', name: 'api_validate', methods: ['GET'])]
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
