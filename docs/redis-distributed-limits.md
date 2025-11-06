# 分布式 Rate Limiting 和 Retry 管理

## 概述

当您的应用程序启动多个进程时，使用本地的 rate limiter 和 retry handler 会导致每个进程独立管理限制，这可能导致：
- 超过实际的 API rate limit（多个进程的总请求数超过限制）
- 重复的重试操作（多个进程同时重试相同的失败操作）

使用 Redis 实现分布式的 rate limiting 和 retry 管理可以解决这些问题。

## 配置

### 基本配置

在配置文件中添加 Redis 设置：

```yaml
# Redis 配置
redis_address: "localhost:6379"
redis_password: ""  # 可选
redis_db: 0
redis_key_prefix: "amapi:"

# 启用分布式功能
use_redis_rate_limit: true
use_redis_retry: true
```

### 环境变量配置

```bash
export AMAPI_REDIS_ADDRESS="localhost:6379"
export AMAPI_REDIS_PASSWORD=""
export AMAPI_REDIS_DB=0
export AMAPI_REDIS_KEY_PREFIX="amapi:"
export AMAPI_USE_REDIS_RATE_LIMIT=true
export AMAPI_USE_REDIS_RETRY=true
```

### 代码配置

```go
cfg := &config.Config{
    ProjectID:       "your-project-id",
    CredentialsFile: "./sa-key.json",

    // Redis 配置
    RedisAddress:      "localhost:6379",
    RedisPassword:     "",
    RedisDB:           0,
    RedisKeyPrefix:    "amapi:",
    UseRedisRateLimit: true,
    UseRedisRetry:     true,

    // Rate limiting 配置（会通过 Redis 在多个进程间共享）
    RateLimit: 100,  // 每分钟 100 个请求（所有进程总计）
    RateBurst: 20,

    // Retry 配置
    RetryAttempts: 3,
    RetryDelay:    1 * time.Second,
    EnableRetry:   true,
}

client, err := client.New(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## 工作原理

### 分布式 Rate Limiting

使用 Redis 的 **滑动窗口计数器** 算法实现：

1. 每个请求在 Redis 中记录一个带时间戳的条目
2. 定期清理超出时间窗口的旧条目
3. 统计当前时间窗口内的请求数
4. 如果超过限制，等待直到有足够的配额

**Redis Key 结构：**
```
amapi:ratelimit:requests  (Sorted Set)
```

每个条目：
- Score: 请求的时间戳（秒）
- Member: 时间戳字符串

**特点：**
- ✅ 所有进程共享同一个 rate limit
- ✅ 精确控制每分钟的请求总数
- ✅ 自动清理过期数据
- ✅ 支持突发流量（burst）

### 分布式 Retry 管理

使用 Redis **分布式锁** 防止多个进程同时重试同一操作：

1. 每个重试操作生成唯一的 operation ID
2. 尝试获取 Redis 锁
3. 如果获取成功，执行重试
4. 如果获取失败，等待一小段时间后检查操作是否已成功
5. 操作完成后释放锁

**Redis Key 结构：**
```
amapi:retry:lock:{operationID}     (String, TTL: 1分钟)
amapi:retry:count:{operationID}    (Integer, TTL: 1小时)
```

**特点：**
- ✅ 防止多个进程同时重试同一操作
- ✅ 减少重复的 API 调用
- ✅ 提供重试次数统计
- ✅ 自动过期清理

## 使用示例

### 示例 1：基本使用

```go
package main

import (
    "log"
    "time"

    "amapi-pkg/pkgs/amapi/client"
    "amapi-pkg/pkgs/amapi/config"
)

