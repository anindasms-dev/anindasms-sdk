"""Resmi AnındaSMS Python SDK'sı (bağımlılıksız, yalnızca standart kütüphane)."""

from .client import AnindaSMSError, Client

__all__ = ["Client", "AnindaSMSError"]
__version__ = "1.0.0"
