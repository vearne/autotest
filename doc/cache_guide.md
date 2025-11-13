# AutoTest 缓存系统指南

## 概述

AutoTest的缓存系统专门优化gRPC测试性能，通过缓存gRPC服务描述符避免重复的反射调用，显著提升大规模gRPC测试的执行效率。

> **重要**：缓存系统**默认开启**！考虑到gRPC描述符缓存收益明显（可节省50-200秒）且开销很小，系统会自动启用缓存功能。

> **注意**：这是经过优化的缓存系统，专注于最有价值的gRPC描述符缓存，移除了收益有限的HTTP响应缓存、模板缓存和Lua脚本缓存，既保留了核心性能优化，又降低了系统复杂性。

## gRPC描述符缓存

### 用途
缓存gRPC服务的反射信息和方法描述符，避免重复的服务发现和反射调用。

### 适用场景
- 多个测试用例调用同一个gRPC服务
- 大规模gRPC测试（1000+测试用例）
- 需要优化测试执行时间的场景

### 性能收益
在1000个gRPC测试用例的场景下：
- **无缓存**: 1000次 × 50-200ms (反射调用) = 50-200秒
- **有缓存**: 1次 × 100ms + 999次 × 0.001ms ≈ 0.1秒
- **节省时间**: 约50-200秒

### 示例配置

```yaml
# 多个测试用例使用同一个gRPC服务
- id: 1
  request:
    address: "{{ GRPC_SERVER }}"  # 第一次会获取服务描述符
    symbol: "Bookstore/GetBook"

- id: 2  
  request:
    address: "{{ GRPC_SERVER }}"  # 从缓存获取，无需重复反射
    symbol: "Bookstore/ListBook"
```

## 配置选项

### 默认行为

**缓存默认开启**，无需任何配置即可享受性能优化！

### 可选配置

如果需要自定义缓存行为，可以在配置文件中（如 `config_files/autotest.yml` 或 `config_files/autotest_enhanced.yml`）添加：

```yaml
global:
  cache:
    enabled: true        # 可选：显式开启缓存（默认已开启）
    ttl: 300s           # 可选：缓存过期时间（默认5分钟）
    max_size: 100       # 可选：最大缓存条目数（默认100）
```

### 禁用缓存

如果您确实需要禁用缓存（不推荐），可以这样配置：

```yaml
global:
  cache:
    enabled: false      # 显式禁用缓存
```

### 配置说明

- **enabled**: 控制是否启用缓存功能（**默认：true**）
- **ttl**: 缓存项的生存时间，过期后自动清理（**默认：5分钟**）  
- **max_size**: 缓存的最大条目数（**默认：100个**，足够gRPC描述符使用）

## 性能优化建议

### 1. 合理设置TTL

```yaml
# 对于相对稳定的gRPC服务，可以设置较长的TTL
cache:
  ttl: 600s    # 10分钟，适合稳定的开发/测试环境
```

### 2. 控制缓存大小

```yaml
# 根据使用的gRPC服务数量调整
cache:
  max_size: 50   # 如果只有几个gRPC服务，可以设置更小的值
```

## 缓存统计

### 查看缓存效果

缓存系统会自动收集统计信息：

- **命中率 (Hit Rate)**: 缓存命中的百分比
- **命中次数 (Hits)**: 成功从缓存获取数据的次数  
- **未命中次数 (Misses)**: 缓存未命中的次数
- **缓存大小 (Size)**: 当前缓存的条目数

### 监控缓存性能

```bash
# 运行gRPC测试并查看缓存效果
./autotest grpc-automate -c config_files/autotest.yml

# 在日志中查看缓存统计信息
grep "Cache" /var/log/test/autotest.log

# 示例输出：
# Cache initialized: TTL=5m0s, MaxSize=100
# Using cached gRPC descriptor for localhost:50051
# Cache HIT: localhost:50051
# Cache [grpc_descriptor]: Hits=999, Misses=1, HitRate=99.90%, Size=1
# gRPC descriptor cache cleared
```

### 验证缓存效果

```bash
# 第一次运行：缓存未命中，会较慢
time ./autotest grpc-automate -c config_files/autotest.yml

# 第二次运行：缓存命中，会明显加速
time ./autotest grpc-automate -c config_files/autotest.yml

# 比较两次执行时间，第二次应该明显更快
```

## 故障排除

### 缓存未生效

1. **检查配置**：确认 `cache.enabled: true`
2. **查看日志**：搜索 "Cache" 关键字
3. **验证服务地址**：确保gRPC服务地址一致

### 内存使用过高

1. **减少缓存大小**：降低 `max_size` 值到 50 或更少
2. **缩短TTL**：减少 `ttl` 时间到 300s 或更短
3. **禁用缓存**：设置 `enabled: false`

### 服务更新后缓存不一致

1. **重启测试程序**：清理所有缓存
2. **缩短TTL**：减少缓存时间，如设置为 60s
3. **手动清理**：程序会在TTL到期后自动清理

## 示例配置

### 高性能配置
```yaml
global:
  cache:
    enabled: true
    ttl: 900s      # 15分钟
    max_size: 100
```

### 保守配置  
```yaml
global:
  cache:
    enabled: true
    ttl: 300s      # 5分钟
    max_size: 50
```

### 禁用缓存
```yaml
global:
  cache:
    enabled: false
```

## 总结

对于1000个测试用例规模的gRPC测试：
- **性能收益**：gRPC描述符缓存可以节省50-200秒的执行时间
- **资源占用**：内存开销很小（通常只需要缓存几个gRPC服务的描述符）
- **零配置**：默认开启，无需任何配置即可享受性能优化
- **适用场景**：特别适合CI/CD环境中的大规模自动化测试

### 快速开始

1. **直接运行**：使用 `./autotest grpc-automate -c your-config.yml`（缓存已默认开启）
2. **查看效果**：检查日志中的缓存统计信息
3. **可选调优**：根据实际使用情况调整 `ttl` 和 `max_size`

> 💡 **提示**：由于缓存默认开启，您无需任何额外配置即可享受性能提升！

通过合理配置gRPC描述符缓存，可以显著提升AutoTest的执行性能，让您的测试更快、更高效！