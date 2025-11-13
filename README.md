# autotest

[![golang-ci](https://github.com/vearne/autotest/actions/workflows/golang-ci.yml/badge.svg)](https://github.com/vearne/autotest/actions/workflows/golang-ci.yml)

* [中文 README](https://github.com/vearne/autotest/blob/master/README_zh.md)

## 1.Overview
An automated testing framework for API services, such as HTTP and gRPC.

## 2.Features
### Core Features
* No program development is required, only configuration files need to be written
* You can specify dependencies between testcases
* Testcases without dependencies can be executed concurrently and execute faster
* Use XPath to extract variables for easy writing
* Supports importing variables from files and extracting variables from response

### Enterprise Features 🚀
* **Auto Retry**: Intelligent retry mechanism for network failures and server errors
* **Smart Caching**: gRPC descriptor caching for significant performance improvements
* **Concurrency Control**: Rate limiting and concurrent request management
* **Multi-format Reports**: Generate HTML, JSON, CSV, and JUnit reports
* **Slack Integration**: Real-time notifications via webhook
* **Environment Management**: Easy environment switching (dev/staging/prod)

## 3.Something you need to know
[XPath Syntax](https://www.w3schools.com/xml/xpath_syntax.asp)

## 4.Install
### 1) use the compiled binary file
[release](https://github.com/vearne/autotest/releases)

Obtain the bin file corresponding to the operating system and CPU architecture from the link above
### 2) compile by yourself
```
make build
```
or
```
go install github.com/vearne/autotest@latest
```


## 5.Usage
### 1) check configuration file
``` 
autotest test --config-file=${CONFIG_FILE}
```

### 2) execute automated tests
``` 
# Basic usage with env file
autotest run --config-file=${CONFIG_FILE} --env-file=${ENV_FILE}

# Advanced usage with environment selection
autotest run --config-file=${CONFIG_FILE} --environment=dev
autotest run --config-file=${CONFIG_FILE} --environment=prod
```
### 3) extract the value corresponding to xpath
``` 
autotest extract --xpath=${XPATH} --json=${JSON}
```

## 6.Example
### 1) start a fake http api service
```
cd ./docker-compose
docker compose up -d
```
#### Add
```
curl -X POST 'http://localhost:8080/api/books' \
--header 'Content-Type: application/json' \
--data '{"title": "book3_title", "author": "book3_author"}'
```

#### Delete
```
curl -X DELETE 'http://localhost:8080/api/books/1'
```

#### Modify
```
curl -X PUT 'localhost:8080/api/books/3' \
--header 'Content-Type: application/json' \
--data '{"title": "book3_title", "author": "book3_author-2"}'
```
#### List
```
curl  'http://localhost:8080/api/books'
```


### 2) run automated test cases
```
autotest run -c=./config_files/autotest.yml -e=./config_files/.env.dev
```

### 3) extract the value corresponding to xpath
get the title of each book in the book list
```
autotest extract -x "//title" -j '[
 {
  "id": 2,
  "title": "Effective Go",
  "author": "The Go Authors"
 },
 {
  "id": 3,
  "title": "book3_title",
  "author": "book3_author-2"
 }
]'
```
## 7.Test Reports & Notifications

### Multi-format Report Generation
AutoTest now supports multiple report formats:

- **HTML**: Beautiful web-based reports with custom templates
- **JSON**: Structured data for programmatic processing  
- **CSV**: Tabular data for Excel analysis
- **JUnit**: CI/CD integration friendly XML format

Configuration example:
```yaml
global:
  report:
    dir_path: "/var/log/test/report/"
    formats: ["html", "json", "csv", "junit"]  # Generate multiple formats
    template_path: "./templates/custom.html"   # Optional custom template
```

### Slack Notifications 📢
Get real-time test results in Slack:

```yaml
global:
  notifications:
    enabled: true
    webhook_url: "https://hooks.slack.com/services/..."
    on_failure: true   # Notify on test failures
    on_success: false  # Optional success notifications
```

### Sample Reports
![CSV Report](https://github.com/vearne/autotest/raw/main/img/result_csv.jpg)
![HTML Report](https://github.com/vearne/autotest/raw/main/img/result_html.jpg)

## 8.Enterprise Configuration

### Performance & Reliability
```yaml
global:
  # Retry mechanism for better stability
  retry:
    max_attempts: 3                        # Auto-retry failed requests
    retry_delay: 1s                        # Wait between retries
    retry_on_status_codes: [500, 502, 503, 504]

  # Concurrency control  
  concurrency:
    max_concurrent_requests: 20            # Limit concurrent requests
    rate_limit_per_second: 50              # Rate limiting

  # Smart caching (enabled by default)
  cache:
    enabled: true                          # gRPC descriptor caching
    ttl: 300s                             # Cache TTL
    max_size: 100                         # Cache entries
```

### Environment Management
```yaml
# Multi-environment support
environments:
  dev:
    HOST: "localhost:8080"
    GRPC_SERVER: "localhost:50031"
  staging:
    HOST: "staging.api.com"  
    GRPC_SERVER: "staging.grpc.com:50031"
  prod:
    HOST: "api.example.com"
    GRPC_SERVER: "grpc.example.com:50031"
```

Usage:
```bash
autotest run -c config.yml --environment=prod
```

## 9.Advanced Usage
In certain scenarios, we may need to use Lua scripts to generate the request body 
or to verify whether the response body meets expectations
```
- id: 6
  desc: "add a new book"
  request:
    # optional
    method: "post"
    url: "http://{{ HOST }}/api/books"
    headers:
      - "Content-Type: application/json"
    luaBody: |
      function body()
        local json = require "json";
        -- the string representation of today's date at 23:59:59.
        local today235959 = os.date("%Y%m%d235959");
        local data = {
          title   = "book4_title-" .. today235959,
          author  = "book4_author"
        };
        return json.encode(data);
      end
  rules:
    - name: "HttpStatusEqualRule"
      expected: 200
    - name: "HttpBodyEqualRule"
      xpath: "/author"
      expected: "book4_author"
    - name: "HttpLuaRule"
      lua: |
        function verify(r)
          local json = require "json";
          local book = json.decode(r:body());
          print("book.title:", book.title);
          print("---1---", 10);
          print("---2---", 20);
          local today235959 = os.date("%Y%m%d235959");
          local title = "book4_title-" .. today235959;
          return book.title == title;
        end
```
**Notice**:
* In LuaBody scripts, the form of functions is fixed.
```lua
function body()
```
* In HttpLuaRule.lua scripts, the form of functions is fixed too.
```
function verify(r)
```