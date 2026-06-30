# go-tools

> a commonly used golang sdk package

<img src="./icon.png" width="200"/>

[![Test](https://github.com/BYT0723/go-tools/actions/workflows/test.yml/badge.svg)](https://github.com/BYT0723/go-tools/actions/workflows/test.yml)

## Modules

| Package | Description |
|---|---|
| `cfg` | Load configurations from local files or etcd using viper |
| `channelx` | Common channel utility functions (context-aware send/receive) |
| `contextx` | Common context utility functions (trace ID, request ID, logger injection) |
| `ds` | Common data structures (cache, counter, pool, stack, queue, map, set, hub, mutex) |
| `funny/graph` | ASCII graph plotting (heart, rose curves) |
| `i18n` | Wrappers for `go-i18n` with template and sprig support |
| `logx` | Logging facade with `zap` and `zerolog` implementations |
| `mathx` | Math utilities |
| `monitor` | Host/service monitoring framework |
| `monitor/component` | Reusable monitor lifecycle and alert rule evaluation |
| `monitor/ping` | ICMP ping monitor |
| `monitor/prometheus` | Prometheus metrics endpoint monitor |
| `monitor/web` | HTTP service monitor |
| `monitor/multi` | Multi-monitor aggregator |
| `osx` | OS utilities (terminal size, codepage decoding) |
| `packer` | Archive utilities (unzip, gzip/bzip2/xz/zstd decompressor) |
| `spider` | Web scraping utilities |
| `srvx` | Service lifecycle management (Init → Run → Destroy) |
| `transport/httpx` | HTTP client wrapper with encoder/decoder/compressor |
| `transport/httpx/middleware` | Gin and Echo middleware (logger, trace context injection) |
| `transport/sshx` | SSH server with PTY shell, exec, and port forwarding (`-L`/`-R`) |
| `transport/frpx` | FRP (Fast Reverse Proxy) V2 client for exposing local services behind NAT |
| `unsafex` | Unsafe operations with benchmarks |

## Development

### Testing

```bash
# All tests
go test ./...

# With race detector
go test -race ./...

# With coverage
go test -coverprofile=cover.cov ./...

# Specific package
go test ./ds/
```

All tests run offline using `httptest.Server` and mock data. No real network endpoints, ICMP, or external services required.

### Local Monitoring Stack (Optional)

For running monitors locally (not required for tests):

```bash
docker-compose up    # Prometheus + node_exporter
docker-compose down
```

## CI

GitHub Actions runs on every push/PR: `go vet`, `go build`, `go test -race -coverprofile`.

## License

MIT License - see [LICENSE](./LICENSE) for details.
