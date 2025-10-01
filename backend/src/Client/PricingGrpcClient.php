<?php

declare(strict_types=1);

namespace App\Client;

use Grpc\ChannelCredentials;
use Pricing\V1\PricingServiceClient;
use Pricing\V1\QuoteRequest;
use RuntimeException;

final class PricingGrpcClient
{
    private PricingServiceClient $client;

    public function __construct(string $target)
    {
        $this->client = new PricingServiceClient($target, [
            'credentials' => ChannelCredentials::createInsecure(),
        ]);
    }

    public function quote(float $amount, int $termMonths, ?array $options = null): array
    {
        $req = new QuoteRequest();
        $req->setAmount($amount);
        $req->setTermMonths($termMonths);

        // Handle options array for riskScore
        $riskScore = $options['riskScore'] ?? null;
        if (null !== $riskScore) {
            $req->setRiskScore($riskScore);
        }

        [$resp, $status] = $this->client->Quote($req)->wait();
        if ($status->code !== \Grpc\STATUS_OK) {
            throw new RuntimeException('gRPC error: '.$status->details, $status->code);
        }

        return [
            'interestRate' => $resp->getInterestRate(),
            'apr' => $resp->getApr(),
            'monthlyPayment' => $resp->getMonthlyPayment(),
        ];
    }
}