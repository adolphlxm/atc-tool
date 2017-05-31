# atc-tool

atc-tool 是一个命令行工具

## 安装ATC

    go get github.com/adolphlxm/atc-tool
    
  **同时你需要安装如下依赖**
  
  ```config 
    // ATC框架,用于ORM反转及日志支持
    github.com/adolphlxm/atc
 
  ```
   
## 命令列表
* orm 数据库表字段反转,生产.go文件

## atc-tool orm

    $ atc-tool orm [dbName] [table] [-c=config path] [-s=Generate the file path] [-j= true|false]
 
 * dbName 数据库名称
 * table 表名
 * -c ATC配置文件路径
    - e.g. ../conf/app.ini
 * -s 数据库表反转后生成的具体文件路径
    - e.g. ../user 则会在该目录下生成一个[表名].go文件
 * -j 是否支持JSON TAG
    - e.g.  true 表示支持JSON TAG
    
举例：

    $ example
    atc-tool orm test user -c="../conf/app.ini" -s ="../users"
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
# LICENSE

ATC is Licensed under the Apache License, Version 2.0 (the "License")
(http://www.apache.org/licenses/LICENSE-2.0.html).