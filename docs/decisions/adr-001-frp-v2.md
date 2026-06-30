# ADR-001: Use FRP V2 Wire Protocol for frpx Client

## Status
Accepted

## Date
2026-06-30

## Context
The `transport/frpx` package needs to connect to standard FRP servers (github.com/fatedier/frp).
FRP supports two wire protocol versions:

- **V1**: No handshake. First byte is the Login message. Uses gob-based binary encoding
  (golib/msg/json library). XOR stream cipher with salt exchange for encryption.
- **V2**: Magic bytes (`FRP\x00\x02\r\n`) → ClientHello/ServerHello negotiation →
  AEAD (AES-256-GCM) encryption with HKDF key derivation. JSON encoding with
  8-byte frame headers (uint16 type + uint16 flags + uint32 length).

## Decision
Use **V2 wire protocol** exclusively.

## Alternatives Considered

### V1 Protocol
- Pros: Simpler frame format (4B + 1B type + body), XOR crypto is trivial to implement
- Cons: Depends on golib's gob encoding (external dependency), being phased out
- Rejected: V2 is the modern path; maintaining V1 backward compat for a new library adds debt

### Custom Protocol (non-FRP-compatible)
- Pros: Maximum simplicity, full control
- Cons: Zero interoperability with existing FRP ecosystem
- Rejected: The entire purpose is FRP compatibility

## Consequences
- **Zero external dependencies** — AEAD (crypto/aes, crypto/cipher), SHA256, HMAC, HKDF all from stdlib
- **Non-trivial crypto setup** — V2 requires HKDF key derivation, AES-256-GCM stream framing,
  nonce management, and transcript hash computation
- **Frame overhead** — 8-byte header per frame + 16-byte AEAD tag per GCM chunk vs V1's 5-byte header + XOR
- **Server compatibility** — Standard frps auto-detects V1/V2 via magic bytes; our V2 client
  connects seamlessly

## Key Implementation Decisions
- HKDF implemented manually (stdlib has no `crypto/hkdf` package) using HMAC-SHA256 extract+expand
- ClientHello sends only `aes-256-gcm` in crypto algorithms (not xchacha20-poly1305)
- ClientHello sets `TLS: false, TCPMux: false` (MVP scope)
- Mock frps returns empty crypto algorithm (no AEAD) to simplify test setup

## See Also
- [FRP V2 protocol spec](https://github.com/fatedier/frp/blob/dev/pkg/proto/wire/wire.go)
- `/docs/ideas/frpx-client-spec.md`
