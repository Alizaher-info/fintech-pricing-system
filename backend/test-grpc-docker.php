<?php

require_once __DIR__ . '/vendor/autoload.php';

use App\Client\PricingGrpcClient;

try {
    echo "Testing gRPC connection...\n";

    // Verwende den Docker Service Namen statt localhost
    $client = new PricingGrpcClient('pricing-api:50051');
    echo "Client created successfully.\n";

    // Test Quote Request
    echo "Sending quote request...\n";
    $result = $client->quote(10000.0, 12);

    echo "Success! Response received:\n";
    print_r($result);

} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
    echo "Code: " . $e->getCode() . "\n";
    echo "File: " . $e->getFile() . ":" . $e->getLine() . "\n";
}
