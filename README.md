# atc-tool

atc-tool 是一个命令行工具

## 安装ATC-TOOL

**通过GO命令安装**

    go get github.com/adolphlxm/atc-tool

**通过源码安装**

    git clone https://github.com/adolphlxm/atc-tool.git ./atc-tool
    go install ./atc-tool
   
## 命令列表
* atc-tool new 创建一个项目
* atc-tool orm 数据库表字段反转,生产.go文件
* atc-tool thrift 定义IDL文件生成.go源码包

## atc-tool new
创建一个ATC项目

    $ atc-tool new [appname]

```
├── conf [配置文件目录]
│   ├── app.ini [核心配置文件]
│   └── error.ini [错误码文件]
├── bin [可执行文件目录]
├── src [源码目录]
│   ├── api [RESTFul API 目录]
│         ├── V1 [版本目录]
│         └── router.go [路由文件]
│   └── model [DB目录]
│   └── service [服务目录]
│   └── thrift [RPC]
│         ├── idl [存放Thrift IDL文件目录]
│         ├── gen-go [存放Thrift 生成 go文件目录]
│         ├── ...(.go)
│         └── router.go [路由文件]
└── atc.go [运行主程序文件]
```

启动运行后，可通过以下生成的RESTFul HTTP API 访问：
* http://localhost/api
* http://localhost/api2
* http://localhost/api2/123
* http://localhost/ap2/test

来验证创建的项目DEMO是否通过

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
* 2017.6
    - 优化生成数据表反转代码格式化(go fmt file.)
* 2017.12
    - 增加`new`创新新项目命令
    
# LICENSE

ATC is Licensed under the Apache License, Version 2.0 (the "License")
(http://www.apache.org/licenses/LICENSE-2.0.html).