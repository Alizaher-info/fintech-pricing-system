<?php

require_once __DIR__ . '/vendor/autoload.php';

// Test Login
$loginData = [
    'email' => 'admin@example.com',
    'password' => 'admin123'
];

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://localhost:8080/api/auth/login');
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($loginData));
curl_setopt($ch, CURLOPT_HTTPHEADER, ['Content-Type: application/json']);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

echo "Login Response (HTTP $httpCode):\n";
echo $response . "\n\n";

if ($httpCode === 200) {
    $data = json_decode($response, true);
    if (isset($data['token'])) {
        $token = $data['token'];
        echo "Token received: " . substr($token, 0, 50) . "...\n\n";

        // Test secured endpoint
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, 'http://localhost:8080/api/auth/validate');
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'Content-Type: application/json',
            'Authorization: Bearer ' . $token
        ]);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

        $validateResponse = curl_exec($ch);
        $validateHttpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        echo "Validate Response (HTTP $validateHttpCode):\n";
        echo $validateResponse . "\n";
    }
}
