<?php

declare(strict_types=1);

namespace App\Interface;

use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;

/**
 * Register Interface
 * Contract for all registration strategy implementations
 */
interface RegisterInterface
{
    /**
     * Register a user with the specific strategy
     *
     * @param Request $request The HTTP request
     * @param array $data Decoded request data
     * @return JsonResponse Registration response with token and user data
     */
    public function register(Request $request, array $data): JsonResponse;
}
