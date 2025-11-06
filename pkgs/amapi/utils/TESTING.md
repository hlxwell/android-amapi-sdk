# Utils 测试文档

## 概述

`utils` 包提供了优先级队列、rate limiting 和 retry 功能的集成测试。测试使用 `miniredis` 作为内存 Redis 服务器，无需外部依赖。

## 运行测试

### 运行所有测试

```bash
cd pkgs/amapi
go test ./utils -v -timeout 120s
```

### 运行特定测试

```bash
# 测试任务序列化
go test ./utils -v -run TestTaskSerialization

# 测试 Redis 优先级队列
go test ./utils -v -run TestRedisPriorityQueue

# 测试 Task Worker 集成
go test ./utils -v -run TestTaskWorkerIntegration

# 测试 429 错误重试
go test ./utils -v -run TestTaskWorker429Retry
```

## 测试覆盖

### 1. Task 序列化和反序列化 (`TestTaskSerialization`)

- ✅ Task 创建和序列化
- ✅ Task 反序列化
- ✅ 优先级验证（0-1000）
- ✅ 任务 ID 和 CallbackID 生成

### 2. Redis Priority Queue (`TestRedisPriorityQueue`)

- ✅ 任务入队
- ✅ 按优先级出队（高优先级优先）
- ✅ 队列大小查询
- ✅ 空队列处理

### 3. Task Worker 集成 (`TestTaskWorkerIntegration`)

- ✅ Worker 启动和停止
- ✅ 任务执行流程
- ✅ 任务状态更新（pending -> processing -> completed）
- ✅ 任务结果存储和查询
- ✅ Rate limiting 应用

### 4. 429 错误重试 (`TestTaskWorker429Retry`)

- ✅ 429 错误检测
- ✅ 自动重试机制
- ✅ 优先级降低（每次重试 -50）
- ✅ 指数退避延迟
- ✅ 重试次数限制

### 5. Priority Queue Rate Limiter (`TestPriorityQueueRateLimiter`)

- ✅ Rate limiter 接口实现
- ✅ 队列大小检查
- ✅ Allow() 方法

### 6. Priority Queue Retry Handler (`TestPriorityQueueRetryHandler`)

- ✅ Retry handler 接口实现
- ✅ 操作执行和重试协调

### 7. 端到端流程 (`TestEndToEndFlow`)

- ✅ 基本流程（任务创建 -> 入队 -> 执行 -> 完成）
- ✅ 优先级排序（多个任务按优先级执行）
- ✅ Rate limiting 验证

### 8. 错误处理 (`TestErrorHandling`)

- ✅ 永久错误处理
- ✅ 任务失败状态
- ✅ 错误信息存储

## Mock 实现

### Google API Mock

测试使用 `mockAPICallExecutor` 函数来模拟 Google API 调用：

```go
func mockAPICallExecutor(success bool, errCode int, delay time.Duration) TaskExecutor
```

参数：
- `success`: 是否成功
- `errCode`: 错误代码（如 `types.ErrCodeTooManyRequests`）
- `delay`: 模拟 API 调用延迟

### Rate Limiting Mock

测试使用真实的 Redis rate limiter（通过 miniredis），但配置了较高的 rate limit（1000 次/秒）以避免测试中的限制。

## 测试配置

### Redis 设置

测试使用 `miniredis` 作为内存 Redis 服务器：
- 无需外部 Redis 服务器
- 自动清理测试数据
- 支持所有 Redis 命令

### Worker 配置

测试中的 Worker 配置：
- `Concurrency`: 2-3 个 worker
- `PollInterval`: 200ms（用于快速测试）
- `RateLimit`: 1000 次/秒（高限制以避免测试阻塞）
- `BaseDelay`: 500ms（用于重试测试）

### 超时设置

- 任务执行超时：10-30 秒
- 等待结果超时：10-20 秒
- 测试上下文超时：30 秒

## 已知限制

1. **miniredis 限制**：
   - 不支持小于 1 秒的过期时间（会警告但不会失败）
   - 某些 Redis 命令可能不完全支持

2. **测试时间**：
   - 由于需要等待 worker 处理和重试，某些测试可能需要较长时间
   - 建议使用 `-timeout` 参数增加超时时间

3. **并发测试**：
   - 某些测试可能因为并发执行而需要更长的等待时间

## 改进建议

1. 添加更多边界情况测试
2. 添加性能测试（压力测试）
3. 添加并发安全测试
4. 添加更详细的错误场景测试
5. 添加 metrics 和监控测试

## 故障排查

### 测试超时

如果测试超时，可以：
1. 增加超时时间：`-timeout 120s`
2. 检查 worker 配置（PollInterval, Concurrency）
3. 检查 rate limit 设置（测试中使用高限制）

### Redis 连接失败

如果 Redis 连接失败：
1. 确保 `miniredis` 已正确安装：`go get github.com/alicebob/miniredis/v2`
2. 检查测试代码中的 Redis 初始化

### 任务未完成

如果任务一直处于 pending 状态：
1. 检查 worker 是否已启动
2. 增加等待时间
3. 检查 executor 是否已注册
4. 检查 rate limiter 是否阻塞


