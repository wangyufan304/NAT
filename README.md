# NAT 内网穿透工具

<img src="./images/logo.png" style="zoom:50%;" />

![](https://img.shields.io/badge/GO-v1.20-blue)

![](https://img.shields.io/badge/release-v0.1.2-green)

![](https://img.shields.io/badge/auth-yyffww-lightgrey)

# 服务端介绍

[详细文档](./NAT-new.md)

## 1. 如何使用

### 1.1. 基本使用

```shell
git clone https://github.com/byteYuFan/NAT.git
cd NAT/server
go build -o server
```

修改配置文件相关信息：

```yaml
Server:
  Name: "Server-NAT"
  ControllerAddr: "0.0.0.0:8080"
  TunnelAddr: "0.0.0.0:8008"
  VisitPort:
    - 25565
    - 57852
    - 64251
    - 12541
  TaskQueueNum: 4
  TaskQueueBuff: 32
  MaxTCPConnNum: 4
  MaxConnNum: 128
  StartAuth: true
Database:
 Username: "root"
 Password: "123456"
 Host: "127.0.0.1:3309"
 DBName: "NAT"

```

```go
// 执行命令
[root@VM-4-7-centos ~]# ./server
  _   _              _______    _____   ______   _____   __      __  ______   _____
 | \ | |     /\     |__   __|  / ____| |  ____| |  __ \  \ \    / / |  ____| |  __ \
 |  \| |    /  \       | |    | (___   | |__    | |__) |  \ \  / /  | |__    | |__) |
 | . ` |   / /\ \      | |     \___ \  |  __|   |  _  /    \ \/ /   |  __|   |  _  /
 | |\  |  / ____ \     | |     ____) | | |____  | | \ \     \  /    | |____  | | \ \
 |_| \_| /_/    \_\    |_|    |_____/  |______| |_|  \_\     \/     |______| |_|  \_\


[ServerName] Server-NAT
[MaxServerConn] 4
[服务端开启端口] [25565 57852 64251 12541]
[服务器控制端开始监听]0.0.0.0:8080
[Start Auth Successfully!] 服务器开启认证请求
```

### 1.2. 命令行参数

命令行的优先级是最高的，您可以在命令行输入相关的配置信息

```shell
[root@VM-4-7-centos ~]# ./server --help
GO language-based Intranet penetration tool that supports multiple connections

Usage:
  Server-NAT [OPTIONS] COMMAND [flags]

Flags:
  -c, --controller-addr string         Server controller address
  -p, --expose-port ints               Server exposed ports
  -h, --help                           help for Server-NAT
  -l, --log-name string                The name of the log.
  -x, --max-conn-num int32             Maximum connection number
  -m, --max-tcp-conn-num int32         Maximum TCP connection number
  -n, --name string                    Server name
  -a, --start-auth string              This is the method that whether the server start the auth. (default "true")
  -b, --task-queue-buffer-size int32   Task queue buffer size
  -q, --task-queue-num int32           Task queue number
  -t, --tunnel-addr string             Server tunnel address

```

如果您开启了不需要进行认证功能，那么您不需要进行配置数据库模块：

```shell

[root@VM-4-7-centos ~]# ./server --start-auth=false
  _   _              _______    _____   ______   _____   __      __  ______   _____
 | \ | |     /\     |__   __|  / ____| |  ____| |  __ \  \ \    / / |  ____| |  __ \
 |  \| |    /  \       | |    | (___   | |__    | |__) |  \ \  / /  | |__    | |__) |
 | . ` |   / /\ \      | |     \___ \  |  __|   |  _  /    \ \/ /   |  __|   |  _  /
 | |\  |  / ____ \     | |     ____) | | |____  | | \ \     \  /    | |____  | | \ \
 |_| \_| /_/    \_\    |_|    |_____/  |______| |_|  \_\     \/     |______| |_|  \_\


[ServerName] Server-NAT
[MaxServerConn] 4
[服务端开启端口] [25565 57852 64251 12541]
[服务器控制端开始监听]0.0.0.0:8080
```

# 客户端介绍

## 1. 如何使用

### 1.1. 基本使用

```shell
git clone https://github.com/byteYuFan/NAT.git
cd NAT/client
go build -o client
```

修改配置文件：

```yaml
Client:
  Name: "Client-NAT"
  PublicServerAddr: "公网服务器域名"
  TunnelServerAddr: "公网服务器隧道端口"
  ControllerAddr: "公网服务器控制端口"
  LocalServerAddr: "本地服务端口"
Auth:
  Username: "用户名"
  Password: "密码"
```

```go
$ ./client.exe 
  _   _              _______    _____   _        _____   ______   _   _   _______ 
 | \ | |     /\     |__   __|  / ____| | |      |_   _| |  ____| | \ | | |__   __|
 |  \| |    /  \       | |    | |      | |        | |   | |__    |  \| |    | |   
 | . ` |   / /\ \      | |    | |      | |        | |   |  __|   | . ` |    | |   
 | |\  |  / ____ \     | |    | |____  | |____   _| |_  | |____  | |\  |    | |   
 |_| \_| /_/    \_\    |_|     \_____| |______| |_____| |______| |_| \_|    |_|   
                                                                                  
                                                                                  
[Client Running Successfully!]
[PublicAddress] 
[TunnelAddress] :8008
[LocalAddress] 127.0.0.1:8080
[Conn Successfully] :8080
[ClientInfoUID] 1
[VisitAddress] :25565
[receive KeepLive package] ping

```

### 1.2. 命令行参数

```shell
$ ./client.exe --help
If the intranet is written in the go language, you need to start the intranet client before you can connect

Usage:
  Client [OPTIONS] COMMAND [flags]

Flags:
  -c, --controller-addr string      The address of the controller channel used to send controller messages to the client
  -h, --help                        help for Client
  -l, --local-server-addr string    The address of the local web server program
  -n, --name string                 Client name
  -P, --password string             the password for auth the server.
  -s, --public-server-addr string   The address of the public server used for accessing the inner web server
  -t, --tunnel-server-addr string   The address of the tunnel server used to connect the local and public networks
  -u, --username string             the name for auth the server.

```

# 案例测试

## 1. 基本web接口代理

内网客户端代码:

```shell
[Client Running Successfully!]
[PublicAddress] pogf.com.cn
[TunnelAddress] pogf.com.cn:8008
[LocalAddress] 127.0.0.1:8080
[Conn Successfully]pogf.com.cn:8080
[ClientInfoUID] 1
[VisitAddress] pogf.com.cn:25565

```

```go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 创建一个Gin的默认引擎
	r := gin.Default()

	// 定义一个路由处理函数
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// 启动Web服务，监听在8080端口
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("启动Web服务失败: ", err)
	}
}
```

```go
D:\goworkplace\src>curl pogf.com.cn:25565
{"message":"Hello, World!"}
```

## 2. 代理web接口

```shell
[root@localhost ~]# ./licent
  _   _              _______    _____   _        _____   ______   _   _   _______
 | \ | |     /\     |__   __|  / ____| | |      |_   _| |  ____| | \ | | |__   __|
 |  \| |    /  \       | |    | |      | |        | |   | |__    |  \| |    | |
 | . ` |   / /\ \      | |    | |      | |        | |   |  __|   | . ` |    | |
 | |\  |  / ____ \     | |    | |____  | |____   _| |_  | |____  | |\  |    | |
 |_| \_| /_/    \_\    |_|     \_____| |______| |_____| |______| |_| \_|    |_|


[Client Running Successfully!]
[PublicAddress] pogf.com.cn
[TunnelAddress] pogf.com.cn:8008
[LocalAddress] 127.0.0.1:80
[Conn Successfully]pogf.com.cn:8080
[receive KeepLive package] ping
[ClientInfoUID] 3
[VisitAddress] pogf.com.cn:60002

```

```go
D:\goworkplace\src>curl pogf.com.cn:60002
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            background-color: #222;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }

        .clock {
            display: flex;
            justify-content: center;
            align-items: center;
            background-color: #fff;
            width: 280px;
            height: 140px;
            border-radius: 10px;
            box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
        }

        .clock .digit {
            font-family: Arial, sans-serif;
            font-size: 80px;
            color: #333;
        }

        .clock .separator {
            font-family: Arial, sans-serif;
            font-size: 80px;
            color: #333;
            margin: 0 10px;
        }

        .clock .hour {
            color: #e67e22;
        }

        .clock .minute {
            color: #2ecc71;
        }

        .clock .second {
            color: #3498db;
        }
    </style>
</head>
<body>
    <div class="clock">
        <div class="hour digit"></div>
        <div class="separator">:</div>
        <div class="minute digit"></div>
        <div class="separator">:</div>
        <div class="second digit"></div>
    </div>
    <script>
        function updateClock() {
            const currentTime = new Date();
            const hourDigit = document.querySelector('.hour');
            const minuteDigit = document.querySelector('.minute');
            const secondDigit = document.querySelector('.second');

            hourDigit.textContent = currentTime.getHours().toString().padStart(2, '0');
            minuteDigit.textContent = currentTime.getMinutes().toString().padStart(2, '0');
            secondDigit.textContent = currentTime.getSeconds().toString().padStart(2, '0');
        }

        setInterval(updateClock, 1000);
    </script>
</body>
</html>
```

## 3. 代理mysql

```shell
[root@localhost ~]# ./licent -l 127.0.0.1:3306
  _   _              _______    _____   _        _____   ______   _   _   _______
 | \ | |     /\     |__   __|  / ____| | |      |_   _| |  ____| | \ | | |__   __|
 |  \| |    /  \       | |    | |      | |        | |   | |__    |  \| |    | |
 | . ` |   / /\ \      | |    | |      | |        | |   |  __|   | . ` |    | |
 | |\  |  / ____ \     | |    | |____  | |____   _| |_  | |____  | |\  |    | |
 |_| \_| /_/    \_\    |_|     \_____| |______| |_____| |______| |_| \_|    |_|


[Client Running Successfully!]
[PublicAddress] pogf.com.cn
[TunnelAddress] pogf.com.cn:8008
[LocalAddress] 127.0.0.1:3306
[Conn Successfully]pogf.com.cn:8080
[receive KeepLive package] ping
[ClientInfoUID] 5
[VisitAddress] pogf.com.cn:60003
```

![](./imagesn-1.png)