func main() {
    cfg := config.DefaultConfig()
    cfg.ProjectID = "your-project-id"
    cfg.CredentialsFile = "./sa-key.json"

    // 启用 Redis 分布式管理
    cfg.RedisAddress = "localhost:6379"
    cfg.UseRedisRateLimit = true
    cfg.UseRedisRetry = true

    // 全局 rate limit: 100 请求/分钟（所有进程共享）
    cfg.RateLimit = 100

    c, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // 启动多个 goroutine 模拟多进程
    for i := 0; i < 10; i++ {
        go func(id int) {
            for j := 0; j < 20; j++ {
                // 所有 goroutine 共享同一个 rate limit
                _, err := c.Enterprises().List(nil)
                if err != nil {
                    log.Printf("Process %d: Error: %v", id, err)
                } else {
                    log.Printf("Process %d: Request %d succeeded", id, j)
                }
                time.Sleep(100 * time.Millisecond)
            }
        }(i)
    }

    // 等待所有 goroutine 完成
    time.Sleep(30 * time.Second)
}
```

### 示例 2：监控 Rate Limit

```go
// 检查当前时间窗口内的请求数
func checkRateLimit(client *redis.Client, prefix string) (int64, error) {
    key := prefix + "ratelimit:requests"
    ctx := context.Background()

    // 获取当前请求数
    count, err := client.ZCard(ctx, key).Result()
    if err != nil {
        return 0, err
    }

    return count, nil
}
```

### 示例 3：监控 Retry 次数

```go
// 使用 RedisRetryHandler 的 GetRetryCount 方法
retryHandler := utils.NewRedisRetryHandler(redisClient, "amapi:", config)
count, err := retryHandler.GetRetryCount(ctx, "operation-123")
if err != nil {
    log.Printf("Failed to get retry count: %v", err)
} else {
    log.Printf("Operation retried %d times", count)
}
```

## 性能考虑

### Rate Limiting

- **Redis 操作**：每次请求需要 2-3 个 Redis 操作（ZRemRangeByScore, ZCard, ZAdd）
- **延迟**：如果 Redis 在本地网络，延迟通常 < 1ms
- **扩展性**：可以处理每秒数千次请求

### Retry Management

- **Redis 操作**：每次重试需要 2-3 个 Redis 操作（SetNX, Del, Incr）
- **锁竞争**：如果多个进程同时重试，只有一个会成功，其他会等待
- **性能影响**：正常情况下影响很小，只有在高并发重试时才会有锁竞争

## 故障处理

### Redis 连接失败

如果 Redis 连接失败，客户端初始化会返回错误：
```go
client, err := client.New(cfg)
// err: "failed to connect to Redis: ..."
```

### Redis 运行时故障

如果 Redis 在运行时故障：
- **Rate Limiting**：会返回错误，请求会被拒绝
- **Retry Management**：会回退到本地重试（如果获取锁失败）

建议：
1. 使用 Redis Sentinel 或 Cluster 实现高可用
2. 监控 Redis 连接状态
3. 配置连接池和重连机制

## 最佳实践

1. **使用连接池**：Redis 客户端默认使用连接池，无需额外配置
2. **设置合理的 Key 前缀**：避免不同环境或项目之间的 key 冲突
3. **监控 Redis 内存使用**：定期清理过期的 key
4. **使用 Redis Sentinel/Cluster**：在生产环境中实现高可用
5. **设置合理的 TTL**：确保过期的 key 能够自动清理

## 与本地实现的对比

| 特性 | 本地实现 | Redis 实现 |
|------|----------|------------|
| 多进程支持 | ❌ 每个进程独立限制 | ✅ 所有进程共享限制 |
| 精确控制 | ⚠️ 可能超过总限制 | ✅ 精确控制总请求数 |
| 重试协调 | ❌ 可能重复重试 | ✅ 防止重复重试 |
| 性能 | ✅ 零延迟 | ⚠️ 需要 Redis 网络调用 |
| 复杂度 | ✅ 简单 | ⚠️ 需要 Redis 服务器 |
| 监控 | ❌ 无法跨进程监控 | ✅ 可以在 Redis 中查看统计 |

## 优先级队列模式

除了基本的分布式 rate limiting 和 retry 管理，还可以使用**优先级队列模式**来管理 API 调用。

### 概述

优先级队列模式将 API 调用封装为任务，按优先级放入 Redis sorted set，由 worker 异步消费执行。执行时会应用 rate limiting（如 1000 次/秒）和 retry 逻辑（针对 429 错误）。

**主要特性：**
- ✅ 按优先级执行任务（数字，0-1000，越大优先级越高）
- ✅ 异步任务执行，提高吞吐量
- ✅ 自动 rate limiting（在 worker 执行时应用）
- ✅ 自动处理 429 错误并重试
- ✅ 支持任务优先级调整（重试时降低优先级）

### 配置

```yaml
# 启用优先级队列模式
use_priority_queue: true

