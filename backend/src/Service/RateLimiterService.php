<?php

namespace App\Service;

use Predis\Client as RedisClient;
use Symfony\Contracts\Cache\CacheInterface;

/**
 * Redis-based rate limiter for brute force protection
 */
class RateLimiterService
{
    private const MAX_ATTEMPTS_PER_IP = 10; // Max attempts per IP in time window
    private const MAX_ATTEMPTS_PER_USER = 5; // Max attempts per user account
    private const TIME_WINDOW = 900; // 15 minutes in seconds
    private const USER_LOCKOUT_TIME = 900; // 15 minutes user lockout
    private const IP_LOCKOUT_TIME = 3600; // 1 hour IP lockout

    public function __construct(
        private CacheInterface $rateLimiterCache,
        private RedisClient $redisClient
    ) {}

    /**
     * Check if IP is rate limited
     */
    public function isIpRateLimited(string $ipAddress): bool
    {
        $key = "ip_attempts:{$ipAddress}";
        $lockoutKey = "ip_locked:{$ipAddress}";
        
        // Check if IP is currently locked out
        if ($this->redisClient->exists($lockoutKey)) {
            return true;
        }
        
        // Get current attempt count
        $attempts = (int) $this->redisClient->get($key);
        
        return $attempts >= self::MAX_ATTEMPTS_PER_IP;
    }

    /**
     * Check if user account is locked
     */
    public function isUserLocked(int $userId): bool
    {
        $lockoutKey = "user_locked:{$userId}";
        return $this->redisClient->exists($lockoutKey);
    }

    /**
     * Record a failed login attempt
     */
    public function recordFailedAttempt(string $ipAddress, ?int $userId = null): array
    {
        $ipKey = "ip_attempts:{$ipAddress}";
        $userKey = $userId ? "user_attempts:{$userId}" : null;
        
        // Increment IP attempt counter
        $ipAttempts = $this->redisClient->incr($ipKey);
        $this->redisClient->expire($ipKey, self::TIME_WINDOW);
        
        $userAttempts = 0;
        if ($userKey) {
            $userAttempts = $this->redisClient->incr($userKey);
            $this->redisClient->expire($userKey, self::TIME_WINDOW);
        }
        
        // Check for lockouts
        $ipLocked = false;
        $userLocked = false;
        
        // Lock IP if exceeded
        if ($ipAttempts >= self::MAX_ATTEMPTS_PER_IP) {
            $ipLockKey = "ip_locked:{$ipAddress}";
            $this->redisClient->setex($ipLockKey, self::IP_LOCKOUT_TIME, '1');
            $ipLocked = true;
        }
        
        // Lock user if exceeded
        if ($userId && $userAttempts >= self::MAX_ATTEMPTS_PER_USER) {
            $userLockKey = "user_locked:{$userId}";
            $this->redisClient->setex($userLockKey, self::USER_LOCKOUT_TIME, '1');
            $userLocked = true;
        }
        
        return [
            'ip_attempts' => $ipAttempts,
            'user_attempts' => $userAttempts,
            'ip_locked' => $ipLocked,
            'user_locked' => $userLocked,
            'ip_remaining' => max(0, self::MAX_ATTEMPTS_PER_IP - $ipAttempts),
            'user_remaining' => $userId ? max(0, self::MAX_ATTEMPTS_PER_USER - $userAttempts) : null
        ];
    }

    /**
     * Clear rate limiting for successful login
     */
    public function clearAttempts(string $ipAddress, ?int $userId = null): void
    {
        $ipKey = "ip_attempts:{$ipAddress}";
        $this->redisClient->del($ipKey);
        
        if ($userId) {
            $userKey = "user_attempts:{$userId}";
            $this->redisClient->del($userKey);
        }
    }

    /**
     * Get remaining attempts info
     */
    public function getRemainingAttempts(string $ipAddress, ?int $userId = null): array
    {
        $ipKey = "ip_attempts:{$ipAddress}";
        $ipAttempts = (int) $this->redisClient->get($ipKey);
        
        $userAttempts = 0;
        if ($userId) {
            $userKey = "user_attempts:{$userId}";
            $userAttempts = (int) $this->redisClient->get($userKey);
        }
        
        return [
            'ip_remaining' => max(0, self::MAX_ATTEMPTS_PER_IP - $ipAttempts),
            'user_remaining' => $userId ? max(0, self::MAX_ATTEMPTS_PER_USER - $userAttempts) : null,
            'ip_locked' => $this->isIpRateLimited($ipAddress),
            'user_locked' => $userId ? $this->isUserLocked($userId) : false
        ];
    }

    /**
     * Get progressive delay based on attempts
     */
    public function getProgressiveDelay(int $attempts): int
    {
        if ($attempts <= 3) return 0;
        if ($attempts <= 5) return 2; // 2 seconds
        if ($attempts <= 7) return 5; // 5 seconds
        if ($attempts <= 10) return 10; // 10 seconds
        return 30; // 30 seconds for many attempts
    }

    /**
     * Get lockout time remaining
     */
    public function getLockoutTimeRemaining(string $ipAddress, ?int $userId = null): array
    {
        $ipLockKey = "ip_locked:{$ipAddress}";
        // Get TTLs
        $ipTtl = $this->redisClient->ttl($ipLockKey);
        // Default user TTL
        $userTtl = -1;
        if ($userId) {
            $userLockKey = "user_locked:{$userId}";
            $userTtl = $this->redisClient->ttl($userLockKey);
        }
        
        return [
            'ip_lockout_remaining' => $ipTtl > 0 ? $ipTtl : 0,
            'user_lockout_remaining' => $userTtl > 0 ? $userTtl : 0
        ];
    }
}