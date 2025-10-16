# AutoTest 缓存系统指南

## 概述

AutoTest的缓存系统旨在提升测试执行性能，减少重复计算和网络请求。通过智能缓存，可以显著提高大规模测试场景的执行效率。

## 缓存类型

### 1. gRPC描述符缓存 (GrpcDescriptorCache)

**用途：** 缓存gRPC服务的反射信息和方法描述符

**场景：**
- 多个测试用例调用同一个gRPC服务
- 避免重复的服务发现和反射调用
- 提升gRPC测试执行速度

**示例：**
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

### 2. HTTP响应缓存 (HttpResponseCache)

**用途：** 缓存GET请求的响应结果

**场景：**
- 获取测试数据或配置信息
- 减少对第三方API的重复调用
- 提升依赖数据获取效率

**示例：**
```yaml
# 获取用户信息，结果会被缓存
- id: 1
  request:
    method: "get"
    url: "http://{{ HOST }}/api/users/123"

# 后续相同请求会使用缓存
- id: 2
  request:
    method: "get"  
    url: "http://{{ HOST }}/api/users/123"  # 从缓存获取
```

### 3. 模板缓存 (TemplateCache)

**用途：** 缓存模板渲染结果

**场景：**
- 相同模板和变量组合的重复渲染
- URL模板、请求体模板等
- 减少模板解析开销

**示例：**
```yaml
# 相同的URL模板会被缓存
- id: 1
  request:
    url: "http://{{ HOST }}/api/books/{{ BOOK_ID }}"  # 首次渲染

- id: 2
  request:
    url: "http://{{ HOST }}/api/books/{{ BOOK_ID }}"  # 使用缓存结果
```

### 4. Lua脚本缓存 (LuaScriptCache)

**用途：** 缓存编译后的Lua脚本

**场景：**
- 相同Lua脚本的重复编译
- 提升脚本执行性能
- 减少Lua虚拟机开销

**示例：**
```yaml
# 相同的Lua脚本会被缓存
- id: 1
  request:
    luaBody: |
      function body()
        return '{"id": 1}'
      end

- id: 2  
  request:
    luaBody: |
      function body()
        return '{"id": 1}'  # 相同脚本，使用缓存的编译结果
      end
```

## 配置选项

```yaml
global:
  cache:
    enabled: true        # 是否启用缓存
    ttl: 300s           # 缓存过期时间（5分钟）
    max_size: 1000      # 最大缓存条目数
```

### 配置说明

- **enabled**: 控制是否启用缓存功能
- **ttl**: 缓存项的生存时间，过期后自动清理
- **max_size**: 缓存的最大条目数，超出时会淘汰最旧的项

## 性能优化建议

### 1. 合理设置TTL

```yaml
# 根据数据更新频率设置TTL
cache:
  ttl: 300s    # 对于相对稳定的数据，可以设置较长的TTL
```

### 2. 控制缓存大小

```yaml
# 根据内存使用情况调整缓存大小
cache:
  max_size: 2000  # 增加缓存大小以提高命中率
```

### 3. 缓存预热

在测试开始前预加载常用数据：

```bash
# 可以通过配置文件预定义常用的模板和数据
environments:
  dev:
    HOST: "localhost:8080"
    GRPC_SERVER: "localhost:50051"
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
# 在日志中查看缓存统计信息
grep "Cache" /var/log/test/autotest.log

# 示例输出：
# Cache HIT: grpc:localhost:50051
# Cache MISS: http:get:localhost:8080/api/new-endpoint
# Cache cleanup: removed 5 expired items
```

## 最佳实践

### 1. 适合缓存的场景

✅ **推荐缓存：**
- GET请求获取参考数据
- gRPC服务反射信息
- 相同的模板渲染
- 重复的Lua脚本

❌ **不推荐缓存：**
- POST/PUT/DELETE等修改操作
- 包含随机数据的请求
- 实时性要求高的数据

### 2. 缓存键设计

缓存系统会自动生成缓存键，包含：
- 请求方法和URL
- 重要的请求头（如Authorization）
- 模板内容和变量
- Lua脚本内容

### 3. 内存管理

```yaml
# 根据测试规模调整缓存配置
cache:
  enabled: true
  ttl: 600s        # 大型测试可以设置更长的TTL
  max_size: 5000   # 增加缓存容量
```

## 故障排除

### 缓存未生效

1. **检查配置**：确认 `cache.enabled: true`
2. **查看日志**：搜索 "Cache" 关键字
3. **验证请求**：确保是相同的GET请求

### 内存使用过高

1. **减少缓存大小**：降低 `max_size` 值
2. **缩短TTL**：减少 `ttl` 时间
3. **禁用缓存**：设置 `enabled: false`

### 数据不一致

1. **清理缓存**：重启测试程序
2. **缩短TTL**：减少缓存时间
3. **排除缓存**：临时禁用缓存验证

## 示例配置

### 高性能配置
```yaml
global:
  cache:
    enabled: true
    ttl: 900s      # 15分钟
    max_size: 3000
```

### 保守配置  
```yaml
global:
  cache:
    enabled: true
    ttl: 60s       # 1分钟
    max_size: 500
```

### 禁用缓存
```yaml
global:
  cache:
    enabled: false
```

通过合理配置和使用缓存系统，可以显著提升AutoTest的执行性能，特别是在大规模测试场景中。
