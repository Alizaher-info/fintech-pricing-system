<?php

declare(strict_types=1);

namespace App\Security;

use Firebase\JWT\JWT;
use Firebase\JWT\Key;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpKernel\Exception\UnauthorizedHttpException;

class JWTAuthenticator
{
    private string $secretKey;
    private string $issuer;

    public function __construct()
    {
        $this->secretKey = $_ENV['JWT_SECRET'] ?? 'your-secret-key';
        $this->issuer = $_ENV['JWT_ISSUER'] ?? 'fintech-api';
    }

    public function authenticate(Request $request): ?array
    {
        $authHeader = $request->headers->get('Authorization');

        if (!$authHeader || !str_starts_with($authHeader, 'Bearer ')) {
            return null;
        }

        $token = substr($authHeader, 7); // Remove "Bearer "

        try {
            // Firebase JWT decode
            $decoded = JWT::decode($token, new Key($this->secretKey, 'HS256'));

            return [
                'user_id' => $decoded->user_id,
                'email' => $decoded->email,
                'role' => $decoded->role,
                'exp' => $decoded->exp,
            ];
        } catch (\Exception $e) {
            throw new UnauthorizedHttpException('Bearer', 'Invalid token: ' . $e->getMessage());
        }
    }

    public function requireAuth(Request $request): array
    {
        $user = $this->authenticate($request);

        if (!$user) {
            throw new UnauthorizedHttpException('Bearer', 'Authentication required');
        }

        return $user;
    }

    public function requireRole(Request $request, string $requiredRole): array
    {
        $user = $this->requireAuth($request);

        if ($user['role'] !== $requiredRole && $user['role'] !== 'admin') {
            throw new UnauthorizedHttpException('Bearer', 'Insufficient permissions');
        }

        return $user;
    }

    public function generateToken(int $userId, string $email, string $role): string
    {
        $payload = [
            'iss' => $this->issuer,
            'aud' => $this->issuer,
            'iat' => time(),
            'exp' => time() + (24 * 60 * 60), // 24 hours
            'nbf' => time(),
            'user_id' => $userId,
            'email' => $email,
            'role' => $role
        ];

        return JWT::encode($payload, $this->secretKey, 'HS256');
    }
}
