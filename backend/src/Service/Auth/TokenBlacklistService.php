<?php

declare(strict_types=1);

namespace App\Service\Auth;

use Predis\Client as RedisClient;
use Firebase\JWT\JWT;
use Firebase\JWT\Key;

class TokenBlacklistService
{
    private RedisClient $redis;
    private string $secretKey;

    public function __construct(RedisClient $redis)
    {
        $this->redis = $redis;
        $this->secretKey = $_ENV['JWT_SECRET'] ?? 'your-secret-key';
    }

    /**
     * Blacklist a specific token
     */
    public function blacklistToken(string $token): void
    {
        try {
            // Decode token to get expiration
            $decoded = JWT::decode($token, new Key($this->secretKey, 'HS256'));
            
            // Calculate TTL (time until token expires)
            $ttl = $decoded->exp - time();
            
            if ($ttl > 0) {
                // Store token hash in Redis with expiration
                $tokenHash = hash('sha256', $token);
                $this->redis->setex("blacklist:token:{$tokenHash}", $ttl, '1');
            }
        } catch (\Exception $e) {
            // Token is invalid, no need to blacklist
            return;
        }
    }

    /**
     * Check if a token is blacklisted
     */
    public function isBlacklisted(string $token): bool
    {
        $tokenHash = hash('sha256', $token);
        return (bool) $this->redis->exists("blacklist:token:{$tokenHash}");
    }

    /**
     * Blacklist all tokens for a specific user
     */
    public function blacklistAllUserTokens(int $userId): void
    {
        // Set a user blacklist timestamp
        // All tokens issued before this timestamp will be considered invalid
        $this->redis->set("blacklist:user:{$userId}", time());
    }

    /**
     * Check if a user has been logged out from all devices
     */
    public function isUserBlacklisted(int $userId, int $tokenIssuedAt): bool
    {
        $blacklistTimestamp = $this->redis->get("blacklist:user:{$userId}");
        
        if (!$blacklistTimestamp) {
            return false;
        }
        
        // If token was issued before the blacklist timestamp, it's invalid
        return $tokenIssuedAt < (int) $blacklistTimestamp;
    }

    /**
     * Remove user from blacklist (optional - for re-enabling old tokens)
     */
    public function removeUserBlacklist(int $userId): void
    {
        $this->redis->del("blacklist:user:{$userId}");
    }
}
