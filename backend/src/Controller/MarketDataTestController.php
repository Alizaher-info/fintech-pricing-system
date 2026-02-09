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
        $startTime = microtime(true);
        $timings = [];

        // Test Redis directly since market-trading writes to it
        try {
            // Timing: Redis Connection
            $connectStart = microtime(true);
            $redis = new RedisClient([
                'scheme' => 'tcp',
                'host'   => 'redis',
                'port'   => 6379,
            ]);
            $timings['redis_connect_ms'] = round((microtime(true) - $connectStart) * 1000, 2);

            // Timing: BTC Price Fetch
            $btcStart = microtime(true);
            $btcPrice = $redis->get('price:BTC');
            $timings['btc_fetch_ms'] = round((microtime(true) - $btcStart) * 1000, 2);

            // Timing: ETH Price Fetch
            $ethStart = microtime(true);
            $ethPrice = $redis->get('price:ETH');
            $timings['eth_fetch_ms'] = round((microtime(true) - $ethStart) * 1000, 2);

            // Timing: SOL Price Fetch
            $solStart = microtime(true);
            $solPrice = $redis->get('price:SOL');
            $timings['sol_fetch_ms'] = round((microtime(true) - $solStart) * 1000, 2);

            // Timing: JSON Parsing
            $parseStart = microtime(true);
            $btcData = $btcPrice ? json_decode($btcPrice, true) : null;
            $ethData = $ethPrice ? json_decode($ethPrice, true) : null;
            $solData = $solPrice ? json_decode($solPrice, true) : null;
            $timings['json_parse_ms'] = round((microtime(true) - $parseStart) * 1000, 2);

            // Calculate data age (how old is the cached data)
            $dataAge = [];
            if ($btcData && isset($btcData['timestamp'])) {
                $dataAge['BTC'] = round((time() - strtotime($btcData['timestamp'])) * 1000, 2) . 'ms ago';
            }
            if ($ethData && isset($ethData['timestamp'])) {
                $dataAge['ETH'] = round((time() - strtotime($ethData['timestamp'])) * 1000, 2) . 'ms ago';
            }
            if ($solData && isset($solData['timestamp'])) {
                $dataAge['SOL'] = round((time() - strtotime($solData['timestamp'])) * 1000, 2) . 'ms ago';
            }

            $totalTime = round((microtime(true) - $startTime) * 1000, 2);
            $timings['total_request_ms'] = $totalTime;

            return $this->json([
                'status' => 'success',
                'source' => 'Redis (cached by market-trading service)',
                'data' => [
                    'BTC' => $btcData,
                    'ETH' => $ethData,
                    'SOL' => $solData,
                ],
                'performance' => [
                    'timings' => $timings,
                    'data_age' => $dataAge,
                    'total_ms' => $totalTime,
                ],
                'timestamp' => date('Y-m-d H:i:s')
            ]);
        } catch (\Exception $e) {
            return $this->json([
                'status' => 'error',
                'message' => $e->getMessage(),
                'timings' => $timings
            ], 500);
        }
    }
}
