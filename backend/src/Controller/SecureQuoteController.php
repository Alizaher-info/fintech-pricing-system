<?php

declare(strict_types=1);

namespace App\Controller;

use App\Client\PricingGrpcClient;
use App\Security\JWTAuthenticator;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;
use Exception;

class SecureQuoteController extends AbstractController
{
    private PricingGrpcClient $pricingClient;
    private JWTAuthenticator $authenticator;

    public function __construct(JWTAuthenticator $authenticator)
    {
        $this->authenticator = $authenticator;

        // Verwende Umgebungsvariable oder Standard
        $target = $_ENV['PRICING_GRPC_TARGET'] ?? 'pricing-api:50051';
        $this->pricingClient = new PricingGrpcClient($target);
    }

    #[Route('/api/secure/quote', name: 'api_secure_quote', methods: ['POST'])]
    public function secureQuote(Request $request): JsonResponse
    {
        try {
            // 1. Authentication Required
            $user = $this->authenticator->requireAuth($request);

            // 2. Parse Request
            $data = json_decode($request->getContent(), true);

            if (!$data || !isset($data['amount']) || !isset($data['termMonths'])) {
                return new JsonResponse([
                    'error' => 'Missing required fields: amount, termMonths'
                ], 400);
            }

            $amount = (float) $data['amount'];
            $termMonths = (int) $data['termMonths'];

            // 3. Extract JWT token from request header
            $authHeader = $request->headers->get('Authorization');
            $jwtToken = null;
            if ($authHeader && str_starts_with($authHeader, 'Bearer ')) {
                $jwtToken = substr($authHeader, 7); // Remove "Bearer " prefix
            }

            // 4. Add user context and JWT token to gRPC call
            $options = [
                'userId' => $user['user_id'],  // ← User context for pricing
                'userRole' => $user['role'],
                'jwt_token' => $jwtToken,      // ← Forward JWT token to gRPC
            ];

            if (isset($data['riskScore'])) {
                $options['riskScore'] = (int) $data['riskScore'];
            }

            // 5. Call gRPC with user context and authentication
            $result = $this->pricingClient->quote($amount, $termMonths, $options);

            // 5. Return enriched response
            return new JsonResponse([
                'success' => true,
                'data' => $result,
                'user' => [
                    'id' => $user['user_id'],
                    'email' => $user['email'],
                    'role' => $user['role']
                ],
                'request' => [
                    'amount' => $amount,
                    'termMonths' => $termMonths,
                    'riskScore' => $options['riskScore'] ?? null
                ]
            ]);

        } catch (Exception $e) {
            return new JsonResponse([
                'success' => false,
                'error' => $e->getMessage(),
                'code' => $e->getCode()
            ], $e->getCode() === 401 ? 401 : 500);
        }
    }

    #[Route('/api/admin/quote', name: 'api_admin_quote', methods: ['POST'])]
    public function adminQuote(Request $request): JsonResponse
    {
        try {
            // Require admin role
            $user = $this->authenticator->requireRole($request, 'admin');

            // Admin-specific logic here...

            return new JsonResponse([
                'success' => true,
                'message' => 'Admin quote endpoint',
                'user' => $user
            ]);

        } catch (Exception $e) {
            return new JsonResponse([
                'success' => false,
                'error' => $e->getMessage()
            ], $e->getCode() === 401 ? 401 : 500);
        }
    }
}
