# g

一个 code generator for golang

## 主要功能

1. 自动生成 if err != nil, 并做了碰撞检查（如果已经写了 if err != nil ，那么不会再生成一次）

2. 自动生成 defer recover

3. 集成了 wire 做编译时依赖注入

TODO

4. 自动生成 log 语句

5. 自动生成 json tag

6. 自动生成 ok check。 if _,ok := a.(*b); ok {}


## 演示

执行本项目下的 .example 的 gen_test.go 文件

将扫描项目下带有 // +build dev 标识的 go 文件

并生成这些文件：`autoerr.gen.go`, `safeerr.gen.go`, `wire.gen.go`

## 声明

请勿直接用在生产环境

需要使用请二次开发