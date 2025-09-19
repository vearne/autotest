# autotest

[![golang-ci](https://github.com/vearne/autotest/actions/workflows/golang-ci.yml/badge.svg)](https://github.com/vearne/autotest/actions/workflows/golang-ci.yml)

## 1.概览
针对api服务，http、gRPC的自动化测试框架

## 2.特点:
* 无需进行程序开发，只需要编写配置文件
* 可以指定testcase之间的依赖关系
* 无依赖关系的testcase可以并发执行，执行速度更快
* 使用XPath提取变量，书写方便
* 支持从文件中导入变量，支持从response中提取变量

## 3.你需要了解的知识
[XPath Syntax](https://www.w3schools.com/xml/xpath_syntax.asp)

## 4.安装
### 1) 使用编译好的bin文件
[release](https://github.com/vearne/autotest/releases)

从上面的链接获取对应操作系统和cpu架构的bin文件
### 2) 手动编译
```
make build
```
或
```
go install github.com/vearne/autotest@latest
```

## 5.用法
### 1) 检查配置文件
``` 
autotest test --config-file=${CONFIG_FILE}
```

### 2) 执行自动化测试
``` 
autotest run --config-file=${CONFIG_FILE} --env-file=${ENV_FILE}
```

### 3) 提取xpath对应的值
``` 
autotest extract --xpath=${XPATH} --json=${JSON}
```

## 6.示例
### 1) 启动一个伪造的http api服务
```
cd ./docker-compose
docker compose up -d
```
#### 添加
```
curl -X POST 'http://localhost:8080/api/books' \
--header 'Content-Type: application/json' \
--data '{"title": "book3_title", "author": "book3_author"}'
```

#### 删除
```
curl -X DELETE 'http://localhost:8080/api/books/1'
```

#### 修改
```
curl -X PUT 'localhost:8080/api/books/3' \
--header 'Content-Type: application/json' \
--data '{"title": "book3_title", "author": "book3_author-2"}'
```
#### 列表
```
curl  'http://localhost:8080/api/books'
```

### 2) 运行自动化测试用例
```
autotest run -c=./config_files/autotest.yml -e=./config_files/.env.dev
```

### 3) 提取xpath对应的值
获取书本列表中，书的title
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
## 7.测试报告
### CSV格式
![report](https://github.com/vearne/autotest/raw/main/img/result_csv.jpg)

### HTML格式
![report](https://github.com/vearne/autotest/raw/main/img/result_html.jpg)

## 8.高级用法
某些场景，我们可能需要使用lua脚本来生成请求的body，或者用脚本来验证响应的body是否符合预期
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
        -- 今天 23:59:59 的字符串
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
**注意**: 
* luaBody脚本中函数形式是固定的
```lua
function body()
```
* HttpLuaRule脚本中函数形式也是是固定的
```
function verify(r)
```


