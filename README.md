# NAT 内网穿透工具

基于`GO`语言开发的一款内网穿透工具。目前版本`v0.1`实现了具体的功能。

## 一、项目详细介绍

[开发文档](./README-CN.MD)

## 二、如何使用

客户端

```go
var (
	// 公网ip
	publicAddr = "8.8.8.8"
	// 本地服务的地址
	localServerAddr = "127.0.0.1:8080"
	// 公网服务端的控制接口
	controllerServerAddr = publicAddr + ":8080"
	// 公网隧道地址
	tunnelServerAddr = publicAddr + ":8008"
)
```

服务器

```go
const (
	// 控制信息地址
	controllerAddr = "0.0.0.0:8009"
	// 隧道地址
	tunnelAddr = "0.0.0.0:8008"
	// 外部访问地址
	visitAddr = "0.0.0.0:8007"
)
```

对应参数修改完成后，启动即可
