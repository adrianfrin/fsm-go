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
go test -count=1 -tags=integration ./test/integration/...
```

集成测试使用 Testcontainers 启动真实 MySQL，覆盖：

- MySQL schema 初始化。
- CAS 状态更新。
- 并发请求只有一个成功。
- 状态日志写入。
- 幂等结果复用。
- Outbox 事务写入。

## Benchmark

```bash
go test -run '^$' -bench BenchmarkRuntimeFire100K -benchtime=1x -benchmem ./test/benchmark
```

Benchmark 会执行 100,000 次状态流转，并分别对比不开启可观测性和开启 Prometheus 可观测性的耗时与内存分配。

更多说明见 [Benchmark](benchmark.md)。

## 本地聚合检查

如果安装了 Taskfile：

```bash
task check
```

该命令会运行格式检查、依赖整理检查、`go vet`、`golangci-lint`、单元测试、竞态测试、Testcontainers 集成测试和 10 万级 Benchmark。

如果不想全局安装 Taskfile，可以直接运行：

```bash
go run github.com/go-task/task/v3/cmd/task@v3.50.0 check
```
