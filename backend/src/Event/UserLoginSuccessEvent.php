<?php

namespace App\Event;

use App\Entity\User;
use Symfony\Contracts\EventDispatcher\Event;

class UserLoginSuccessEvent extends Event
{
    public function __construct(
        private User $user,
        private string $ipAddress,
        private ?string $userAgent = null,
        private \DateTimeImmutable $occurredAt = new \DateTimeImmutable()
    ) {}

    public function getUser(): User
    {
        return $this->user;
    }

    public function getIpAddress(): string
    {
        return $this->ipAddress;
    }

    public function getUserAgent(): ?string
    {
        return $this->userAgent;
    }

    public function getOccurredAt(): \DateTimeImmutable
    {
        return $this->occurredAt;
    }
}