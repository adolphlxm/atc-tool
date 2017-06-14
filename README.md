# atc-tool

atc-tool 是一个命令行工具

## 安装ATC

    go get github.com/adolphlxm/atc-tool
    
  **同时你需要安装如下依赖**
  
  ```config 
    // ATC, orm 包
    github.com/adolphlxm/atc/orm
    // ATC, logs 日志包
    github.com/adolphlxm/atc/logs
  ```
   
## 命令列表
* orm 数据库表字段反转,生产.go文件
* thrift 定义IDL文件生成.go源码包

## atc-tool orm

    $ atc-tool orm [-json] driverName datasourceName tableName [generatedPath]
 
 *Options:*
 
 * -json [是否支持JSON TAG]
    - true or false(默认) 
 * driverName [数据库驱动]
 * datasourceName [数据库连接资源]
 * tableName [表名]
 * generatedPath [数据库表反转后生成文件]
    - ../user 则会在该目录下生成一个[表名].go文件
    
举例：

    $ example
    atc-tool orm mysql "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8" user ./model/
    ...
    [ATC] [Trace] 2017/05/31 19:19:39.329806 reverse.go#225: Reverse create ./users/user.go successfully created!


`user.go`文件内容如下：

```go
package conf

type User struct {
	Id     int64 `xorm:"int(11)"`
	Number int64 `xorm:"int(11)"`
}

```

## atc-tool thrift

* 依赖thrift
* 解决问题
    - 由于ATC框架修改了Thrift go server源码
    - 服务端thrift源码 引用的是 `github.com/adolphlxm/rpc/thrift/lib/thrift`包
    - 使用该工具生成，可自动替换 thrift 生成.go文件的引用包问题。

thrift命令行工具，用于 thrift IDL .go 文件生产

    $ atc-tool thrift [options] file
    
具体thrift命令可使用 `thrift --help` 查看

举例：

    $ atc-tool thrift -r --gen go xxx.thrift
    
即可在该目录生成一个gen-go/xxx 的文件

## 更新日志

* 2017.5 
    - 修复表名带`_` 大写转换成struct问题
    
# LICENSE

ATC is Licensed under the Apache License, Version 2.0 (the "License")
(http://www.apache.org/licenses/LICENSE-2.0.html).