<?php

namespace App\Controller;

use Predis\Client as RedisClient;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\Routing\Annotation\Route;

class MarketDataTestController extends AbstractController
{
    #[Route('/api/market-data/test', name: 'market_data_test', methods: ['GET'])]
    public function test(): JsonResponse
    {
        // Test Redis directly since market-trading writes to it
        try {
            $redis = new RedisClient([
                'scheme' => 'tcp',
                'host'   => 'redis',
                'port'   => 6379,
            ]);

            $btcPrice = $redis->get('price:BTC');
            $ethPrice = $redis->get('price:ETH');
            $solPrice = $redis->get('price:SOL');

            return $this->json([
                'status' => 'success',
                'source' => 'Redis (cached by market-trading service)',
                'data' => [
                    'BTC' => $btcPrice ? json_decode($btcPrice, true) : null,
                    'ETH' => $ethPrice ? json_decode($ethPrice, true) : null,
                    'SOL' => $solPrice ? json_decode($solPrice, true) : null,
                ],
                'timestamp' => date('Y-m-d H:i:s')
            ]);
        } catch (\Exception $e) {
            return $this->json([
                'status' => 'error',
                'message' => $e->getMessage()
            ], 500);
        }
    }
}
