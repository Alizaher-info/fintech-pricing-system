<?php

declare(strict_types=1);

namespace App\Service\Register;

use App\Enum\LoginStrategy;
use App\Interface\RegisterInterface;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;

/**
 * Register Strategy Factory
 * Returns the appropriate registration strategy based on LoginStrategy enum
 */
class RegisterStrategyFactory
{
    public function __construct(
        private EmailRegisterStrategy $emailRegisterStrategy
    ) {}

    /**
     * Register user using the specified strategy
     *
     * @param string $strategyValue Strategy string value (email, google, etc.)
     * @param Request $request HTTP request
     * @param array $data Decoded request data
     * @return JsonResponse Registration response
     */
    public function register(string $strategyValue, Request $request, array $data): JsonResponse
    {
        // Only email registration supported
        // Google OAuth handled by login endpoint (auto-creates user if doesn't exist)
        if ($strategyValue !== 'email') {
            return new JsonResponse([
                'error' => "Only email registration is supported. Use 'Sign in with Google' to create account via Google.",
                'hint' => 'Google OAuth registration is handled through the login endpoint'
            ], 400);
        }

        // Execute email registration
        return $this->emailRegisterStrategy->register($request, $data);
    }
}
