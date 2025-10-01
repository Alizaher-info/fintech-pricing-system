<?php

declare(strict_types=1);

namespace App\Controller;

use App\Client\PricingGrpcClientFactory;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\Routing\Annotation\Route;

final class QuoteTestController
{
    #[Route('/quote-test', methods: ['GET'])]
    public function __invoke(): JsonResponse
    {
        $client = PricingGrpcClientFactory::create();
        $data = $client->quote(10000.0, 24, null);

        return new JsonResponse($data);
    }
}
