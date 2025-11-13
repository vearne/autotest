# AutoTest 使用指南

## 目录
- [快速开始](#快速开始)
- [配置文件详解](#配置文件详解)
- [测试用例编写](#测试用例编写)
- [高级功能](#高级功能)
- [最佳实践](#最佳实践)
- [故障排除](#故障排除)
- [企业级使用场景](#企业级使用场景)

## 快速开始

### 1. 安装
```bash
# 方式1：下载预编译二进制文件
wget https://github.com/vearne/autotest/releases/latest/download/autotest-linux-amd64
chmod +x autotest-linux-amd64
mv autotest-linux-amd64 /usr/local/bin/autotest

# 方式2：从源码编译
git clone https://github.com/vearne/autotest.git
cd autotest
make build
```

### 2. 验证配置
```bash
autotest test --config-file=./config_files/autotest.yml
```

### 3. 运行测试
```bash
# 基础用法：使用环境文件
autotest run --config-file=./config_files/autotest.yml --env-file=./config_files/.env.dev

# 高级用法：环境选择（推荐）
autotest run --config-file=./config_files/autotest.yml --environment=dev
autotest run --config-file=./config_files/autotest.yml --environment=staging  
autotest run --config-file=./config_files/autotest.yml --environment=prod
```

## 配置文件详解

### 全局配置 (autotest.yml)
```yaml
global:
  worker_num: 5                    # 并发工作协程数
  ignore_testcase_fail: true       # 是否忽略测试用例失败
  debug: false                     # 调试模式
  request_timeout: 5s              # 请求超时时间
  
  # 重试机制：提升测试稳定性（自动启用，可选配置）
  retry:
    max_attempts: 3                # 最大重试次数（默认3次）
    retry_delay: 1s                # 重试间隔（默认1秒）
    retry_on_status_codes: [500, 502, 503, 504, 408, 429]  # HTTP重试状态码
  
  # 并发控制：性能和稳定性平衡（可选，有合理默认值）
  concurrency:
    max_concurrent_requests: 20    # 最大并发请求数（默认20）
    rate_limit_per_second: 50      # 每秒请求数限制（默认50）
  
  # 智能缓存：性能优化（默认开启）
  cache:
    enabled: true                  # gRPC描述符缓存（默认开启）
    ttl: 300s                     # 缓存过期时间（默认5分钟）
    max_size: 100                 # 缓存条目数（默认100个）
  
  logger:
    level: "debug"                 # 日志级别: debug/info/error
    file_path: "/var/log/test/autotest.log"
    
  # 多格式报告生成
  report:
    dir_path: "/var/log/test/report/"        # 报告输出目录
    formats: ["html", "json", "csv", "junit"]  # 报告格式（可选）
    template_path: "./templates/custom.html"    # 自定义模板（可选）
  
  # Slack 通知集成
  notifications:
    enabled: true                            # 是否启用通知
    webhook_url: "https://hooks.slack.com/services/..."  # Slack Webhook URL
    on_failure: true                         # 失败时通知
    on_success: false                        # 成功时通知（可选）

# 环境管理：支持多环境切换
environments:
  dev:
    HOST: "localhost:8080"
    GRPC_SERVER: "localhost:50031"
    API_KEY: "dev_api_key"
  staging:
    HOST: "staging.api.com"
    GRPC_SERVER: "staging.grpc.com:50031"
    API_KEY: "staging_api_key"
  prod:
    HOST: "api.example.com"
    GRPC_SERVER: "grpc.example.com:50031"
    API_KEY: "prod_api_key"

http_rule_files:                   # HTTP测试用例文件列表
  - "./config_files/my_http_api.yml"

grpc_rule_files:                   # gRPC测试用例文件列表
  - "./config_files/my_grpc_api.yml"
```

### 环境变量文件 (.env.dev)
```bash
HOST=localhost:8080
GRPC_SERVER=localhost:50031
API_KEY=your_api_key_here
```

## 测试用例编写

### HTTP 测试用例
```yaml
- id: 1
  desc: "创建新书籍"
  request:
    method: "post"
    url: "http://{{ HOST }}/api/books"
    headers:
      - "Content-Type: application/json"
      - "Authorization: Bearer {{ API_KEY }}"
    body: |
      {
        "title": "Go语言编程",
        "author": "张三"
      }
  rules:
    - name: "HttpStatusEqualRule"
      expected: 200
    - name: "HttpBodyEqualRule"
      xpath: "/title"
      expected: "Go语言编程"
  export:
    xpath: "/id"
    exportTo: "BOOK_ID"
    type: integer
```

### gRPC 测试用例
```yaml
- id: 1
  desc: "获取书籍信息"
  request:
    address: "{{ GRPC_SERVER }}"
    symbol: "Bookstore/GetBook"
    body: |
      {
        "id": 1
      }
  rules:
    - name: "GrpcCodeEqualRule"
      expected: "OK"
    - name: "GrpcBodyEqualRule"
      xpath: "/data/title"
      expected: "Go语言编程"
```

## 高级功能

### 1. 环境管理
```bash
# 选择不同环境运行测试
autotest run -c config.yml --environment=dev      # 开发环境
autotest run -c config.yml --environment=staging  # 测试环境  
autotest run -c config.yml --environment=prod     # 生产环境
```

配置多环境变量：
```yaml
environments:
  dev:
    HOST: "localhost:8080"
    API_KEY: "dev_key"
  prod:
    HOST: "api.example.com" 
    API_KEY: "prod_key"
```

### 2. 重试机制
自动重试提升测试稳定性：
```yaml
global:
  retry:
    max_attempts: 5              # 失败后最多重试5次
    retry_delay: 2s              # 每次重试间隔2秒
    retry_on_status_codes: [500, 502, 503]  # 仅这些状态码重试
```

**适用场景**：
- 网络不稳定环境
- 服务偶尔抖动
- 负载过高导致超时

### 3. 智能缓存
gRPC描述符缓存优化性能：
```yaml
global:
  cache:
    enabled: true    # 默认开启，显著提升性能
    ttl: 600s       # 缓存10分钟
    max_size: 200   # 最多缓存200个描述符
```

**性能提升**：
- 1000个测试用例可节省50-200秒
- 减少重复的gRPC反射调用
- 自动启用，无需额外配置

### 4. 并发控制
平衡性能与服务器负载：
```yaml
global:
  concurrency:
    max_concurrent_requests: 10   # 限制最大并发数
    rate_limit_per_second: 30     # 每秒最多30个请求
```

**使用建议**：
- 测试环境可提高并发数
- 生产环境建议保守设置
- 根据目标服务性能调整

### 5. 多格式报告
生成多种格式的测试报告：
```yaml
global:
  report:
    dir_path: "./reports/"
    formats: ["html", "json", "junit"]  # 选择需要的格式
    template_path: "./custom.html"      # 自定义HTML模板
```

**报告格式**：
- **HTML**: 美观的网页报告  
- **JSON**: 程序化处理
- **CSV**: Excel分析
- **JUnit**: CI/CD集成

### 6. Slack通知
实时获取测试结果：
```yaml
global:
  notifications:
    enabled: true
    webhook_url: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
    on_failure: true   # 失败时通知
    on_success: true   # 成功时也通知
```

**通知内容**：
- 测试统计信息
- 失败用例详情
- 执行时间
- 通过率

### 7. 依赖关系
```yaml
- id: 2
  desc: "更新书籍"
  dependOnIDs: [1]  # 依赖测试用例1
  request:
    method: "put"
    url: "http://{{ HOST }}/api/books/{{ BOOK_ID }}"
    # ...
```

### 8. Lua脚本支持
```yaml
- id: 3
  desc: "动态生成请求体"
  request:
    method: "post"
    url: "http://{{ HOST }}/api/books"
    luaBody: |
      function body()
        local json = require "json"
        local timestamp = os.time()
        local data = {
          title = "Book-" .. timestamp,
          author = "Author-" .. timestamp
        }
        return json.encode(data)
      end
  rules:
    - name: "HttpLuaRule"
      lua: |
        function verify(r)
          local json = require "json"
          local body = json.decode(r:body())
          return body.title:match("Book%-") ~= nil
        end
```

### 9. XPath 表达式
```yaml
rules:
  # 精确匹配
  - name: "HttpBodyEqualRule"
    xpath: "/data/books[1]/title"
    expected: "特定书名"
    
  # 至少包含一个
  - name: "HttpBodyAtLeastOneRule"
    xpath: "//title"
    expected: "Go语言编程"
```

## 最佳实践

### 1. 测试用例组织
- 按功能模块分组测试用例
- 使用有意义的ID和描述
- 合理设置依赖关系

### 2. 环境管理
- **推荐使用** `--environment` 参数替代环境文件
- 在配置文件中定义所有环境的变量
- 避免在不同环境间硬编码差异

```yaml
# 推荐做法
environments:
  dev:
    HOST: "localhost:8080"
    TIMEOUT: "5s"
  prod:  
    HOST: "api.example.com"
    TIMEOUT: "30s"
```

### 3. 稳定性最佳实践
- **启用重试机制**：默认已开启，针对网络问题自动重试
- **合理设置并发**：根据目标服务能力调整并发数
- **监控缓存效果**：查看日志中的缓存命中率

```yaml
# 生产环境推荐配置
global:
  retry:
    max_attempts: 3
    retry_delay: 2s
  concurrency:
    max_concurrent_requests: 10  # 保守设置
    rate_limit_per_second: 20
```

### 4. 报告和通知
- **多格式报告**：根据用途选择合适格式
  - CI/CD流水线：使用 `junit` 格式
  - 人员查看：使用 `html` 格式
  - 数据分析：使用 `json` 或 `csv` 格式
- **智能通知**：只在失败时通知，减少噪音

```yaml
global:
  report:
    formats: ["html", "junit"]  # 兼顾查看和CI集成
  notifications:
    on_failure: true    # 只通知失败
    on_success: false   # 减少噪音
```

### 5. 性能优化
- **缓存默认开启**：gRPC描述符缓存自动优化性能
- 合理设置worker_num：通常设置为CPU核数的2倍
- 避免过多的依赖关系
- 大批量测试时适当降低并发数以保护目标服务

### 6. 变量管理  
- 通过export功能在测试用例间传递数据
- 避免硬编码敏感信息
- 使用环境变量进行配置隔离

## 故障排除

### 常见问题

1. **配置文件解析错误**
   ```bash
   # 检查YAML语法
   autotest test --config-file=./config.yml
   ```

2. **环境选择问题**
   ```bash
   # 检查可用环境
   autotest run --config-file=./config.yml --environment=nonexist
   # 错误信息会显示可用环境列表
   ```

3. **网络连接超时**
   ```yaml
   # 方案1: 增加超时时间
   global:
     request_timeout: 30s
   
   # 方案2: 启用重试机制（默认已开启）
   global:
     retry:
       max_attempts: 5
       retry_delay: 3s
   ```

4. **并发控制问题**
   ```yaml
   # 目标服务压力过大时，降低并发
   global:
     concurrency:
       max_concurrent_requests: 5
       rate_limit_per_second: 10
   ```

5. **缓存问题**
   ```yaml
   # 如需禁用缓存进行调试
   global:
     cache:
       enabled: false
   ```

6. **通知发送失败**
   - 检查Slack webhook URL是否正确
   - 验证网络连接
   - 确认webhook权限设置

7. **报告生成失败**
   - 检查报告目录是否存在和有写权限
   - 验证自定义模板文件路径
   - 确认所选报告格式受支持

8. **XPath表达式错误**
   ```bash
   # 测试XPath表达式
   autotest extract --xpath="//title" --json='{"title":"test"}'
   ```

9. **Lua脚本错误**
   - 检查语法错误
   - 确保函数名正确（body/verify）
   - 验证JSON解析逻辑

### 调试技巧

#### 基础调试
- 开启debug模式查看详细日志
- 使用extract命令测试XPath
- 逐步增加测试用例复杂度
- 检查环境变量是否正确设置

#### 企业功能调试
- **重试调试**：查看日志中的重试次数和原因
- **缓存调试**：查看缓存命中率统计
- **并发调试**：观察请求速率是否符合限制
- **通知调试**：先测试简单的成功场景

#### 性能分析
```bash
# 查看缓存统计
# 日志中会显示缓存命中率
autotest run -c config.yml --environment=dev

# 测试不同并发设置的性能影响
# 调整 max_concurrent_requests 参数对比执行时间
```

#### 环境隔离调试
```bash
# 分别测试不同环境
autotest run -c config.yml --environment=dev
autotest run -c config.yml --environment=staging

# 对比环境间的差异
```

## 企业级使用场景

### CI/CD 集成
```yaml
# 适用于持续集成的配置
global:
  retry:
    max_attempts: 3
    retry_delay: 1s
  concurrency:
    max_concurrent_requests: 15
    rate_limit_per_second: 40
  report:
    formats: ["junit", "json"]  # CI友好格式
  notifications:
    enabled: true
    on_failure: true
    on_success: false
```

```bash
# CI流水线中的用法
autotest run -c config.yml --environment=staging
# 生成JUnit报告供CI系统解析
```

### 大规模测试场景（1000+用例）
```yaml
# 针对大批量测试的优化配置  
global:
  worker_num: 10
  concurrency:
    max_concurrent_requests: 8   # 适度并发保护目标服务
    rate_limit_per_second: 25    # 控制请求频率
  cache:
    enabled: true               # 必须开启，显著提升性能
    ttl: 600s                  # 延长缓存时间
    max_size: 300              # 增加缓存容量
  retry:
    max_attempts: 2            # 减少重试次数节省时间
    retry_delay: 500ms         # 缩短重试间隔
```

### 生产环境测试
```yaml
# 生产环境安全配置
global:
  concurrency:
    max_concurrent_requests: 3   # 保守并发设置
    rate_limit_per_second: 10    # 严格限制频率
  retry:
    max_attempts: 1             # 避免对生产服务造成压力
  notifications:
    enabled: true
    on_failure: true            # 及时通知生产问题
    on_success: true            # 确认生产测试成功
```

### 团队协作场景
```yaml
# 多团队共享配置
environments:
  # 开发团队
  team-a-dev:
    HOST: "team-a.dev.internal"
    API_KEY: "${TEAM_A_DEV_KEY}"
  
  # 测试团队  
  qa-staging:
    HOST: "qa.staging.internal"
    API_KEY: "${QA_STAGING_KEY}"
    
  # 运维团队
  ops-prod:
    HOST: "api.production.com"
    API_KEY: "${OPS_PROD_KEY}"

global:
  notifications:
    webhook_url: "${TEAM_SLACK_WEBHOOK}"  # 团队专属通知
```

### 性能基准测试
```bash
# 测试不同配置的性能差异

# 基础配置
autotest run -c basic-config.yml --environment=dev

# 优化配置（启用缓存+适度并发）
autotest run -c optimized-config.yml --environment=dev

# 对比执行时间和成功率
```
