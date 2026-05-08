# 测试说明

本项目要求核心逻辑和真实依赖都要覆盖。

## 单元测试

```bash
go test ./...
```

覆盖内容：

- DSL 校验。
- Guard 判断。
- 状态迁移。
- 幂等命中。
- Action 执行。
- 状态日志和 Outbox 写入。

## 竞态测试

```bash
go test -race ./...
```

用于检查状态机运行时和测试辅助实现是否存在数据竞争。

## 集成测试

```bash
go test -tags=integration ./test/integration/...
```

集成测试使用 Testcontainers 启动真实 MySQL，覆盖：

- MySQL schema 初始化。
- CAS 状态更新。
- 并发请求只有一个成功。
- 状态日志写入。
- 幂等结果复用。
- Outbox 事务写入。

## 本地聚合检查

如果安装了 Taskfile：

```bash
task check
```

该命令会运行格式检查、`go vet`、`golangci-lint`、单元测试和竞态测试。
