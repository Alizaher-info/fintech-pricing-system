<?php

declare(strict_types=1);

namespace App\Interface;

use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;

/**
 * Login Interface
 * Contract for all login strategy implementations
 */
interface LoginInterface
{
    /**
     * Execute login with specific strategy
     *
     * @param Request $request HTTP request containing credentials
     * @param array $data Request payload with login credentials
     * @return JsonResponse Login response with JWT token or error
     */
    public function login(Request $request, array $data): JsonResponse;
}
