<?php

declare(strict_types=1);

namespace App\Client;

use Pricing\V1\PricingServiceClient;
use Pricing\V1\QuoteRequest;
use Pricing\V1\QuoteResponse;
use RuntimeException;

final class PricingGrpcClient
{
    private PricingServiceClient $client;

    public function __construct(string $target)
    {
        $this->client = new PricingServiceClient($target, [
            'credentials' => \Grpc\ChannelCredentials::createInsecure(),
        ]);
    }

    public function quote(float $amount, int $termMonths, ?array $options = null): array
    {
        $req = new QuoteRequest();
        $resp = new QuoteResponse();
        $req->setAmount($amount);
        $req->setTermMonths($termMonths);

        // Handle options array for riskScore
        $riskScore = $options['riskScore'] ?? null;
        if (null !== $riskScore) {
            $req->setRiskScore($riskScore);
        }

        // Add JWT token to metadata for authentication
        $metadata = [];
        if (isset($options['jwt_token'])) {
            $metadata['authorization'] = ['Bearer ' . $options['jwt_token']];
        }

        [$resp, $status] = $this->client->Quote($req, $metadata)->wait();
        if ($status->code !== 0) { // 0 is STATUS_OK
            throw new RuntimeException('gRPC error: '.$status->details, $status->code);
        }

        return [
            'interestRate' => $resp->getInterestRate(),
            'apr' => $resp->getApr(),
            'monthlyPayment' => $resp->getMonthlyPayment(),
        ];
    }
}
