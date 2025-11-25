<?php

declare(strict_types=1);

namespace App\Service;

use Symfony\Contracts\HttpClient\HttpClientInterface;

/**
 * Lightweight Google OAuth verification service
 * No heavy dependencies - uses HTTP client to verify tokens directly with Google
 */
class GoogleOAuthService
{
    private const GOOGLE_TOKEN_INFO_URL = 'https://oauth2.googleapis.com/tokeninfo';
    private const GOOGLE_USERINFO_URL = 'https://www.googleapis.com/oauth2/v2/userinfo';

    public function __construct(
        private HttpClientInterface $httpClient,
        private string $googleClientId,
        private string $googleClientSecret
    ) {}

    /**
     * Verify Google ID token and extract user information
     *
     * @param string $idToken The ID token from Google OAuth
     * @return array User info: ['email', 'name', 'picture', 'email_verified']
     * @throws \RuntimeException If token is invalid
     */
    public function verifyIdToken(string $idToken): array
    {
        try {
            $response = $this->httpClient->request('GET', self::GOOGLE_TOKEN_INFO_URL, [
                'query' => ['id_token' => $idToken]
            ]);

            if ($response->getStatusCode() !== 200) {
                throw new \RuntimeException('Invalid Google ID token');
            }

            $data = $response->toArray();

            // Verify the token belongs to our app
            if ($data['aud'] !== $this->googleClientId) {
                throw new \RuntimeException('Token audience mismatch');
            }

            // Verify token is not expired
            if (isset($data['exp']) && $data['exp'] < time()) {
                throw new \RuntimeException('Token expired');
            }

            return [
                'email' => $data['email'] ?? null,
                'name' => $data['name'] ?? null,
                'email_verified' => ($data['email_verified'] ?? 'false') === 'true',
                'google_id' => $data['sub'] ?? null,
            ];

        } catch (\Exception $e) {
            throw new \RuntimeException('Failed to verify Google token: ' . $e->getMessage());
        }
    }

    /**
     * Exchange authorization code for access token and get user info
     * This is called when user returns from Google OAuth consent screen
     *
     * @param string $authCode Authorization code from Google
     * @param string $redirectUri Must match the one registered in Google Console
     * @return array User info: ['email', 'name', 'picture', 'email_verified']
     * @throws \RuntimeException If exchange fails
     */
    public function exchangeCodeForUserInfo(string $authCode, string $redirectUri): array
    {
        try {
            // Step 1: Exchange code for tokens
            $tokenResponse = $this->httpClient->request('POST', 'https://oauth2.googleapis.com/token', [
                'body' => [
                    'code' => $authCode,
                    'client_id' => $this->googleClientId,
                    'client_secret' => $this->googleClientSecret,
                    'redirect_uri' => $redirectUri,
                    'grant_type' => 'authorization_code',
                ]
            ]);

            if ($tokenResponse->getStatusCode() !== 200) {
                throw new \RuntimeException('Failed to exchange authorization code');
            }

            $tokens = $tokenResponse->toArray();
            $accessToken = $tokens['access_token'] ?? null;
            $idToken = $tokens['id_token'] ?? null;

            if (!$accessToken) {
                throw new \RuntimeException('No access token received');
            }

            // Step 2: Get user info using access token
            $userInfoResponse = $this->httpClient->request('GET', self::GOOGLE_USERINFO_URL, [
                'headers' => [
                    'Authorization' => 'Bearer ' . $accessToken
                ]
            ]);

            if ($userInfoResponse->getStatusCode() !== 200) {
                throw new \RuntimeException('Failed to get user info');
            }

            $userInfo = $userInfoResponse->toArray();

            return [
                'email' => $userInfo['email'] ?? null,
                'name' => $userInfo['name'] ?? null,
                'picture' => $userInfo['picture'] ?? null,
                'email_verified' => $userInfo['verified_email'] ?? false,
                'google_id' => $userInfo['id'] ?? null,
            ];

        } catch (\Exception $e) {
            throw new \RuntimeException('OAuth exchange failed: ' . $e->getMessage());
        }
    }
}
