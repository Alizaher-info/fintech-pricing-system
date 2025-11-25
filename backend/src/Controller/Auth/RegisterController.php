<?php

declare(strict_types=1);

namespace App\Controller\Auth;

use App\Service\Register\RegisterStrategyFactory;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;
use Exception;

class RegisterController extends AbstractController
{
    public function __construct(
        private RegisterStrategyFactory $registerFactory
    ) {}

    #[Route('/api/register/{strategy}', name: 'api_register', methods: ['POST'])]
    public function register(Request $request, string $strategy): JsonResponse
    {
        try {
            $data = json_decode($request->getContent(), true);

            if ($data === null) {
                return new JsonResponse([
                    'error' => 'Invalid JSON in request body'
                ], 400);
            }

            // Use factory to handle registration strategy
            return $this->registerFactory->register($strategy, $request, $data);

        } catch (Exception $e) {
            return new JsonResponse([
                'error' => 'Registration failed: ' . $e->getMessage()
            ], 500);
        }
    }
}