# Worker 配置
queue_worker_concurrency: 10      # Worker 并发数
queue_poll_interval: 100ms         # 队列轮询间隔
default_task_priority: 500         # 默认任务优先级（0-1000）
max_queue_size: 10000              # 最大队列大小

# 优先级队列的 rate limit
priority_queue_rate_limit: 1000    # 请求/秒（如果未设置，使用 rate_limit/60）
priority_queue_burst: 100          # Burst 容量
```

### 使用示例

```go
cfg := &config.Config{
    ProjectID:       "your-project-id",
    CredentialsFile: "./sa-key.json",

    // Redis 配置
    RedisAddress: "localhost:6379",

    // 启用优先级队列模式
    UsePriorityQueue: true,
    QueueWorkerConcurrency: 10,
    DefaultTaskPriority: 500,
    PriorityQueueRateLimit: 1000,  // 1000 请求/秒
}

client, err := client.New(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 使用默认优先级执行 API 调用
_, err = client.Enterprises().List(nil)

// 使用指定优先级执行（高优先级）
// 注意：需要 Client 支持优先级参数
```

### 工作原理

1. **任务创建**：API 调用被封装为 Task，包含操作信息和优先级
2. **入队**：Task 按优先级放入 Redis sorted set（score = priority）
3. **Worker 消费**：多个 worker goroutine 并发从队列中按优先级取出任务
4. **Rate Limiting**：Worker 执行任务前检查 rate limit（如 1000 次/秒），超过限制时等待
5. **执行任务**：Worker 执行任务（调用 API）
6. **处理 429**：如果返回 429 错误：
   - 计算重试延迟（指数退避）
   - 降低优先级（原优先级 - 50）重新入队
   - 等待延迟后重试
7. **存储结果**：任务执行结果存储在 Redis 中

**Redis Key 结构：**
```
amapi:queue:priority            (Sorted Set)
  - Score: 优先级（0-1000）
  - Member: Task JSON 字符串

amapi:task:result:{callbackID}  (Hash)
  - status: pending/processing/completed/failed
  - result: 结果 JSON（如果成功）
  - error: 错误信息（如果失败）
  - created_at: 创建时间
  - completed_at: 完成时间
```

### 优先级策略

- **默认优先级**：500
- **高优先级任务**：700-1000（重要操作）
- **普通任务**：400-600（常规操作）
- **低优先级任务**：0-300（后台任务）
- **429 重试任务**：原优先级 - 50（每次重试降低）

### Rate Limiting

优先级队列模式使用**滑动窗口算法**，支持每秒限制（如 1000 次/秒）：

- 使用 1 秒窗口
- Worker 执行任务前检查当前窗口内的请求数
- 超过限制时等待，直到有配额
- 所有 worker 共享同一个 rate limit

### Retry 策略（429 错误）

当请求返回 429（Too Many Requests）时：

1. 检测 HTTP 429 响应
2. 如果重试次数 < MaxRetries：
   - 计算延迟（指数退避：baseDelay * 2^attempt）
   - 降低优先级（原优先级 - 50）重新入队
   - 或等待延迟后直接重试
3. 使用 Redis 锁防止重复重试
4. 如果重试耗尽，标记任务失败

### Worker 并发控制

- 使用 channel 限制并发数（semaphore pattern）
- 支持优雅关闭（context cancellation）
- 任务超时处理（context timeout）

### 与标准模式的对比

| 特性 | 标准模式 | 优先级队列模式 |
|------|----------|----------------|
| 执行方式 | 同步执行 | 异步执行（通过队列） |
| 优先级支持 | ❌ | ✅ 支持任务优先级 |
| Rate Limiting | 在入队前检查 | 在执行时检查（worker） |
| 429 重试 | 立即重试 | 降低优先级重新入队 |
| 吞吐量 | 受限于同步等待 | 更高（异步执行） |
| 复杂度 | ✅ 简单 | ⚠️ 需要 worker 管理 |
| 适用场景 | 简单场景 | 高并发、需要优先级的场景 |

## 总结

使用 Redis 实现分布式 rate limiting 和 retry 管理是在多进程环境中保证 API 调用合规性的最佳实践。通过简单的配置就可以启用这些功能，而无需修改业务代码。

对于需要更高吞吐量和优先级管理的场景，可以使用优先级队列模式，它提供了更灵活的任务调度和更好的并发处理能力。
