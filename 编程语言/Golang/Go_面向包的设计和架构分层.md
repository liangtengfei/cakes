# Go 面向包的设计和架构分层

## Overview（概述）

*This is a basic layout for Go application projects. It's not an official standard defined by the core Go dev team; however, it is a set of common historical and emerging project layout patterns in the Go ecosystem. Some of these patterns are more popular than others. It also has a number of small enhancements along with several supporting directories common to any large enough real world application.*
**（这是 Go 应用程序项目的基本布局。它不是核心 Go 开发团队定义的官方标准；但是，它是 Go 生态系统中一组常见的历史和新兴项目布局模式。其中一些模式比其他模式更受欢迎。它还具有许多小的增强功能，以及任何足够大的现实世界应用程序共有的几个支持目录。）**

## Go Directories（目录）

### /cmd
Main applications for this project.
### /internal
Private application and library code. 
### /pkg
Library code that's ok to use by external applications (e.g., /pkg/mypubliclib). 

## Service Application Directories（程序服务目录）
### /api
OpenAPI/Swagger specs, JSON schema files, protocol definition files.

See the /api directory for examples.

## 结尾
软件大师 Kent Beck 在《重构Refactoring》一书中描述
> 1. 先让代码工作起来-如果代码不能工作，就不能产生价值
> 1. 然后再试图将它变好-通过对代码进行重构，让我们自己和其他人更好地理解代码，并能按照需求不断地修改代码。
> 1. 最后再试着让它运行得更快-按照性能提升的需求来重构代码。

## 结构地址
[Standard Go Project Layout](https://github.com/golang-standards/project-layout)