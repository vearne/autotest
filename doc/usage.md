# AutoTest 使用指南

## 目录
- [快速开始](#快速开始)
- [配置文件详解](#配置文件详解)
- [测试用例编写](#测试用例编写)
- [高级功能](#高级功能)
- [最佳实践](#最佳实践)
- [故障排除](#故障排除)

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
autotest run --config-file=./config_files/autotest.yml --env-file=./config_files/.env.dev
```

## 配置文件详解

### 全局配置 (autotest.yml)
```yaml
global:
  worker_num: 5                    # 并发工作协程数
  ignore_testcase_fail: true       # 是否忽略测试用例失败
  debug: false                     # 调试模式
  request_timeout: 5s              # 请求超时时间
  
  logger:
    level: "debug"                 # 日志级别: debug/info/error
    file_path: "/var/log/test/autotest.log"
    
  report:
    dir_path: "/var/log/test/report/"  # 报告输出目录

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

### 1. 依赖关系
```yaml
- id: 2
  desc: "更新书籍"
  dependOnIDs: [1]  # 依赖测试用例1
  request:
    method: "put"
    url: "http://{{ HOST }}/api/books/{{ BOOK_ID }}"
    # ...
```

### 2. Lua脚本支持
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

### 3. XPath 表达式
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

### 2. 变量管理
- 使用环境变量文件管理不同环境配置
- 通过export功能在测试用例间传递数据
- 避免硬编码敏感信息

### 3. 错误处理
- 设置合理的超时时间
- 使用重试机制处理网络波动
- 编写详细的验证规则

### 4. 性能优化
- 合理设置worker_num
- 避免过多的依赖关系
- 使用缓存减少重复请求

## 故障排除

### 常见问题

1. **配置文件解析错误**
   ```bash
   # 检查YAML语法
   autotest test --config-file=./config.yml
   ```

2. **网络连接超时**
   ```yaml
   # 增加超时时间
   global:
     request_timeout: 30s
   ```

3. **XPath表达式错误**
   ```bash
   # 测试XPath表达式
   autotest extract --xpath="//title" --json='{"title":"test"}'
   ```

4. **Lua脚本错误**
   - 检查语法错误
   - 确保函数名正确（body/verify）
   - 验证JSON解析逻辑

### 调试技巧
- 开启debug模式查看详细日志
- 使用extract命令测试XPath
- 逐步增加测试用例复杂度
- 检查环境变量是否正确设置
