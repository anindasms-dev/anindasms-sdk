<?php
declare(strict_types=1);

namespace AnindaSms;

/**
 * API hatası. getStatus() HTTP kodunu, getData() varsa ek bilgiyi döner
 * (ör. 402'de ['balance' => .., 'required' => ..]).
 */
class AnindaSmsException extends \RuntimeException
{
    private int $status;
    private ?array $data;

    public function __construct(string $message, int $status = 0, ?array $data = null)
    {
        parent::__construct($message, $status);
        $this->status = $status;
        $this->data = $data;
    }

    public function getStatus(): int
    {
        return $this->status;
    }

    public function getData(): ?array
    {
        return $this->data;
    }
}
