# autotest

[![golang-ci](https://github.com/vearne/autotest/actions/workflows/golang-ci.yml/badge.svg)](https://github.com/vearne/autotest/actions/workflows/golang-ci.yml)

## 概览
针对api服务，http、gRPC的自动化测试框架

## 特点:
* 无需进行程序开发，只需要编写配置文件
* 可以指定testcase之间的依赖关系
* 无依赖关系的testcase可以并发执行，执行速度更快
* 使用XPath提取变量，书写方便
* 支持从文件中导入变量，支持从response中提取变量

## 你需要了解的知识
[XPath Syntax](https://www.w3schools.com/xml/xpath_syntax.asp)

## 安装
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

## 用法
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

## 示例
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
## 测试报告
### CSV格式
![report](https://github.com/vearne/autotest/raw/main/img/result_csv.jpg)

### HTML格式
![report](https://github.com/vearne/autotest/raw/main/img/result_html.jpg)

## TODO
* [x] 1) 支持使用脚本语言Lua判断HTTP response是否符合预期
* [x] 2) 输出report到文件中
* [x] 3) 支持对gRPC协议的API服务进行自动化测试



