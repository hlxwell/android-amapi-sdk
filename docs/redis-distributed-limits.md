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

## 总结

使用 Redis 实现分布式 rate limiting 和 retry 管理是在多进程环境中保证 API 调用合规性的最佳实践。通过简单的配置就可以启用这些功能，而无需修改业务代码。
