<?php

namespace App\Event;

use App\Entity\User;
use Symfony\Contracts\EventDispatcher\Event;

class UserLoginFailedEvent extends Event
{
    public function __construct(
        private string $email,
        private string $reason,
        private string $ipAddress,
        private ?string $userAgent = null,
        private ?User $user = null,
        private \DateTimeImmutable $occurredAt = new \DateTimeImmutable()
    ) {}

    public function getEmail(): string
    {
        return $this->email;
    }

    public function getReason(): string
    {
        return $this->reason;
    }

    public function getIpAddress(): string
    {
        return $this->ipAddress;
    }

    public function getUserAgent(): ?string
    {
        return $this->userAgent;
    }

    public function getUser(): ?User
    {
        return $this->user;
    }
    //to get the exact time of the event occurrence
    public function getOccurredAt(): \DateTimeImmutable
    {           
        return $this->occurredAt;
    }
}