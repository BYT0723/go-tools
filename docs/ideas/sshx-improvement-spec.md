# Spec: sshx SSH Server Improvement

## Objective

Transform `transport/sshx` from an SSH server skeleton into a fully functional SSH server suitable for embedding
as a management channel in Go services. Core capabilities: PTY shell, remote command execution, and TCP port
forwarding (-L/-R) on the server side.

**User:** Go developers embedding SSH for production debugging/management of their services.

**Success looks like:** `ssh user@host` gets a working bash shell, `ssh user@host <cmd>` executes commands,
`ssh -L 8080:localhost:80 user@host` tunnels TCP, `ssh -R 9999:localhost:22 user@host` reverse-tunnels TCP.

## Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.25 |
| SSH protocol | `golang.org/x/crypto/ssh` | v0.46.0 |
| PTY | `github.com/creack/pty` | latest (new dependency) |
| Testing | `github.com/stretchr/testify` | v1.11.1 |
| Logging | `github.com/BYT0723/go-tools/logx` | (internal) |
| Lifecycle | `github.com/BYT0723/go-tools/srvx` | (internal) |

**Dependency decision:** Use `creack/pty` for PTY allocation. It abstracts Linux `/dev/ptmx` and macOS `openpty`
behind a single `pty.Start(cmd)` call, reducing ~200 lines of platform-specific code.

## Commands

```
Build:     go build ./transport/sshx/...
Test:      go test -race ./transport/sshx/...
Test All:  go test -race -count=1 -timeout 120s $(go list ./... | grep -v transport/ssh)  (legacy)
           go test -race -count=1 -timeout 120s ./...  (target, when CI is updated)
Vet:       go vet ./transport/sshx/...
Lint:      go vet ./...
```

## Project Structure

```
transport/sshx/
├── server.go         → Server struct, NewServer, Start/Stop, Init/Run/Destroy,
│                        handleConn, defaultHandler (shell/exec/pty)
├── forward.go        → direct-tcpip handler (-L), tcpip-forward handler (-R),
│                        forwarding lifecycle (bind/unbind per-connection)
├── option.go         → Option type, WithAddr, WithHostKey, WithPasswordAuth,
│                        WithPublicKeyAuth, WithUser
├── server_test.go    → Existing tests + new tests for shell, exec, forwarding
└── doc.go            → Package documentation
```

**Rationale:** `forward.go` is extracted because TCP forwarding involves distinct concerns —
listener management, goroutine lifecycle, channel piping — that benefit from isolation.
Shell/exec/PTY stays in `server.go` as it builds on the existing `defaultHandler`.

## Code Style

Follow existing project conventions:
- No comments unless critical (package doc is OK)
- Addressable types: `Server`, `Option`, `Handler`
- Unexported internals: handler functions, forwarding helpers
- Functional options pattern for configuration
- Thread safety via `sync.Mutex` on `Server` struct

Example of new Handler pattern expected:

```go
type Server struct {
    // ... existing fields ...
    publicKeyAuth func(conn xssh.ConnMetadata, key xssh.PublicKey) bool
}
```

## Testing Strategy

- **Framework:** `stretchr/testify/assert` + Go standard `testing`
- **Test location:** In-package (`package sshx`), file `server_test.go`
- **Network:** All network tests use loopback; PTY tests spawn bash locally
- **Test levels:**
  - Unit: Option behavior, auth callbacks, handler dispatch logic
  - Integration: Full SSH client→server dial, session with shell/exec/forwarding
- **Coverage target:** >80% on server.go, forward.go, option.go

## Boundaries

### Always Do
- Run `go test -race ./transport/sshx/...` before declaring a step complete
- Follow existing naming conventions (`TestXxx` for tests, `WithXxx` for options)
- Validate inputs in public functions (nil checks, empty addr)
- Close channels, listeners, and connections on cleanup

### Ask First
- Adding new dependencies beyond `creack/pty`
- Changing the public API (exported types/functions)
- Modifying CI configuration (`.github/workflows/test.yml`)

### Never Do
- Commit secrets or private keys
- Panic in library code (return errors)
- Rely on real network endpoints in tests
- Remove existing tests

## Success Criteria

- [ ] `ssh user@host` with `-p 2222` connects and provides interactive bash shell
- [ ] `ssh user@host echo hello` returns "hello\n" with exit code 0
- [ ] `ssh -L 8080:internal:80 user@host` forwards TCP connections correctly
- [ ] `ssh -R 9999:localhost:22 user@host` reverse-forwards correctly
- [ ] Server cleans up all goroutines and listeners on Stop/Destroy
- [ ] `go test -race ./transport/sshx/...` passes with zero races
- [ ] CI includes `transport/sshx` in build and test matrix

## Open Questions

- Should we support `certificate` authentication type? Deferring — public key callback can handle it.
- Should shell path be configurable? Yes, via `WithShell` option.
- How many concurrent forwards? Unlimited, bounded only by OS file descriptors.
