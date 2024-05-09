# autotest

* [中文 README](https://github.com/vearne/autotest/blob/master/README_zh.md)

## Overview
Automated testing framework for api services, http, gRPC (may support in the future)

## Features
* No program development is required, only configuration files need to be written
* You can specify dependencies between testcases
* Testcases without dependencies can be executed concurrently and execute faster
* Use xpath to extract variables for easy writing
* supports importing variables from files and extracting variables from response

## Install
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


## Usage
### 1) check configuration file
``` 
autotest test --config-file=${CONFIG_FILE}
```

### 2) execute automated tests
``` 
autotest run --config-file=${CONFIG_FILE} --env-file=${ENV_FILE}
```
### 3) extract the value corresponding to xpath
``` 
autotest extract --xpath=${XPATH} --json=${JSON}
```

## Example
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
