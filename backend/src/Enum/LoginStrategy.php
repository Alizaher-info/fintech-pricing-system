<?php

declare(strict_types=1);

namespace App\Enum;

/**
 * Login Strategy Enum
 * Defines all available login strategies
 */
enum LoginStrategy: string
{
    case EMAIL = 'email';
    case GOOGLE = 'google';
    case MICROSOFT = 'microsoft';
    case FACEBOOK = 'facebook';
    case OTP = 'otp';
    case BIOMETRIC = 'biometric';

    /**
     * Get all strategy values as array
     */
    public static function values(): array
    {
        return array_map(fn($case) => $case->value, self::cases());
    }

    /**
     * Check if strategy is implemented
     */
    public function isImplemented(): bool
    {
        return match ($this) {
            self::EMAIL, self::GOOGLE => true,
            default => false
        };
    }

    /**
     * Get strategy from string value
     */
    public static function fromString(string $value): ?self
    {
        return self::tryFrom($value);
    }
}
