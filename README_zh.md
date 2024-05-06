# autotest

## 概览
针对api服务，http、gRPC(暂时不支持)的自动化测试框架

## 特点:
* 无需进行程序开发，只需要编写配置文件
* 可以指定testcase之间的依赖关系
* 无依赖关系的testcase可以并发执行，执行速度更快
* 使用xpath提取变量书写方便

## 用法
### 1) 检查配置文件
``` 
autotest --config-file=${CONFIG_FILE}
```

### 2) 执行自动化测试
``` 
autotest --config-file=${CONFIG_FILE} --env-file=${ENV_FILE}
```

## 示例
### 1) 启动一个api服务

### 2) 运行自动化测试用例
```
make build
./autotest run -c=./config_files/autotest.yml -e=./config_files/.env.dev
```

