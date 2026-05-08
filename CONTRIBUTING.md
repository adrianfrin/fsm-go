# 贡献指南

## 本地检查

提交前请至少运行：

```bash
go test ./...
go test -race ./...
```

如果本机有 Taskfile 和 golangci-lint，运行：

```bash
task check
```

涉及 MySQL Repository、事务、并发、幂等或 Outbox 的改动，需要运行集成测试：

```bash
go test -tags=integration ./test/integration/...
```

## 代码要求

- 状态变化必须走状态机入口。
- 新增 DSL 能力必须补校验。
- 新增存储行为必须补 Testcontainers 集成测试。
- 不要把业务租户字段写死进核心表结构。
- 不要把核心库绑定成必须部署的服务。
