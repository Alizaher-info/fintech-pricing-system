<?php

namespace App\EventListener;

use App\Entity\SecurityEvent;
use App\Event\UserLoginSuccessEvent;
use App\Event\UserLoginFailedEvent;
use App\Repository\SecurityEventRepository;
use App\Repository\UserRepository;
use App\Service\RateLimiterService;
use Symfony\Component\EventDispatcher\Attribute\AsEventListener;

class SecurityEventListener
{
    public function __construct(
        private SecurityEventRepository $securityEventRepository,
        private UserRepository $userRepository,
        private RateLimiterService $rateLimiterService
    ) {}

    #[AsEventListener(event: UserLoginSuccessEvent::class)]
    public function onLoginSuccess(UserLoginSuccessEvent $event): void
    {
        $user = $event->getUser();
        $ipAddress = $event->getIpAddress();
        
        // Clear rate limiting attempts for successful login
        $this->rateLimiterService->clearAttempts($ipAddress, $user->getId());
        
        // Update user's last login time
        $user->setLastLoginAt($event->getOccurredAt());
        $this->userRepository->save($user, true);
        
        // Create security event record
        $securityEvent = new SecurityEvent();
        $securityEvent->setEventType('LOGIN_SUCCESS');
        $securityEvent->setDescription('User successfully logged in');
        $securityEvent->setIpAddress($ipAddress);
        $securityEvent->setUserAgent($event->getUserAgent());
        $securityEvent->setSeverity('LOW');
        $securityEvent->setCreatedAt($event->getOccurredAt());
        $securityEvent->setUser($user);
        
        $this->securityEventRepository->save($securityEvent, true);
    }

    #[AsEventListener(event: UserLoginFailedEvent::class)]
    public function onLoginFailed(UserLoginFailedEvent $event): void
    {
        $ipAddress = $event->getIpAddress();
        $user = $event->getUser();
        $userId = $user?->getId();
        
        // Record failed attempt in Redis for rate limiting
        $rateLimitInfo = $this->rateLimiterService->recordFailedAttempt($ipAddress, $userId);
        
        // Determine severity based on rate limiting status
        $severity = 'MEDIUM';
        if ($rateLimitInfo['ip_locked'] || $rateLimitInfo['user_locked']) {
            $severity = 'HIGH';
        } elseif ($rateLimitInfo['ip_attempts'] > 5 || $rateLimitInfo['user_attempts'] > 3) {
            $severity = 'MEDIUM';
        }
        
        // Enhanced description with rate limit info
        $description = sprintf(
            'Login failed for email %s: %s (IP attempts: %d/%d, User attempts: %d/%d)',
            $event->getEmail(),
            $event->getReason(),
            $rateLimitInfo['ip_attempts'],
            10, // MAX_ATTEMPTS_PER_IP
            $rateLimitInfo['user_attempts'],
            5   // MAX_ATTEMPTS_PER_USER
        );
        
        if ($rateLimitInfo['ip_locked']) {
            $description .= ' - IP ADDRESS LOCKED';
        }
        if ($rateLimitInfo['user_locked']) {
            $description .= ' - USER ACCOUNT LOCKED';
        }
        
        // Create security event record for failed login
        $securityEvent = new SecurityEvent();
        $securityEvent->setEventType('LOGIN_FAILED');
        $securityEvent->setDescription($description);
        $securityEvent->setIpAddress($ipAddress);
        $securityEvent->setUserAgent($event->getUserAgent());
        $securityEvent->setSeverity($severity);
        $securityEvent->setCreatedAt($event->getOccurredAt());
        $securityEvent->setUser($user); // null if user doesn't exist
        
        $this->securityEventRepository->save($securityEvent, true);
    }
}