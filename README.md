# NAT 内网穿透工具

<img src="./images/logo.png" style="zoom:50%;" />

![](https://img.shields.io/badge/GO-v1.20-blue)

![](https://img.shields.io/badge/release-v0.05-green)

![](https://img.shields.io/badge/auth-pogf-lightgrey)

## 1. 如何使用

**服务器端**

```shell
git clonehttps://github.com/byteYuFan/NAT.git
cd NAT/server
go build .
server 
```

```yaml
Server:
  Name: "Server-NAT"
  ControllerAddr: "0.0.0.0:8007"
  TunnelAddr: "0.0.0.0:8008"
  VisitPort:
    - 60000
    - 60001
    - 60002
    - 60003
  TaskQueueNum: 4
  TaskQueueBuff: 32
  MaxTCPConnNum: 4
  MaxConnNum: 256
```

根据具体的情况修改`yaml`文件的参数

现象：

```shell
root@VM-4-7-centos ~]# ./server --help
GO language-based Intranet penetration tool that supports multiple connections

Usage:
  Server-NAT [OPTIONS] COMMAND [flags]

Flags:
  -c, --controller-addr string         Server controller address
  -p, --expose-port ints               Server exposed ports
  -h, --help                           help for Server-NAT
  -x, --max-conn-num int32             Maximum connection number
  -m, --max-tcp-conn-num int32         Maximum TCP connection number
  -n, --name string                    Server name
  -b, --task-queue-buffer-size int32   Task queue buffer size
  -q, --task-queue-num int32           Task queue number
  -t, --tunnel-addr string             Server tunnel address
[root@VM-4-7-centos ~]# ./server
2023/06/05 13:36:04 [公网服务器控制端开始监听]0.0.0.0:8080
[ListenTaskQueue] 监听工作队列传来的消息
```

客户端:

说明：在启动客户端之前需要需改相应的配置文件

```yaml	
Client:
  Name: "Client-NAT"
  PublicServerAddr: "8.8.8.8"
  TunnelServerAddr: "8.8.8.8:8008"
  ControllerAddr: "8.8.8.8:8007"
  LocalServerAddr: "127.0.0.1:8080"
```

前三个地址和公网服务器所对应的地址相同即可，` LocalServerAddr`为本地服务器地址



启动测试服务:

```go
D:\goworkplace\src\test>test.exe
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func1 (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080


```

启动客户端程序

```shell
git clonehttps://github.com/byteYuFan/NAT.git
cd NAT/client
go build .
./client 
$ ./client.exe  --help
If the intranet is written in the go language, you need to start the intranet client before you can connect

Usage:
  Client [OPTIONS] COMMAND [flags]

Flags:
  -c, --controller-addr string      The address of the controller channel used to send controller messages to the client
  -h, --help                        help for Client
  -l, --local-server-addr string    The address of the local web server program
  -n, --name string                 Client name
  -s, --public-server-addr string   The address of the public server used for accessing the inner web server
  -t, --tunnel-server-addr string   The address of the tunnel server used to connect the local and public networks

$ ./client.exe  
2023/06/05 13:45:26 [Conn Successfully]公网ip(8.8.8.8):8080
[Heart] ping
[Byte] [0 0 0 0 0 0 0 1 0 0 234 96]
[ClientInfo] &{1 60000}
[ClientInfo] &{1 60000}
[Heart] ping
[Heart] ping
[Heart] ping
[Heart] ping
```

测试结果:

```shell
D:\goworkplace\src\github.com\test>curl 公网ip(8.8.8.8):60000
{"message":"Hello, World!"}

[GIN] 2023/06/05 - 13:46:24 |[97;42m] 200 [0m|     912.5µs |     127.0.0.1 |[97;44m GET]  [0m] "/"
```



