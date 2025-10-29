<?php

// Simple gRPC connection test
require_once 'vendor/autoload.php';

echo "Testing gRPC Connection...\n";

try {
    // Test if protobuf classes exist
    if (!class_exists('Pricing\V1\PricingServiceClient')) {
        echo "ERROR: PricingServiceClient class not found\n";
        echo "Available classes:\n";
        $grpcFiles = glob('src/Grpc/Pricing/V1/*.php');
        foreach ($grpcFiles as $file) {
            echo "Found: " . basename($file) . "\n";
        }
        exit(1);
    }
    
    echo "✓ PricingServiceClient class found\n";
    
    // Test gRPC connection
    $client = new \Pricing\V1\PricingServiceClient('pricing-api:50051', [
        'credentials' => \Grpc\ChannelCredentials::createInsecure(),
    ]);
    
    echo "✓ gRPC client created\n";
    
    // Test simple request with JWT token
    $request = new \Pricing\V1\QuoteRequest();
    $request->setAmount(1000);
    $request->setTermMonths(12);
    $request->setRiskScore(700);
    
    echo "✓ Request created\n";
    
    // Add JWT metadata for authentication
    $metadata = [
        'authorization' => ['Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJmaW50ZWNoLWFwaSIsImF1ZCI6ImZpbnRlY2gtYXBpIiwiaWF0IjoxNzYxNzM2MTIzLCJleHAiOjE3NjE4MjI1MjMsIm5iZiI6MTc2MTczNjEyMywidXNlcl9pZCI6MSwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsInJvbGUiOiJhZG1pbiJ9.dU8fJHxgVUnGiEnqkSbM0Zj9bxIo2cEsPBv']
    ];
    
    [$response, $status] = $client->Quote($request, $metadata)->wait();
    
    if ($status->code !== 0) {
        echo "ERROR: gRPC call failed - Code: {$status->code}, Details: {$status->details}\n";
        exit(1);
    }
    
    echo "✓ gRPC call successful!\n";
    echo "Interest Rate: " . $response->getInterestRate() . "\n";
    echo "APR: " . $response->getApr() . "\n";
    echo "Monthly Payment: " . $response->getMonthlyPayment() . "\n";
    
} catch (Exception $e) {
    echo "ERROR: " . $e->getMessage() . "\n";
    echo "File: " . $e->getFile() . "\n";
    echo "Line: " . $e->getLine() . "\n";
    echo "Trace:\n" . $e->getTraceAsString() . "\n";
    exit(1);
}

echo "✓ All tests passed!\n";