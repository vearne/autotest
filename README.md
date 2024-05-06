# autotest

* [中文 README](https://github.com/vearne/autotest/blob/master/README_zh.md)

## Overview
Automated testing framework for api services, http, gRPC (may support in the future)

## Features
* No program development is required, only configuration files need to be written
* You can specify dependencies between testcases
* Testcases without dependencies can be executed concurrently and execute faster
* Use xpath to extract variables for easy writing

## Usage
### 1) check configuration file
``` 
autotest --config-file=${CONFIG_FILE}
```

### 2) execute automated tests
``` 
autotest --config-file=${CONFIG_FILE} --env-file=${ENV_FILE}
```
## Example
### 1) start a fake api service

### 2) run automated test cases
```
make build
./autotest run -c=./config_files/autotest.yml -e=./config_files/.env.dev
```
