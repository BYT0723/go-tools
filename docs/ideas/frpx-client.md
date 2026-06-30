# frpx FRP Client

## Problem Statement
如何将 FRP 客户端能力嵌入 Go 服务，使 NAT 后的服务能通过标准 frps 暴露到公网，无需单独部署 frpc 二进制？

## Recommended Direction
兼容标准 frps (github.com/fatedier/frp) 的 TCP + HTTP/HTTPS 客户端库。
遵循 sshx 模式：功能选项（`Option`）、`srvx.Service` 生命周期、`WaitGroup` 优雅关闭。

选择原因：fatedier/frp 只有独立二进制，没有可嵌入的 Go 库。这个是填补空白的唯一方案。
TCP + HTTP/HTTPS 覆盖 95% 使用场景，token 认证与标准 frps 一致。

## Key Assumptions to Validate
- [ ] FRP 二进制协议（4B 帧长 + 1B 类型 + JSON body）与目标 frps 版本兼容
- [ ] WorkConn 回连模型（frps 主动向 client 发起 TCP 连接）在目标网络环境下可行
- [ ] HTTP vhost 路由（frps → WorkConn → client → local service）延迟可接受
- [ ] Token 认证足够，v1 不需要 mTLS/OIDC

## MVP Scope

### 包含
- **TCP 代理**: 注册 TCP 端口映射，双向 pipe 流量转发
- **HTTP 代理**: vhost 域名路由到本地 HTTP 服务
- **HTTPS 代理**: vhost 域名路由到本地 HTTPS 服务
- **Token 认证**: Login 消息携带 token
- **心跳**: 定期 ping/pong 保活
- **多代理**: 单个 Client 实例管理多个 proxy 配置
- **srvx 集成**: Name/Init/Run/Destroy
- **功能选项**: WithServerAddr, WithToken, WithProxy 等
- **断线重连**: 基础指数退避

### 不包含（v1）
- UDP/STCP/XTCP 代理 — TCP+HTTP 覆盖主要场景
- TLS 加密/mTLS — 内网场景不需要，公网可走外部 TLS 隧道
- OIDC/Token 之外的认证 — token 是 frps 最常用方式
- 压缩 — 控制通道数据量小，WorkConn 是直通不处理的
- 多 frps 故障转移 — 单实例但可部署多个 client 实例
- FRP 服务端 (frps) — 只做 client 端

## Not Doing (and Why)
- **FRP 服务端实现** — fatedier/frp 的 frps 已经很成熟，不需要重复实现。client 是缺失的拼图
- **自定义隧道协议** — 不搞"更高层次抽象"，先让 FRP 跑起来。抽象等第二个协议实现再提取
- **WebSocket 控制通道** — 会失去与标准 frps 的兼容性。如果确实需要，可作为后续的 `transport/` 参数
- **Dashboard/管理 API** — 嵌入库不需要自己的 Web UI，宿主服务自行管理

## Open Questions
- 包名：`transport/frpx` 还是 `transport/tunnel/frp`？倾向 `transport/frpx` 与 sshx 一致
- FRP 协议版本锁定：固定兼容某个 frps 版本范围？还是跟踪最新版本？
- WorkConn 模型在 Docker/K8s 环境下的可用性：frps 需要能回连到 client pod
- HTTP vhost 代理中 `Host` 头改写：由 frps 处理还是 client 处理？
