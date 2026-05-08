# Benchmark

本项目提供 10 万级状态流转 Benchmark，用来观察核心运行时在批量数据下的耗时、吞吐和内存分配。

Benchmark 覆盖两组场景：

- `without_observability`：只运行核心状态流转。
- `with_prometheus_observability`：开启 Prometheus 可观测性后运行同样的状态流转。

## 运行命令

```bash
go test -run '^$' -bench BenchmarkRuntimeFire100K -benchtime=1x -benchmem ./test/benchmark
```

如果使用 Taskfile：

```bash
task test:benchmark
```

或者不全局安装 Taskfile：

```bash
go run github.com/go-task/task/v3/cmd/task@v3.50.0 test:benchmark
```

## 测试内容

单次 Benchmark 会预先创建 100,000 条状态实体，然后对这 100,000 条实体逐条触发一次状态流转。

每次流转都会经过：

- DSL 编译后的迁移规则匹配。
- Guard 表达式判断。
- CAS 状态更新。
- 状态日志写入。
- 幂等结果保存。

开启可观测性时，还会额外记录：

- 流转总次数。
- 流转耗时。
- 错误次数。
- 幂等命中次数。
- 正在执行的流转数量。

## 本机结果示例

以下结果来自本机 Apple M1 Pro，单次运行仅作为参考，真实结果会随机器、Go 版本和运行负载变化。

| 场景 | 总耗时 | 单次流转耗时 | 吞吐 | 内存分配 |
|---|---:|---:|---:|---:|
| 不开启可观测性 | 561.80 ms | 5,618 ns | 178,005 transitions/s | 699.27 MB |
| 开启 Prometheus 可观测性 | 601.15 ms | 6,011 ns | 166,354 transitions/s | 699.28 MB |

这次结果里，开启 Prometheus 可观测性后，单次流转耗时约增加 7.0%。这个开销主要来自指标计数、标签写入和耗时统计。

## 使用建议

Benchmark 用来观察趋势，不应该作为固定性能承诺。实际业务接入时，建议结合自己的 DSL 复杂度、Action 数量、存储实现、数据库延迟和并发模型重新运行。
