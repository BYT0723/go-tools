# sshx SSH Server Improvement

## Problem Statement
如何将当前 transport/sshx 的 SSH 服务器骨架（仅接受 session channel，shell/exec 均为占位符，
拒绝非 session 通道）改造为功能完整的 SSH 服务器，支持交互式终端、远程命令执行和 TCP 端口转发？

## Recommended Direction
完整功能实现（Shell + Exec + -R + -L）+ 插件化 Handler 架构（v2），接入 srvx 生命周期。

选择原因：sshx 嵌入 Go 服务作为管理通道是一个真实的生产需求。完整功能提供了立即可用的价值，
插件化架构（后续提取）允许使用者按需替换子组件，避免大而全的包袱。

目标架构：Init() 配置资源 → Run() 阻塞接受连接 → Destroy() 清理。
认证支持密码 + 公钥回调。PTY 支持 Linux 和 macOS，Windows 仅 exec + 转发。

## Key Assumptions to Validate
- [ ] PTY 分配在 Linux 和 macOS 上均工作正常（Linux 走 /dev/ptmx，macOS 走 openpty）
- [ ] 端口转发的 goroutine 生命周期能通过 SSH 连接的 context 级联取消来正确管理
- [ ] golang.org/x/crypto/ssh 的 PublicKeyCallback 能满足公钥认证需求（签名验证由库完成，
      我们只需做授权匹配）
- [ ] srvx.Init/Run/Destroy 模式与当前 Start/Stop 的语义能平滑映射（Run 阻塞等待 ctx.Done）

## MVP Scope

### 包含
- **密码认证**（已有，增强健壮性）
- **公钥认证**（callback 模式，用户自定义授权逻辑）
- **Shell**（PTY 分配 + /bin/bash，支持 pty-req、window-change、signal）
- **Exec**（os/exec 执行命令，返回 stdout/stderr + exit-status）
- **-L（direct-tcpip）**（接受 direct-tcpip 通道，出站连接目标，双向 pipe）
- **-R（tcpip-forward）**（处理 tcpip-forward 全局请求，监听绑定端口，
     将 inbound 连接转为 forwarded-tcpip 通道）
- **srvx 集成**（Init / Run / Destroy）
- **当前 server_test.go 扩展**（端到端测试 Shell、Exec、转发）
- **CI 恢复**（将 transport/sshx 从 CI 排除清单中移除）

### 不包含（v1）
- 插件化 Handler 接口提取（v2 做，等完整实现跑通后按实际 pattern 提取）
- Windows PTY（Windows 仅支持 exec + 转发）
- SFTP 支持
- X11 转发
- SSH 客户端实现
- WebSocket/HTTP 升级入口
- Session 录制/回放
- 命令白名单（可在 execHandler 回调中由使用者自行实现）

## Not Doing (and Why)
- **Windows PTY** — Windows 没有原生 PTX，需要 winpty 等额外依赖，复杂性远高于收益。
     如果确实需要，使用者可以自行实现 shellHandler。
- **SFTP** — 这不是管理通道的核心需求。文件传输可以通过 exec（scp 模拟）或专门的子协议做。
- **插件化接口 v1** — 先写完整功能，让实现驱动接口设计，而不是反过来。
     避免"第一次就做对"的过度设计陷阱。
- **客户端 API** — 问题范围限定在 sshd 端。客户端做不做是另一个独立需求。

## Open Questions
- PTY 分配使用 `github.com/creack/pty` 还是直接操作 `/dev/ptmx` + `syscall`？
     creack/pty 封装了跨平台逻辑，减少 ~200 行平台特定代码。但它是一个新依赖。
- 转发端口的生命周期：当 SSH 连接断开时，需要关闭所有已绑定的转发端口。
     如何做到可靠清理？通过 per-connection context 派生 + waitGroup 等待所有转发 goroutine 退出。
- 默认 handler 的 fallback 行为：自定义 handler 未设置时，默认提供完整实现？
     还是默认拒绝所有请求（更安全）？倾向默认提供完整实现，符合直觉。
