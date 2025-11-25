<?php

declare(strict_types=1);

namespace App\Service\Login;

use App\Enum\LoginStrategy;
use App\Interface\LoginInterface;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;

/**
 * Login Strategy Factory
 * Returns the appropriate login strategy implementation based on strategy type
 * Similar to PaymentStrategyFactory pattern
 */
class LoginStrategyFactory
{
    public function __construct(
        private EmailLoginStrategy $emailLoginStrategy,
        private GoogleLoginStrategy $googleLoginStrategy
    ) {}

    /**
     * Execute login based on strategy
     *
     * @param string $strategyValue Login strategy value
     * @param Request $request HTTP request
     * @param array $data Request payload
     * @return JsonResponse Login response with JWT token or error
     */
    public function login(string $strategyValue, Request $request, array $data): JsonResponse
    {
        // Get the strategy implementation
        $strategy = $this->getStrategy($strategyValue);

        if (!$strategy) {
            return new JsonResponse([
                'error' => "Invalid or unimplemented login strategy '{$strategyValue}'"
            ], 400);
        }

        // Execute login with the strategy
        return $strategy->login($request, $data);
    }

    /**
     * Get the appropriate login strategy implementation
     *
     * @param string $strategyValue Strategy value from request
     * @return LoginInterface|null Strategy implementation or null if invalid
     */
    private function getStrategy(string $strategyValue): ?LoginInterface
    {
        // Convert string to enum
        $strategyEnum = LoginStrategy::fromString($strategyValue);

        if (!$strategyEnum) {
            return null;
        }

        // Check if strategy is implemented
        if (!$strategyEnum->isImplemented()) {
            return null;
        }

        // Return appropriate strategy implementation
        return match ($strategyEnum) {
            LoginStrategy::EMAIL => $this->emailLoginStrategy,
            LoginStrategy::GOOGLE => $this->googleLoginStrategy,
            default => null
        };
    }
}
