<?php
namespace App\Client;

/**
 * Factory for creating gRPC-based pricing clients
 *
 * This factory provides static methods to create pricing clients
 * without requiring Symfony dependency injection.
 */
final class PricingGrpcClientFactory
{
        public static function create(): PricingGrpcClient
    {
        $target = 'pricing-api:50051';
        return new PricingGrpcClient($target);
    }
}
