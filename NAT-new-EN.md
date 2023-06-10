# Inner net penetration

## 1. Concept

Intranet Penetration is a technology that allows devices or services located in a private network, such as a home or corporate network, to be accessed over a public network, such as the Internet. In a traditional network environment, devices in a private network are often not directly accessible from the public network because they are in the private IP address range and are restricted by Network Address Translation (NAT) and firewalls. Internal network penetration technology overcomes these limitations by various means, so that devices in the public network can directly access devices in the private network.

A popular example is camera access in a home network. Let's say you have a camera in your home and want to be able to access it remotely over the Internet, such as viewing live video from your phone or computer. However, the camera is connected to a private IP address in the home network, which is not directly accessible by the Internet. At this time, you can use the intranet penetration technology to transmit the data stream of the camera to your mobile phone or computer through the public network. By configuring the appropriate intranet penetration tools and settings, you can connect to the camera in the private network through the public network in different network environments to achieve remote access and monitoring.

![](D:\goworkplace\src\github.com\byteYuFan\Intranet-penetration\images\1-NAT概述.png)

Next, I will develop an `NAT` intranet penetration plug-in by GO language.

The overall program is divided into three main parts, `通用部分`, `客户端部分`, `服务器部分`

Below I will introduce the implementation of the whole program from the general part.

[项目github地址](https://github.com/byteYuFan/NAT)

## 2. The formulation of the agreement

### 2.1. Agreement for sending and receiving a package

We know that `TCP` the transmission of data is specific `字节流`, so when we send some information, it is possible that these packets will be stuck together, and the client receives information that should not appear, which we absolutely do not want to appear, so we define the relevant protocol of sending and receiving packets in this area to avoid these problems.

See this article for `TCP` the features of the byte stream (https://juejin.cn/post/7239996748319899703)

#### 1. Message interface definition


```go
// IMessage 将请求的消息封装到一个Message中，定义抽象的接口
type IMessage interface {
	//	GetMsgID 获取消息的ID
	GetMsgID() uint32
	// GetMsgDataLen 获取消息长度
	GetMsgDataLen() uint32
	// GetMsgData 获取消息内容
	GetMsgData() []byte
	//	SetMsgID 设置消息的ID
	SetMsgID(uint32)
	// SetMsgLen 设置消息长度
	SetMsgLen(uint32)
	// SetMsgData 设置消息内容
	SetMsgData([]byte)
}
```

This interface `IMessage` defines the method that encapsulates the request message to `Message`, and provides some abstract methods to manipulate the message. The following is a description of the interface:

-  `GetMsgID() uint32` : Get the ID of the message and return the message ID of a `uint32` type.
-  `GetMsgDataLen() uint32` : Get the length of the message and return the message length of a `uint32` type.
-  `GetMsgData() []byte` : Gets the content of the message, returning a `[]byte` type of message content.
-  `SetMsgID(uint32)` : Sets the ID of the message. Takes a `uint32` parameter of type to set the ID of the message.
-  `SetMsgLen(uint32)` : Sets the length of the message. Takes a `uint32` parameter of type to set the length of the message.
-  `SetMsgData([]byte)` : Sets the content of the message. Takes a `[]byte` parameter of type to set the content of the message.

#### 2. Message interface implementation


```go
// Message 该结构体为消息类型，其中包括消息的ID编号，
// 消息的长度 消息的内容
// 其中 ID 和 DataLen 为固定长度大小 为 4 + 4 = 8 字节
// Data 为消息的具体内容 其中DataLen为Data的大小
type Message struct {
	// ID 消息ID
	ID uint32
	// DataLen 消息长度
	DataLen uint32
	// Data 消息的内容
	Data []byte
}
// NewMsgPackage 新建一个Message消息的包 传入的参数为消息的type-id 消息的实际内容
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		ID:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// GetMsgID 获取消息的ID
func (m *Message) GetMsgID() uint32 {
	return m.ID
}

// GetMsgDataLen 获取消息长度
func (m *Message) GetMsgDataLen() uint32 {
	return m.DataLen
}

// GetMsgData 获取消息内容
func (m *Message) GetMsgData() []byte {
	return m.Data
}

// SetMsgID 设置消息的ID
func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

// SetMsgLen 设置消息长度
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}

// SetMsgData 设置消息内容
func (m *Message) SetMsgData(data []byte) {
	m.Data = data
}
```

-  `ID uint32` : ID number of the message, occupying a fixed length of 4 bytes.
-  `DataLen uint32` : The length of the message, occupying 4 bytes of the fixed length.
-  `Data []byte` : The specific content of the message, and the length `DataLen` is specified by.

In addition, the code provides methods to manipulate `Message` the structure:

-  `NewMsgPackage(id uint32, data []byte) *Message` : Create a new `Message` message object with the ID and actual content of the incoming message as parameters.
-  `GetMsgID() uint32` : Get the ID of the message.
-  `GetMsgDataLen() uint32` Gets the length of the message.
-  `GetMsgData() []byte` Gets the content of the message.
-  `SetMsgID(id uint32)` : Set the ID of the message.
-  `SetMsgLen(len uint32)` : Set the length of the message.
-  `SetMsgData(data []byte)` : Set the content of the message.

#### 3. Message passing interface


```go
// IDataPack 抽象层封装包解决TCP黏包问题的拆包封装包的模块
// 针对Message进行TLV格式的封装
// 针对Message进行TLV格式的拆包
// 先读取固定长度的head-->消息的长度和消息的类型
// 在根据消息内容的长度，在读取内容
// 直接面向TCP连接的数据流 TCP stream
type IDataPack interface {
	// GetHeadLen 获取包的头长度
	GetHeadLen() uint32
	// Pack 封包
	Pack(msg IMessage) ([]byte, error)
	// Unpack 拆包
	Unpack([]byte) (IMessage, error)
}
```

The above code defines an `IDataPack` interface named, which is used to encapsulate the module of unpacking and encapsulating packets at the abstraction layer to solve the TCP sticky packet problem. The interface provides the following methods:

-  `GetHeadLen() uint32` : Get the header length of the package. Used to determine the fixed length of the message header.
- The `Pack(msg IMessage) ([]byte, error)` encapsulation method encapsulates the message into the form of a byte stream for sending.
-  `Unpack([]byte) (IMessage, error)` : The unpacking method parses the received byte stream and restores it to a message object.

The implementation class `Message` of the interface encapsulates and unpacks the message in TLV format. The TLV format means that the message header contains the type (Type), the length (Length), and the actual message content (Value).

The specific packaging process is as follows:

1. According to the content of the message object, a message header is constructed, which contains the length and type information of the message.
2. The message header and the message content are spliced together according to certain rules to form complete packet data.

The specific unpacking process is as follows:

1. First, read the header data with fixed length to obtain the length and type information of the message.
2. And reading the message content with the corresponding length according to the message length information.
3. According to the read header and content, it is restored to a complete message object.

Through the implementation `IDataPack` of the specific class of the interface, the data stream in the TCP connection can be packaged and unpackaged to solve the problem of TCP sticky packets and ensure the complete transmission of messages.

#### 4. An entity that deliver a message


```go

// DataPackage 封包解包的结构体
type DataPackage struct {
}

// NewDataPackage 创建一个封包拆包的实例
func NewDataPackage() *DataPackage {
	return &DataPackage{}
}

// GetHeadLen 获取包头的长度 根据我们的协议定义直接返回8就可以了
func (dp *DataPackage) GetHeadLen() uint32 {
	return uint32(8)
}

// Pack 将 ninterfance.IMessage 类型的结构封装为字节流的形式
// 字节流形式 [ 数据长度 + ID + 真实数据 ]
func (dp *DataPackage) Pack(msg ninterfance.IMessage) ([]byte, error) {
	// 创建一个字节流的缓存，将msg的信息一步一步的填充到里面去
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgDataLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (dp *DataPackage) Unpack(data []byte) (ninterfance.IMessage, error) {
	// 创建一个从data里面读取的ioReader
	dataBuffer := bytes.NewBuffer(data)
	msg := &Message{}
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	return msg, nil
}

```

This code implements the function of Pack and Unpack of a data packet, which is used to encapsulate a Message into a byte stream and parse the byte stream into a message object.

1.  `DataPackage` A struct is an implementation of packing and unpacking, in which no state information is stored.
2. The instance () function `NewDataPackage()` creates an `DataPackage` instance and returns a pointer.
3.  `GetHeadLen()` Method is used to obtain the length of the packet header. Here, a fixed value of 8 is directly returned, indicating that the length of the packet header is 8 bytes.
4.  `Pack(msg ninterfance.IMessage) ([]byte, error)` Method encapsulates the passed `IMessage` object as a byte stream. The specific implementation is as follows:
   1. Create a byte stream cache `dataBuff`.
   2. The use `binary.Write` method writes `dataBuff` the data length, ID and real data of the message in turn according to the big-endian byte order.
   3. Finally, the encapsulated byte stream is returned `dataBuff.Bytes()`.
5.  `Unpack(data []byte) (ninterfance.IMessage, error)` Method is used to parse a byte stream into a message object. The specific implementation is as follows:
   1. Create a `data` to read `dataBuffer` from.
   2. Creates an empty `Message` object `msg`.
   3. The usage `binary.Read` method reads the data length and ID from `dataBuffer` the according to the big-endian byte order, and assigns them to `msg` the corresponding fields.
   4. Finally, the resolved `msg` object is returned.

### 2.1. Message type protocol


```go
const (
	USER_AUTHENTICATION_SUCCESSFULLY = 10001
)
const (
	NEW_CONNECTION     = 20001
	USER_INFORMATION   = 20002
	KEEP_ALIVE         = 20003
	CONNECTION_IF_FULL = 20004
)

const (
	USER_NOT_EXIST     = 30001
	USER_ALREADY_EXIST = 30002
	USER_EXPIRED       = 30003
	Password_INCORRET  = 30004
	AUTH_FAIL          = 30005
)

const (
	USER_REQUEST_AUTH = 60001
)

```

1. In the first constant block, `USER_AUTHENTICATION_SUCCESSFULLY` the value of is a 10001, indicating successful user authentication.
2. In the second constant block, the following constants are defined:
   -  `NEW_CONNECTION` : The value is a 20001, indicating a new connection.
   -  `USER_INFORMATION` : a value of 2000 2 represents user information.
   -  `KEEP_ALIVE` : The value is the 20003, indicating that the connection is maintained.
   -  `CONNECTION_IF_FULL` : The value is a 20004, indicating that the connection is full.
3. In the third constant block, the following constants are defined:
   -  `USER_NOT_EXIST` : a value of 3000 1 indicates that the user does not exist.
   -  `USER_ALREADY_EXIST` : a value of 3000 2 indicates that the user already exists.
   -  `USER_EXPIRED` : The value is a 30003 indicating that the user has expired.
   -  `PASSWORD_INCORRET` : The value is a 30004 indicating that the password is incorrect.
   -  `AUTH_FAIL` : a value of 3000 5 indicates that authentication failed.
4. In the fourth constant block, the following constants are defined:
   -  `USER_REQUEST_AUTH` : The value is a 60001 indicating that the user requested authentication.

In the subsequent process, more systematic design and related business will be constantly updated.


```go
var ProtocolMap map[uint32]interface{}

func init() {
	ProtocolMap = make(map[uint32]interface{})
	// 添加创建协议
	ProtocolMap[NEW_CONNECTION] = "NEW_CONNECTION"
	ProtocolMap[USER_INFORMATION] = ClientConnInfo{}
	ProtocolMap[KEEP_ALIVE] = "ping"
	ProtocolMap[CONNECTION_IF_FULL] = "the-connection-is-full."

	ProtocolMap[USER_ALREADY_EXIST] = "user already exist!"
	ProtocolMap[USER_NOT_EXIST] = "user not exist."
	ProtocolMap[USER_EXPIRED] = "user not expired."
	ProtocolMap[PASSWORD_INCORRET] = "The password is incorrect."

	ProtocolMap[USER_AUTHENTICATION_SUCCESSFULLY] = "Authentication successful!"
}
```

This code defines a `ProtocolMap` variable named, which is a map type that stores the mapping of the message protocol.

In the `init()` function, first use `make()` the function to create an empty `ProtocolMap` one, and then associate different message protocols with the corresponding values by means of key-value pairs.

The specific mapping relationship is as follows:

- The `NEW_CONNECTION` corresponding value is a string `"NEW_CONNECTION"` indicating the protocol of the new connection.
- The `USER_INFORMATION` corresponding value is `ClientConnInfo{}` the protocol representing the user information, which `ClientConnInfo` is a struct type.
- The `KEEP_ALIVE` corresponding value is a string `"ping"` indicating the protocol that maintains the connection.
- The `CONNECTION_IF_FULL` corresponding value is a string `"the-connection-is-full."` that represents the protocol with the full connection.

- The `USER_ALREADY_EXIST` corresponding value is a string `"user already exist!"` indicating that the user already exists.
- The `USER_NOT_EXIST` corresponding value is a string `"user not exist."` indicating that the user does not exist.
- The `USER_EXPIRED` corresponding value is a string `"user not expired."` indicating that the user is not expired.
- The `PASSWORD_INCORRET` corresponding value is a string `"The password is incorrect."` indicating that the password is incorrect
- The `USER_AUTHENTICATION_SUCCESSFULLY` corresponding value is a string `"Authentication successful!"`, indicating that `ProtocolMap` the corresponding message or status information can be obtained according to the specific protocol value, which is convenient for handling different protocol situations in the code.

## 3. General module configuration

### 3.1. Exchange of data letter

This step is the core process of the whole intranet penetration, no matter how the whole program is written later, it is very important here. In this step, we need to encapsulate a function to exchange data, that is, to forward data from one `tcpConn` to `tcpConn` another. Before implementing these contents, let's first look at `go` a function `io.Copy` provided:


```go
func Copy(dst Writer, src Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}
```

Parameter description:

-  `dst` Is the destination `io.Writer` interface that receives the copied data.
-  `src` Is the source `io.Reader` interface that provides the data to be copied.

Return value:

-  `written` Is the number of bytes copied.
-  `err` Is an error that can occur. If the copy operation is successful, the value is `nil`.

The `io.Copy` function continually reads data from the source `src` and writes it to the target `dst` until the source `src` 's data ends or an error occurs. It automates the process of buffering and replicating data, simplifying the operation of data copying.

After reading this function, do we immediately understand that the exchange of data is like this `Easy`!

 `TCPConn` The type implements the `Write(p []byte) (n int, err error)` and `Read(p []byte) (n int, err error)` method, so the type also implements the `Writer` and `Reader` interface, so the function can be called directly. I like the go language feature very much.

So we have the encapsulated function:


```go
// SwapConnDataEachOther 通讯双方相互交换数据
func SwapConnDataEachOther(local, remote *net.TCPConn) {
	go swapConnData(local, remote)
	go swapConnData(remote, local)
}

// SwapConnData 这个函数是交换两个连接数据的函数
func swapConnData(local, remote *net.TCPConn) {
	// 关闭本地和远程连接通道
	defer local.Close()
	defer remote.Close()
	// 将remote的数据拷贝到local里面
	_, err := io.Copy(local, remote)
	if err != nil {
		return
	}
}
```

### 3.2. Create a TCP listener and connection


```go
func CreateTCPListener(addr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	return tcpListener, nil
}

// CreateTCPConn 连接指定的TCP
func CreateTCPConn(addr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return tcpConn, nil
}

```

These two functions are relatively simple, I will not repeat, `CreateTCPListener` is a new `tcp` one `listener`, `CreateTCPConn` is a new `tcp` one `connecter`, very simple.

### 3.2. Client connection information module

This module is an agreement between the client and the server on how to transmit the user's data, how to convert the user's information into a byte stream, and how to unpack the byte stream into user information. It can also be called a protocol.


```go
// ClientConnInfo 客户端连接信息
type ClientConnInfo struct {
	UID  int64
	Port int32
}

// NewClientConnInstance 新建一个实体
func NewClientConnInstance(id int64, port int32) *ClientConnInfo {
	return &ClientConnInfo{
		UID:  id,
		Port: port,
	}
}

// ToBytes 将 ClientConnInfo 结构体转换为字节流
func (info *ClientConnInfo) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 使用 binary.Write 将字段逐个写入字节流
	if err := binary.Write(buf, binary.BigEndian, info.UID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, info.Port); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// FromBytes 从字节流中恢复 ClientConnInfo 结构体
func (info *ClientConnInfo) FromBytes(data []byte) error {
	buf := bytes.NewReader(data)
	// 使用 binary.Read 从字节流中读取字段值
	if err := binary.Read(buf, binary.BigEndian, &info.UID); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &info.Port); err != nil {
		return err
	}
	return nil
}
```

Unit tests:


```go
func TestClientConnInfo_FromBytes(t *testing.T) {
	ci := ClientConnInfo{
		UID:  1,
		Port: 8080,
	}
	data, err := ci.ToBytes()
	if err != nil {
		fmt.Println("[ToBytes Err]", err)
	} else {
		fmt.Println("[ToBytes Successfully]", data)
	}
	nci := new(ClientConnInfo)
	err = nci.FromBytes(data)
	if err != nil {
		fmt.Println("[FromBytes Err]", err)
	} else {
		fmt.Println("[FromBytes Successfully]", nci)
	}
}

```


```shell
=== RUN   TestClientConnInfo_FromBytes
[ToBytes Successfully] [0 0 0 0 0 0 0 1 0 0 31 144]
[FromBytes Successfully] &{1 8080}
--- PASS: TestClientConnInfo_FromBytes (0.00s)
```

### 3.3. User information module

These are similar, and you will consider encapsulating these as an interface later.


```go
// UserInfo 用户信息模块
type UserInfo struct {
	// UserName
	UserName string
	// Password
	Password string
	// ExpireTime
	ExpireTime time.Time
}

// NewUserInfoInstance 新建一个实体
func NewUserInfoInstance(username, password string) *UserInfo {
	return &UserInfo{
		UserName: username,
		Password: password,
	}
}

func (info *UserInfo) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 将用户名长度编码到字节流中
	userNameLen := len(info.UserName)
	if err := binary.Write(buf, binary.BigEndian, int32(userNameLen)); err != nil {
		return nil, err
	}

	// 将用户名内容编码到字节流中
	if err := binary.Write(buf, binary.BigEndian, []byte(info.UserName)); err != nil {
		return nil, err
	}

	// 将密码长度编码到字节流中
	passwordLen := len(info.Password)
	if err := binary.Write(buf, binary.BigEndian, int32(passwordLen)); err != nil {
		return nil, err
	}

	// 将密码内容编码到字节流中
	if err := binary.Write(buf, binary.BigEndian, []byte(info.Password)); err != nil {
		return nil, err
	}

	// 将过期时间编码到字节流中
	if err := binary.Write(buf, binary.BigEndian, info.ExpireTime.Unix()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (info *UserInfo) FromBytes(data []byte) error {
	buf := bytes.NewReader(data)

	// 从字节流中读取用户名长度
	var userNameLen int32
	if err := binary.Read(buf, binary.BigEndian, &userNameLen); err != nil {
		return err
	}

	// 从字节流中读取用户名内容，并根据长度进行截取
	userNameBytes := make([]byte, userNameLen)
	if err := binary.Read(buf, binary.BigEndian, userNameBytes); err != nil {
		return err
	}
	info.UserName = string(userNameBytes)

	// 从字节流中读取密码长度
	var passwordLen int32
	if err := binary.Read(buf, binary.BigEndian, &passwordLen); err != nil {
		return err
	}

	// 从字节流中读取密码内容，并根据长度进行截取
	passwordBytes := make([]byte, passwordLen)
	if err := binary.Read(buf, binary.BigEndian, passwordBytes); err != nil {
		return err
	}
	info.Password = string(passwordBytes)

	// 从字节流中读取过期时间
	var expireTimeUnix int64
	if err := binary.Read(buf, binary.BigEndian, &expireTimeUnix); err != nil {
		return err
	}
	info.ExpireTime = time.Unix(expireTimeUnix, 0)

	return nil
}

```

This code defines a `UserInfo` structure that represents the user information module. It contains fields for user name, password, and expiration time.

 `NewUserInfoInstance` Function is a constructor that creates `UserInfo` an instance of a struct and initializes the user name and password.

 `ToBytes` Method `UserInfo` encodes the fields of a struct as a stream of bytes. First, create a `bytes.Buffer` buffer. The length of the user name is then written to the buffer in `int32` big-endian encoding of type. Next, the contents of the username are written to the buffer as a byte stream. The password length and password contents are then written to the buffer. Finally, the UNIX timestamp of the expiration time is written to the buffer in `int64` big-endian encoding of type. Finally, the byte stream in the buffer is returned.

 `FromBytes` Method to decode a byte stream into `UserInfo` the fields of a struct. First, create a `bytes.Reader` reader to read the byte stream. The username length is then read from the reader and a byte slice is created based on the length. Next, the user name content is read from the reader and converted to a string to assign to the `UserName` field. The password length is then read from the reader and a byte slice is created based on the length. Next, the password content is read from the reader and converted to a string and assigned to the `Password` field. Finally, the UNIX timestamp of the expiration time is read from the reader and converted to a `time.Time` type to assign to the `ExpireTime` field using `time.Unix` the function.

Therefore, the mutual conversion between the structure body and the byte stream can be realized `UserInfo` through `ToBytes` the and `FromBytes` methods, and the structure body and the byte stream can be conveniently used in network transmission or storage.

Unit tests:


```go
func TestNewUserInfoInstance(t *testing.T) {
	ui := NewUserInfoInstance("wyfld", "yfw123456789")
	ui.ExpireTime = time.Now()
	d, err := ui.ToBytes()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d)
	un := new(UserInfo)
	un.FromBytes(d)
	fmt.Println(un)
}
```


```shell
=== RUN   TestNewUserInfoInstance
[0 0 0 5 119 121 102 108 100 0 0 0 12 121 102 119 49 50 51 52 53 54 55 56 57 0 0 0 0 100 132 113 19]
&{wyfld yfw123456789 2023-06-10 20:48:19 +0800 CST}
--- PASS: TestNewUserInfoInstance (0.02s)
PASS
```

### 3.4. Send information module


```go
// SendAndReceiveInstance 实体
type SendAndReceiveInstance struct {
	Conn *net.TCPConn
}

// NewSendAndReceiveInstance 新建一个控制层实例对象
func NewSendAndReceiveInstance(conn *net.TCPConn) *SendAndReceiveInstance {
	return &SendAndReceiveInstance{
		Conn: conn,
	}
}

// SendDataToClient 向客户端发送消息，此处应该指定协议
func (csi *SendAndReceiveInstance) SendDataToClient(dataType uint32, msg []byte) (int, error) {
	msgInstance := NewMsgPackage(dataType, msg)
	pkgInstance := NewDataPackage()
	dataStream, err := pkgInstance.Pack(msgInstance)
	if err != nil {
		return 0, err
	}
	count, err := csi.Conn.Write(dataStream)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (csi *SendAndReceiveInstance) ReadHeadDataFromClient() (ninterfance.IMessage, error) {
	// 根据我们协议规定的内容每次数据来的时候，先读取头部的8个字节的数据
	headData := make([]byte, 8)
	// 从Conn中读取数据到headData中去
	if _, err := io.ReadFull(csi.Conn, headData); err != nil {
		return nil, err
	}
	// 先创建一个解包的实例
	dp := NewDataPackage()
	// 解封装这个包
	return dp.Unpack(headData)
}
func (csi *SendAndReceiveInstance) ReadRealDataFromClient(msg ninterfance.IMessage) (ninterfance.IMessage, error) {
	if msg.GetMsgDataLen() < 0 {
		return msg, nil
	}
	// 新建一个data，长度为msg头部的长度
	data := make([]byte, msg.GetMsgDataLen())
	if _, err := io.ReadFull(csi.Conn, data); err != nil {
		return nil, err
	}
	msg.SetMsgData(data)
	return msg, nil
}
```

This code defines a structure named `SendAndReceiveInstance` that represents an instance object that sends and receives messages. It contains a `Conn` field that stores the TCP connection to the client.

 `NewSendAndReceiveInstance` Function is a constructor that creates `SendAndReceiveInstance` an instance of a struct and initializes `Conn` a field.

 `SendDataToClient` Method is used to send a message to the client. It takes a `dataType` parameter indicating the type of the message and a `msg` parameter indicating the content of the message to be sent. First, use the `NewMsgPackage` function to create a message package instance that contains the message type and the message content. Then, use the `NewDataPackage` function to create an unpacking instance. Next, the message instance is packaged as a byte stream and the packaged data stream is obtained. Finally, the method using `Conn` the `Write` field writes the data stream to the connection and returns the number of bytes written and any errors that may have occurred.

 `ReadHeadDataFromClient` Method is used to read the header data from the client. According to the protocol, every time data arrives, it first needs to read 8 bytes of header data. This method creates a byte slice `headData` of length 8 and then uses the `io.ReadFull` function to read the data from `Conn` the field and store it in `headData`. Next, create an instance `dp` of packet parsing and unpack `headData` using its `Unpack` methods. Finally, the unpacked message and any errors that may have occurred are returned.

 `ReadRealDataFromClient` Method is used to read real data from the client. According to the protocol, the length of the real data is recorded in the message header, so the corresponding data needs to be read according to the length. First, check the data length of the message. If the length is less than 0, return the message directly. Otherwise, create a byte slice `data` based on the data length of the message, and then use the `io.ReadFull` function to read the data from `Conn` the field and store it in `data`. Finally, the read data is set to the data part of the message, and the message and possible errors are returned.

These methods are combined to realize the function of sending and receiving messages from the client. Through the process of encapsulation and unpacking, the messages are converted into byte streams for transmission, and then parsed and processed at the receiving end.

## 4. Client

### 4.1. Configuration file interpretation


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

The client configuration information includes:

-  `Name` : The name of the client, here "Client-NAT".
-  `PublicServerAddr` : The domain name of the public network server, used to connect to the public network server.
-  `TunnelServerAddr` : The tunnel port of the public network server, used to establish a tunnel connection.
-  `ControllerAddr` : The control port of the public network server is used to control and manage the client.
-  `LocalServerAddr` : The port of the local service where the client will provide the local service.

Authentication information includes:

-  `Username` : Username of the user, used for authentication.
-  `Password` : Password of the user, used for authentication.

Through these configuration and authentication information, the client can connect to the public network server and authenticate to obtain access rights. The client can then use the tunnel connection, the control port, and the local service port for network communication and service provisioning.

### 4.2. Configuration file read

The framework is used `viper` to read the contents of the configuration file.


```go
type ParseConfigFromYML struct {
	ViperInstance *viper.Viper
}

func ParseFile(configFileName string) *ParseConfigFromYML {
	pcfy := new(ParseConfigFromYML)
	pcfy.ViperInstance = viper.New()
	pcfy.ViperInstance.SetConfigName(configFileName)
	pcfy.ViperInstance.AddConfigPath(".")
	pcfy.ViperInstance.AddConfigPath("./config")
	pcfy.ViperInstance.AddConfigPath("../config")
	pcfy.ViperInstance.AddConfigPath("../../config")
	pcfy.ViperInstance.SetConfigType("yml")
	err := pcfy.ViperInstance.ReadInConfig()
	if err != nil {
		fmt.Println("[ReadInConfig]")
		return nil
	}
	return pcfy
}
```


```go
type objectConfigData struct {
	// Name 客户端名称
	Name string
	// LocalServerAddr 本地服务端地址
	LocalServerAddr string
	// TunnelServerAddr 隧道地址用于交换数据
	TunnelServerAddr string
	// PublicServerAddr 公网服务器地址
	PublicServerAddr string
	// ControllerAddr 服务器控制端地址
	ControllerAddr string
	// UserName 登录用户名
	UserName string
	// Password 密码
	Password string
}
```


```go
func initConfig() {
	objectConfig = new(objectConfigData)
	config := utils.ParseFile("client.yml")
	viper := config.ViperInstance
	objectConfig.Name = viper.GetString("Client.Name")
	objectConfig.PublicServerAddr = viper.GetString("Client.PublicServerAddr")
	objectConfig.ControllerAddr = viper.GetString("Client.ControllerAddr")
	objectConfig.LocalServerAddr = viper.GetString("Client.LocalServerAddr")
	objectConfig.TunnelServerAddr = viper.GetString("Client.TunnelServerAddr")
	objectConfig.UserName = viper.GetString("Auth.Username")
	objectConfig.Password = viper.GetString("Auth.Password")
}
```

1.  `objectConfig` Is a global variable used to store an object of configuration data of type `objectConfigData`.
2.  `config := utils.ParseFile("client.yml")` This line of code uses a `utils.ParseFile` function to parse a configuration file named client. Yml "and assign the result to a variable `config`.
3.  `viper := config.ViperInstance` This line of code gets the configuration parser instance, which uses a configuration parser library called `Viper`.
4. Next, the code gets the values of each configuration item from the configuration file through `viper.GetString` the method and assigns them to `objectConfig` the corresponding fields in.
   - Get the value of the Client. Name " `objectConfig.Name = viper.GetString("Client.Name")` from the configuration file and assign it to `objectConfig` the `Name` field.
   -  `objectConfig.PublicServerAddr = viper.GetString("Client.PublicServerAddr")` Get the value of the Client. Public ServerAddr "and assign it to `objectConfig` the `PublicServerAddr` field.
   -  `objectConfig.ControllerAddr = viper.GetString("Client.ControllerAddr")` Get the value of the Client. ControllerAddr "and assign it to `objectConfig` the `ControllerAddr` field.
   -  `objectConfig.LocalServerAddr = viper.GetString("Client.LocalServerAddr")` Get the value of the Client. LocalServerAddr "and assign it to `objectConfig` the `LocalServerAddr` field.
   -  `objectConfig.TunnelServerAddr = viper.GetString("Client.TunnelServerAddr")` Get the value of the Client. TunnelServerAddr "and assign it to `objectConfig` the `TunnelServerAddr` field.
   -  `objectConfig.UserName = viper.GetString("Auth.Username")` Get the value of the Auth. Username "and assign it to `objectConfig` the `UserName` field.
   -  `objectConfig.Password = viper.GetString("Auth.Password")` Get the value of the Auth. Password "and assign it to `objectConfig` the `Password` field.

The purpose of this code is to read the values from the configuration file into the `objectConfig` object so that other parts of the code can use the configuration data.

### 4.3. Command Line Controller

The command line is done with a `cobra` frame.

 `Cobra` Is an application framework for the CLI interface pattern, and is the command-line tool that generates it. Users can quickly see how the binary is used by using the help method. Command: generally represents an action, that is, a running binary command service. You can also have children commands.


```go

func initCobra() {
	object = new(objectConfigData)
	rootCmd.Flags().StringVarP(&object.Name, "name", "n", "", "Client name")
	rootCmd.Flags().StringVarP(&object.LocalServerAddr, "local-server-addr", "l", "", "The address of the local web server program")
	rootCmd.Flags().StringVarP(&object.TunnelServerAddr, "tunnel-server-addr", "t", "", "The address of the tunnel server used to connect the local and public networks")
	rootCmd.Flags().StringVarP(&object.ControllerAddr, "controller-addr", "c", "", "The address of the controller channel used to send controller messages to the client")
	rootCmd.Flags().StringVarP(&object.PublicServerAddr, "public-server-addr", "s", "", "The address of the public server used for accessing the inner web server")
	rootCmd.Flags().StringVarP(&object.UserName, "username", "u", "", "the name for auth the server.")
	rootCmd.Flags().StringVarP(&object.Password, "password", "P", "", "the password for auth the server.")
	// 添加其他字段...
}

func exchange() {
	if object.Name != "" {
		objectConfig.Name = object.Name
	}
	if object.TunnelServerAddr != "" {
		objectConfig.TunnelServerAddr = object.TunnelServerAddr
	}
	if object.ControllerAddr != "" {
		objectConfig.ControllerAddr = object.ControllerAddr
	}
	if object.PublicServerAddr != "" {
		objectConfig.PublicServerAddr = object.PublicServerAddr
	}
	if object.LocalServerAddr != "" {
		objectConfig.LocalServerAddr = object.LocalServerAddr
	}
	if object.UserName != "" {
		objectConfig.UserName = object.UserName
	}
	if object.Password != "" {
		objectConfig.Password = object.Password
	}
}

```

This code uses the Cobra library to handle command-line arguments and assign the values of the arguments to `objectConfigData` objects of type.

First, in the `initCobra` function, `rootCmd.Flags().StringVarP` the method defines the corresponding command-line flags for each parameter and specifies which field in the object the parameter value should be assigned to `object`. For example, `rootCmd.Flags().StringVarP(&object.Name, "name", "n", "", "Client name")` to represent a field to `object` `Name` which the value of a command line parameter `--name` or `-n` is assigned.

Then, in the `exchange` function, it is decided whether to update `objectConfig` the corresponding field value of the object according to `object` whether each field value in the object is null or not. If `object` a field of an object is not null, its value is assigned to `objectConfig` the corresponding field of the object.

The purpose of this code is to swap the values of the command-line arguments into `objectConfig` an object for subsequent use with the configuration data. Command-line arguments can be used to override values in the configuration file, providing a flexible way to configure.

### 4.4. main

The logic for the client is as follows:

1. Fill in the global object `objectConfig` correctly according to the configuration file and command line parameters
2. Print Client-Related Information
3. Connect to the public network server through the server control interface in the configuration file
4. Nding an authentication request to the server
5. Open a big loop, continuously receive the messages sent by the user, and make different responses according to the relevant content specified in the protocol.

Next, I will introduce the response to receiving different data:

#### 1. Received authentication failed


```go
if msg.GetMsgID() == network.AUTH_FAIL {
				fmt.Println("[auth  fail]", "认证失败")
				break
}
```

If the message type received by the client is `AUTH_FAIL`, it means that the user does not have permission to access the server, so it exits directly. You can expand the content here.

#### 2. New connection received


```go
case network.NEW_CONNECTION:
				processNewConnection(msg.GetMsgData())
func processNewConnection(data []byte) {
	// TODO 目前从服务端发送来的信息没有进行处理，后续考虑进行处理
	go connectLocalAndTunnel()
}

```

When the client receives `NEW_CONNECTION` the message, it will establish a tunnel between the intranet server and the public network connection for data exchange, which will be described in detail later.

#### 3. User information received


```go
case network.USER_INFORMATION:
				err := processUserInfo(msg.GetMsgData())
				if err != nil {
					fmt.Println("[User Info]", err)
					continue
				}
```

When the client receives `USER_INFORMATION` the message, it means that the server has configured the port number for the intranet server, and the external user can access the internal service program according to this port.

#### 4. Heartbeat packet received


```go
case network.KEEP_ALIVE:
				processKeepLive(msg.GetMsgData())
func processKeepLive(data []byte) {
	// TODO 目前只简简单单接收服务端发来的请求，简单的打印一下
	fmt.Println("[receive KeepLive package]", string(data))
}
```

When the client receives `KEEP_ALIVE` the heartbeat packet, it only prints it at present, and tells the server that I still have it.

#### 5. Connection full message received


```go
case network.CONNECTION_IF_FULL:
				processConnIsFull(msg.GetMsgData())
				break receiveLoop
			}
```

When the client receives `CONNECTION_IS_FULL` it, it `break` loops until it can receive the connection.

### 4.5. Certification


```go
// authTheServer 向服务器发送认证消息
func authTheServer(conn *net.TCPConn) error {
	// 新建一个数据结构体
	ui := network.NewUserInfoInstance(objectConfig.UserName, objectConfig.Password)
	byteStream, err := ui.ToBytes()
	if err != nil {
		return err
	}
	nsi := instance.NewSendAndReceiveInstance(conn)
	_, err = nsi.SendDataToClient(network.USER_REQUEST_AUTH, byteStream)
	if err != nil {
		return err
	}
	return nil
}

```

First, an `UserInfo` instance `ui` is created by calling `network.NewUserInfoInstance` the function, which contains the username and password required for authentication, obtained from global variables `objectConfig`.

A method is then called `ui.ToBytes` to convert the `ui` instance into a byte stream for transmission over the network. If an error occurs during the conversion, an error message is returned.

Next, create an `SendAndReceiveInstance` instance `nsi` by calling `instance.NewSendAndReceiveInstance` the function, which is used to send and receive messages. Use the `conn` parameter as `nsi` the connection object for.

The calling `nsi.SendDataToClient` method then sends an authentication message to the server. Set the message type to `network.USER_REQUEST_AUTH`, indicating that this is an authentication request message and the previously converted byte stream `byteStream` is sent as the message content. The method returns the number of bytes sent and the error that may have occurred.

Finally, the function returns any errors that may have occurred, and if there were no errors during the send, the return `nil` indicates that the send was successful.

Overall, what this code does is send an authentication message to the server, packing the username and password into a byte stream and sending it with `SendAndReceiveInstance` an instance.

### 4.6. Master function


```go
func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		// 只打印帮助信息，不执行命令
		rootCmd.SetArgs(os.Args[1:])
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		exchange()
		art()
		printRelationInformation()
		// 连接服务器控制接口
		controllerTCPConn, err := network.CreateTCPConn(objectConfig.ControllerAddr)
		if err != nil {
			log.Println("[CreateTCPConn]" + objectConfig.ControllerAddr + err.Error())
			return
		}
		fmt.Println("[Conn Successfully]" + objectConfig.ControllerAddr)
		err = authTheServer(controllerTCPConn)
		if err != nil {
			fmt.Println("[authTheServer]", err)
			return
		}
		nsi := instance.NewSendAndReceiveInstance(controllerTCPConn)
	receiveLoop:
		for {
			msg, err := nsi.ReadHeadDataFromClient()
			if err == io.EOF {
				break
			}
			if opErr, ok := err.(*net.OpError); ok {
				fmt.Println("[err]", err)
				if strings.Contains(opErr.Error(), "An existing connection was forcibly closed by the remote host") {
					// 远程主机关闭连接，退出连接处理循环
					fmt.Println("远程主机关闭连接")
					break
				}
			}
			if err != nil {
				fmt.Println("[err]", err)
				continue
			}
			msg, err = nsi.ReadRealDataFromClient(msg)
			if err != nil {
				fmt.Println("[readReal]", msg)
				continue
			}
			if msg.GetMsgID() == network.AUTH_FAIL {
				fmt.Println("[auth  fail]", "认证失败")
				break
			}
			switch msg.GetMsgID() {
			case network.NEW_CONNECTION:
				processNewConnection(msg.GetMsgData())
			case network.USER_INFORMATION:
				err := processUserInfo(msg.GetMsgData())
				if err != nil {
					fmt.Println("[User Info]", err)
					continue
				}
			case network.KEEP_ALIVE:
				processKeepLive(msg.GetMsgData())
			case network.CONNECTION_IF_FULL:
				processConnIsFull(msg.GetMsgData())
				break receiveLoop
			}

		}
	}
	fmt.Println("[客户端退出，欢迎您的使用。]", "GoodBye, Have a good time!!!")
}
```

## 5. Server side

### 5.1. Configuration file interpretation


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
  LogFilename: "server.log"
  StartAuth: true
Database:
  Username: "root"
  Password: "123456"
  Host: "127.0.0.1:3309"
  DBName: "NAT"
```


The contents of this configuration file describe the configuration of a server.

Server section:

- Name: Set the name of the server to "Server-NAT"
- ControllerAddr: The address of the specified controller is "0.0.0.0: 8007"
- TunnelAddr: The address of the specified tunnel is "0.0.0.0: 8008"
- VisitPort: Defines a list of access ports, including 60000, 60001, 60002, and 60003.
- TaskQueueNum: Set the number of task queues to 4
- TaskQueueBuff: Set the buffer size of the task queue to 32
- Max TCPCon nNum: Set the maximum number of TCP connections to 4
- MaxConnNum: Set the maximum number of connections to 256
- LogFilename: Sets the name of the log file to server. Log ".
- StartAuth: Enable authentication, set to true

Database section:

- Username: Username of the database is "root"
- Password: The password for the database is 123456.
- Host: The host address of the database is "127.0.0.1: 3309"
- DBName: The name of the database is "NAT"

To summarize, this configuration file describes a Server named "Server-nat" with a specified list of controller addresses, tunnel addresses, and access ports. The server is configured with parameters such as task queues, maximum connections, and log file names. In addition, the configuration file specifies the connection information to the database, including user name, password, host address, and database name.

### 5.2. Global configuration object


```go
// ObjectConfigData 全局配置的对象,里面存储着服务端所有的配置信息
type objectConfigData struct {
	// ServerName 服务端名称
	Name string
	// ControllerAddr 服务器控制端地址
	ControllerAddr string
	// TunnelAddr 隧道地址交换数据隧道
	TunnelAddr string
	// ExposePort 服务端向外暴露的端口号
	ExposePort []int
	// TaskQueueNum 任务队列的数量
	TaskQueueNum int32
	// TaskQueueBufferSize 缓冲区最大的数量
	TaskQueueBufferSize int32
	// MaxTCPConnNum  一次性最大处理的并发连接数量，等同于任务队列的大小和服务端暴露的端口号数量
	MaxTCPConnNum int32
	// MaxConnNum 整个系统所能接收到的并发数量 为工作队列大小和工作队列缓冲区之积
	MaxConnNum int32
	// 	LogFilename string 日志文件名称
	LogFilename string
	// StartAuth 是否开启认证功能
	StartAuth string
	// DB 如果开启认证功能就得从配置文件中读取相关的配置信息
	DB DataBase
}

// DataBase 数据库相关信息
type DataBase struct {
	Username string
	Password string
	Host     string
	DBName   string
}
```

This code segment defines a `objectConfigData` structure named that represents a global configuration information object that stores various configuration information for the server.

 `objectConfigData` The struct contains the following fields:

-  `Name` : Server name, used to identify the server.
-  `ControllerAddr` : Server control end address, indicating the address of the controller.
-  `TunnelAddr` : Tunnel address, the tunnel used to exchange data.
-  `ExposePort` : a list of port numbers to which the server is exposed, storing multiple integer values.
-  `TaskQueueNum` : The number of task queues, which represents the number of queues processing tasks concurrently.
-  `TaskQueueBufferSize` : The size of the buffer, which represents the maximum number of buffers for the task queue.
- The maximum number of concurrent connections processed `MaxTCPConnNum` at one time, which is the same as the size of the task queue and the number of port numbers exposed by the server.
- The amount of concurrency the `MaxConnNum` entire system can receive, which is the product of the work queue size and the work queue buffers.
-  `LogFilename` : Log File Name. Indicates the name of the log file.
-  `StartAuth` : Flag indicating whether the authentication function is enabled. The type is string.
-  `DB` : Database-related information, including user name, password, host address, and database name, stored in the `DataBase` struct.

Through this structure, you can easily manage and access the configuration information of the server.

### 5.3. viper and cobra


```go
func initConfig() {
	// 读取配置文件内容
	config := utils.ParseFile("server.yml")
	viper := config.ViperInstance
	viper.SetDefault("Server.Name", "Server-NAT")
	viper.SetDefault("Server.ControllerAddr", "0.0.0.0:8007")
	viper.SetDefault("Server.TunnelAddr", "0.0.0.0:8008")
	viper.SetDefault("Server.VisitPort", []uint16{60000, 60001, 60002, 60003})
	viper.SetDefault("Server.TaskQueueNum", 4)
	viper.SetDefault("Server.TaskQueueBuff", 32)
	viper.SetDefault("Server.MaxTCPConnNum", 4)
	viper.SetDefault("Server.MaxConnNum", 128)
	viper.SetDefault("Server.LogFilename", "server.log")
	viper.SetDefault("Server.StartAuth", true)
	// 读取配置值并存入 objectConfig
	objectConfig.Name = viper.GetString("Server.Name")
	objectConfig.ControllerAddr = viper.GetString("Server.ControllerAddr")
	objectConfig.TunnelAddr = viper.GetString("Server.TunnelAddr")
	objectConfig.ExposePort = viper.GetIntSlice("Server.VisitPort")
	objectConfig.TaskQueueNum = viper.GetInt32("Server.TaskQueueNum")
	objectConfig.MaxTCPConnNum = viper.GetInt32("Server.MaxTCPConnNum")
	objectConfig.TaskQueueBufferSize = viper.GetInt32("Server.TaskQueueBuff")
	objectConfig.MaxConnNum = viper.GetInt32("Server.MaxConnNum")
	objectConfig.LogFilename = viper.GetString("Server.LogFilename")
	objectConfig.StartAuth = viper.GetString("Server.StartAuth")
	objectConfig.DB.Username = viper.GetString("Database.Username")
	objectConfig.DB.Password = viper.GetString("Database.Password")
	objectConfig.DB.Host = viper.GetString("Database.Host")
	objectConfig.DB.DBName = viper.GetString("Database.DBName")
}

```

The main steps of the function are as follows:

1. Read the contents of the configuration file through the call `utils.ParseFile("server.yml")` and return a configuration object `config`.
2. Get the Viper instance of the configuration object and assign it to the variable `viper`.
3. Use the `viper.SetDefault()` method to set the default value for each CI in case the corresponding value is not specified in the configuration file.
4. The configuration values are read from the configuration object and stored in `objectConfig` the corresponding fields in the object.
   - Assigns `objectConfig.Name` the `Server.Name` value of.
   - Assigns `objectConfig.ControllerAddr` the `Server.ControllerAddr` value of.
   - Assigns `objectConfig.TunnelAddr` the `Server.TunnelAddr` value of.
   - Assigns `objectConfig.ExposePort` the `Server.VisitPort` value of.
   - Assigns `objectConfig.TaskQueueNum` the `Server.TaskQueueNum` value of.
   - Assigns `objectConfig.MaxTCPConnNum` the `Server.MaxTCPConnNum` value of.
   - Assigns `objectConfig.TaskQueueBufferSize` the `Server.TaskQueueBuff` value of.
   - Assigns `objectConfig.MaxConnNum` the `Server.MaxConnNum` value of.
   - Assigns `objectConfig.LogFilename` the `Server.LogFilename` value of.
   - Assigns `objectConfig.StartAuth` the `Server.StartAuth` value of.
   - Assigns `objectConfig.DB.Username` the `Database.Username` value of.
   - Assigns `objectConfig.DB.Password` the `Database.Password` value of.
   - Assigns `objectConfig.DB.Host` the `Database.Host` value of.
   - Assigns `objectConfig.DB.DBName` the `Database.DBName` value of.

Through this function, the configuration value can be read from the configuration file and stored in the global `objectConfig` object for subsequent use.


```go
func initCobra() {
	object = &objectConfigData{
		// 初始化对象的字段
	}

	// 将命令行参数与对象的字段绑定
	rootCmd.Flags().StringVarP(&object.Name, "name", "n", "", "Server name")
	rootCmd.Flags().StringVarP(&object.ControllerAddr, "controller-addr", "c", "", "Server controller address")
	rootCmd.Flags().StringVarP(&object.TunnelAddr, "tunnel-addr", "t", "", "Server tunnel address")
	rootCmd.Flags().IntSliceVarP(&object.ExposePort, "expose-port", "p", nil, "Server exposed ports")
	rootCmd.Flags().Int32VarP(&object.TaskQueueNum, "task-queue-num", "q", 0, "Task queue number")
	rootCmd.Flags().Int32VarP(&object.TaskQueueBufferSize, "task-queue-buffer-size", "b", 0, "Task queue buffer size")
	rootCmd.Flags().Int32VarP(&object.MaxTCPConnNum, "max-tcp-conn-num", "m", 0, "Maximum TCP connection number")
	rootCmd.Flags().Int32VarP(&object.MaxConnNum, "max-conn-num", "x", 0, "Maximum connection number")
	rootCmd.Flags().StringVarP(&object.LogFilename, "log-name", "l", "", "The name of the log.")
	rootCmd.Flags().StringVarP(&object.StartAuth, "start-auth", "a", "true", "This is the method that whether the server start the auth.")

	// 打印绑定后的对象

	// 将参数赋值给目标配置对象

	// 添加其他字段...
}

func exchange() {
	if object.Name != "" {
		objectConfig.Name = object.Name
	}
	if object.ControllerAddr != "" {
		objectConfig.ControllerAddr = object.ControllerAddr
	}
	if object.TunnelAddr != "" {
		objectConfig.TunnelAddr = object.TunnelAddr
	}
	if object.LogFilename != "" {
		objectConfig.LogFilename = object.LogFilename
	}
	if object.ExposePort != nil {
		objectConfig.ExposePort = object.ExposePort
	}
	if object.TaskQueueNum != 0 {
		objectConfig.TaskQueueNum = object.TaskQueueNum
	}
	if object.TaskQueueBufferSize != 0 {
		objectConfig.TaskQueueBufferSize = object.TaskQueueBufferSize
	}
	if object.MaxTCPConnNum != 0 {
		objectConfig.MaxTCPConnNum = object.MaxTCPConnNum
	}
	if object.MaxConnNum != 0 {
		objectConfig.MaxConnNum = object.MaxConnNum
	}
	if object.StartAuth != "true" {
		objectConfig.StartAuth = object.StartAuth
	}
}
```

### 5.4. User connection pool


```go
// UserConnInfo 用户连接信息,此处保存的是用户访问公网web所对应的那个接口
type UserConnInfo struct {
	// visit 用户访问web服务的时间
	visit time.Time
	// conn tcp连接句柄
	conn *net.TCPConn
}

// userConnPool 用户连接池
type userConnPool struct {
	// UserConnectionMap 连接池map，存放着用户的连接信息 key-时间戳 val userConnInfo
	UserConnectionMap map[string]*UserConnInfo
	// Mutex 读写锁，用来保证map的并发安全问题
	Mutex sync.RWMutex
}

var userConnPoolInstance *userConnPool

// NewUserConnPool 新建一个连接池对象
func NewUserConnPool() *userConnPool {
	return &userConnPool{
		UserConnectionMap: make(map[string]*UserConnInfo),
		Mutex:             sync.RWMutex{},
	}
}

// AddConnInfo 向连接池添加用户的信息
func (ucp *userConnPool) AddConnInfo(conn *net.TCPConn) {
	// 加写锁保护并发的安全性
	ucp.Mutex.Lock()
	defer ucp.Mutex.Unlock()
	nowTime := time.Now()
	uci := &UserConnInfo{
		visit: nowTime,
		conn:  conn,
	}
	ucp.UserConnectionMap[strconv.FormatInt(nowTime.UnixNano(), 10)] = uci
}
func initUserConnPool() {
	userConnPoolInstance = NewUserConnPool()
}

```

The above code defines a user connection pool `userConnPool` that is used to manage the connection information of users. The user connection information includes the time when the user accesses the public web ( `visit`) and the TCP connection handle ( `conn`).

Here is an explanation of the important parts in the code:

-  `UserConnInfo` The structure is used to store user connection information. It contains a `visit` field that represents the time the user accessed the web service and `conn` a field that represents the TCP connection handle.
-  `userConnPool` The struct represents the user connection pool. It contains a `UserConnectionMap` field that stores a mapping of user connection information, with the timestamp as the key and `UserConnInfo` the object as the value. In addition, it contains a `Mutex` field to ensure `UserConnectionMap` concurrency security for.
-  `NewUserConnPool` Is a factory function that creates and returns a new connection pool object. It initializes to `UserConnectionMap` an empty map.
-  `AddConnInfo` Method is used to add user connection information to the connection pool. It takes the current time as the timestamp, creates a new `UserConnInfo` object, and adds it to `UserConnectionMap` the, with the timestamp as the key.
-  `initUserConnPool` The function is used to initialize a global `userConnPoolInstance` variable, create a new connection pool object by calling `NewUserConnPool`, and assign a value to `userConnPoolInstance`.

The user connection pool is used to manage the user's connection information. When the user establishes a connection with the server, the connection information can be added to the connection pool through `AddConnInfo` the method for subsequent use and management. By using connection pooling, you can manage and control user connections more efficiently and provide concurrent, secure access.

### 5.5. Work Pool


```go
// Worker 真正干活的工人，数量和TCP MAX有关
type Worker struct {
	// ClientConn 客户端连接
	ClientConn *net.TCPConn
	// ServerListener 服务端监听端口
	ServerListener *net.TCPListener
	// Port 服务端对应端口
	Port int32
}

type Workers struct {
	Mutex        sync.RWMutex
	WorkerStatus map[int32]*Worker
}

// NewWorkers 新建workers
func NewWorkers() *Workers {
	return &Workers{
		Mutex:        sync.RWMutex{},
		WorkerStatus: make(map[int32]*Worker),
	}
}
func (workers *Workers) Add(port int32, w *Worker) {
	workers.Mutex.Lock()
	defer workers.Mutex.Unlock()
	serverInstance.Counter++
	workers.WorkerStatus[port] = w
	serverInstance.PortStatus[port] = true
}
func (workers *Workers) Remove(port int32) {
	workers.Mutex.Lock()
	defer workers.Mutex.Unlock()
	workers.WorkerStatus[port].ServerListener.Close()
	delete(workers.WorkerStatus, port)
	serverInstance.PortStatus[port] = false
}
func (workers *Workers) Get(port int32) *Worker {
	workers.Mutex.RLock()
	defer workers.Mutex.RUnlock()
	return workers.WorkerStatus[port]
}

func NewWorker(l *net.TCPListener, c *net.TCPConn, port int32) *Worker {
	return &Worker{
		ClientConn:     c,
		ServerListener: l,
		Port:           port,
	}
}

```

The above code defines two structures: `Worker` and `Workers`.

 `Worker` The struct represents a worker and contains the following fields:

- The TCP connection handle for the `ClientConn` client connection.
-  `ServerListener` : TCP listener for server-side listening port.
-  `Port` : The port number corresponding to the server.

 `Workers` A struct represents a group of workers and contains the following fields:

-  `Mutex` A read-write lock that guarantees `WorkerStatus` the concurrent security access of the.
-  `WorkerStatus` A map that stores the relationship between a port number and the corresponding worker object.

Here is an explanation of the important methods in the code:

-  `NewWorkers` Is a factory function that creates and returns a new `Workers` object. It initializes to `WorkerStatus` an empty map.
-  `Add` Method is used to add a worker to the worker collection. It takes the port number `port` and the worker object `w` and adds it to `WorkerStatus` the, keyed by the port number. At the same time, it increments `serverInstance.Counter` the count of and sets the state of the corresponding port to `true`.
-  `Remove` Method is used to remove a worker of a specified port number from the worker collection. It closes the listener for the worker object and removes the worker from `WorkerStatus` it. At the same time, it sets the status of the corresponding port to `false`.
-  `Get` Method is used to get the corresponding worker object based on the port number. It is accessed `WorkerStatus` as read-only and returns the worker object for the specified port number.
-  `NewWorker` Is a factory function that creates and returns a new worker object. It takes the listener `l`, client connection `c`, and port number `port` and assigns them as fields to the newly created worker object.

These structures and methods are used to manage and manipulate worker objects. The worker object represents the corresponding port number of the server and maintains the connection with the client and the listener of the server. `Workers` The struct manages a set of worker objects and provides operations for adding, removing, and retrieving worker objects. Through these structures and methods, we can easily manage worker objects, achieve concurrent processing and control connections.

### 5.6. server object

This is the core of the whole code.


```go
var serverInstance *Server

// Server 服务端程序的实例
type Server struct {
	// Mutex 保证并发安全的锁
	Mutex sync.RWMutex
	// Counter 目前服务器累计接收到了多少次连接
	Counter int64
	// 最大连接数量
	MaxTCPConnSize int32
	// 最大连接数量
	MaxConnSize int32
	// ExposePort 服务端暴露端口
	ExposePort []int
	// ProcessingMap
	ProcessingMap map[string]*net.TCPConn
	// WorkerBuffer 整体工作队列的大小
	WorkerBuffer chan *net.TCPConn
	// 实际处理工作的数据结构
	ProcessWorker *Workers
	// 端口使用情况
	PortStatus map[int32]bool
}

func initServer() {
	serverInstance = &Server{
		Mutex:          sync.RWMutex{},
		Counter:        0,
		MaxTCPConnSize: objectConfig.MaxTCPConnNum,
		MaxConnSize:    objectConfig.MaxConnNum,
		ExposePort:     objectConfig.ExposePort,
		ProcessingMap:  make(map[string]*net.TCPConn),
		WorkerBuffer:   make(chan *net.TCPConn, objectConfig.MaxConnNum),
		ProcessWorker:  NewWorkers(),
		PortStatus:     make(map[int32]bool),
	}

	// 初始化端口状态
	for i := 0; i < int(serverInstance.MaxTCPConnSize); i++ {
		serverInstance.PortStatus[int32(serverInstance.ExposePort[i])] = false
	}
}

func (s *Server) PortIsFull() bool {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	for _, v := range s.PortStatus {
		if v == false {
			return false
		}
	}
	return true
}

func (s *Server) GetPort() int32 {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	for k, v := range s.PortStatus {
		if v == true {
			continue
		} else {
			return k
		}
	}
	return -1
}
```


The above code defines a `Server` struct and associated methods.

 `Server` The struct represents an instance of a server-side program and contains the following fields:

-  `Mutex` Read-write locks, which guarantee concurrent security access to the server state.
-  `Counter` : Cumulative number of connections received.
-  `MaxTCPConnSize` : Maximum number of TCP connections.
-  `MaxConnSize` : Maximum number of connections.
-  `ExposePort` : The port number exposed by the server.
-  `ProcessingMap` : Stores a map of the connection being processed, with the connection ID in string format as the key and the corresponding TCP connection as the value.
-  `WorkerBuffer` : a buffer for the overall work queue to store pending TCP connections.
-  `ProcessWorker` : The set of workers that process connections.
-  `PortStatus` : The mapping of the port usage, with the port number as the key, indicating whether the port is occupied.

Here is an explanation of the important methods in the code:

-  `initServer` Functions are used to initialize `serverInstance`, create, and set `Server` the fields of an object. Where `MaxTCPConnNum` the values for the, `MaxConnNum`, `ExposePort`, and so on fields come from the global configuration object `objectConfig`. At the same time, it also initializes the port state mapping `PortStatus`, initializing the state of each port to `false`.
-  `PortIsFull` Method is used to check whether all ports are full. It is accessed `PortStatus` as read-only and returns `true` if all ports have a status `true` of, indicating that all ports are occupied; otherwise, it returns `false` indicating that at least one port is available.
-  `GetPort` Method is used to get an available port number. It iterates through the port numbers in the map with read-only access `PortStatus`, and returns the port number if it finds a port with a status of `false`, or if all ports are in use, it returns `-1` that no ports are available.

These methods provide the ability to initialize the server instance, manage the port state, and obtain available ports. Through these methods, the number of connections and the status of ports can be controlled and monitored.

### 5.7. Tablishing control connection channel


```go
// createControllerChannel 创建一个控制信息的通道，用于接收内网客户端的连接请求
// 当内网客户端向服务端的控制接口发送请求建立连接时，控制端会直接向全局的工作队列中添加这个连接信息
// 可以在此进行用户权限的界别与控制
// Create a control information channel to receive connection requests from intranet clients;
// When an intranet client sends a connection request to the control interface of the server,
// the control side will directly add this connection information to the global work queue.
// You can implement user-level permissions and control at this point.
func createControllerChannel() {
	controllerListener, err := network.CreateTCPListener(objectConfig.ControllerAddr)
	if err != nil {
		fmt.Println("[createTCPListener]", err)
		panic(err)
	}
	fmt.Println("[服务器控制端开始监听]" + objectConfig.ControllerAddr)
	if objectConfig.StartAuth == "true" {
		// 获取用户发送来的数据
		fmt.Println("[Start Auth Successfully!]", "服务器开启认证请求")
	}
	for {
		tcpConn, err := controllerListener.AcceptTCP()
		if err != nil {
			fmt.Println("[AcceptTCP]", err)
			continue
		}
		if objectConfig.StartAuth == "true" {
			err := authUser(tcpConn)
			if err != nil {
				fmt.Println(err)
				usi := instance.NewSendAndReceiveInstance(tcpConn)
				_, err = usi.SendDataToClient(network.AUTH_FAIL, []byte{})
				fmt.Println("[AUTH_FAIL]", "发送认证失败消息")
				_ = tcpConn.Close()
				continue
			}
		}
		// 给客户端发送该消息
		fmt.Println("[控制层接收到新的连接]", tcpConn.RemoteAddr())
		// 将新地连接推入工作队列中去
		serverInstance.WorkerBuffer <- tcpConn
		log.Infoln("[%s] %s\n", tcpConn.RemoteAddr().String(), "已推入工作队列中。")
	}
}
```

The main purpose of this function is to create a listener for the control port and handle the connection request.

The function first calls the `network.CreateTCPListener` method to create a TCP listener `controllerListener` to listen to the specified address `objectConfig.ControllerAddr`. If an error occurs while creating the listener, an error message is printed and an exception is thrown.

Next, the function prints the information that the control port started listening.

If `objectConfig.StartAuth` the value of is `"true"`, the server opens the authentication request. In this case, a message that the authentication was successful is printed.

Next, the function enters an infinite loop, continuously accepting connection requests from the control port. Accepts a new TCP connection request by calling `controllerListener.AcceptTCP()` a method. If an error occurs while accepting the connection, an error message is printed and the next cycle continues.

If `objectConfig.StartAuth` the value of is `"true"`, the server opens the authentication request. In this case, a method is called `authUser` to authenticate the user. If authentication fails, an error message is printed, an authentication failure message is sent to the client, and the connection is closed.

If the authentication passes or the server does not turn on the authentication request, a message that a new connection has been received is printed.

The new connection `tcpConn` is then pushed into the server instance's work `serverInstance.WorkerBuffer` queue for subsequent processing. At the same time, a log is kept, printing the remote address of the connection and information pushed to the work queue.

In general, this function creates a listener for the control port and handles the connection request. If the authentication request is turned on, the connection is authenticated. Legitimate connections are then pushed into the server instance's work queue for subsequent processing.

### 5.8. Authenticate the client


```go
func authUser(conn *net.TCPConn) error {
	nsi := instance.NewSendAndReceiveInstance(conn)
	msg, err := nsi.ReadHeadDataFromClient()
	if err != nil {
		return err
	}
	msg, err = nsi.ReadRealDataFromClient(msg)
	if err != nil {
		return err
	}
	if msg.GetMsgID() == network.USER_REQUEST_AUTH {
		// 获取其真实数据
		ui := new(network.UserInfo)
		err := ui.FromBytes(msg.GetMsgData())
		if err != nil {
			return err
		}
		dbInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", objectConfig.DB.Username, objectConfig.DB.Password, objectConfig.DB.Host, objectConfig.DB.DBName)
		cui := network.NewControllerUserInfo([]byte(network.KEY), "mysql", dbInfo)
		err = cui.CheckUser(ui)
		if err != nil {
			return err
		}

	}
	return nil
}

```

This function is used to authenticate the user. It receives a `net.TCPConn` parameter `conn` of type that represents the TCP connection to the client.

The function first uses the `instance.NewSendAndReceiveInstance` method to create a send receive instance `nsi` for reading and sending data to and from the client.

Next, the function reads the message header sent by the client by calling `nsi.ReadHeadDataFromClient` the method. If an error occurs during the read, the error is returned.

The function then reads the real data sent by the client by calling `nsi.ReadRealDataFromClient` the method. If an error occurs during the read, the error is returned.

Next, the function determines whether the message ID of the received message is `network.USER_REQUEST_AUTH`, that is, the message for which the user requests authentication.

If the message is a user request for authentication, the function will parse out the user information `ui`. Then create an `network.ControllerUserInfo` instance `cui` based on the database information and key in the configuration file.

Next, the function calls `cui.CheckUser` a method to authenticate the user. If authentication fails, an appropriate error is returned.

Finally, if the function completes without error, it returns `nil`, indicating that the authentication was successful.

In general, this function authenticates the client. It reads the message header and real data from the client, determines whether it is the message that the user requests authentication, and then uses the database information to authenticate the user and returns the corresponding results.

### 5.9. Heartbeat bag


```go
// keepAlive 心跳包检测,函数负责向客户端发送保活消息以确保连接处于活动状态(每三秒发送一次)。如果在此过程中发生错误，它会检查错误是否表示客户端已关闭连接。
// 如果是，则会记录相应的日志，并从工作队列中移除相应的端口。然后函数返回。
// The keepAlive function is responsible for sending a keep-alive message to the client to ensure the connection is active.
// If an error occurs during the process, it checks if the error indicates that the client has closed the connection.
// If so, it logs the appropriate message and removes the corresponding port from the worker queue.
// The function then returns.
func keepAlive(conn *net.TCPConn, port int32) {
	for {
		nsi := instance.NewSendAndReceiveInstance(conn)
		_, err := nsi.SendDataToClient(network.KEEP_ALIVE, []byte("ping"))
		log.Infoln("SendData [ping] Successfully.")
		if err != nil {
			log.Errorln("[检测到客户端关闭]", err)
			serverInstance.ProcessWorker.Remove(port)
		}
		time.Sleep(time.Minute*5)
	}
}

```

 `keepAlive` Function is responsible for sending keep-alive messages to the client to ensure that the connection is alive. This function sends a keep-alive message every five seconds. If an error occurs during this process and the error indicates that the client has closed the connection, the function logs the appropriate information and removes the corresponding port from the work queue. The function then returns.

Specifically, the function first creates an `SendAndReceiveInstance` object `nsi` that is used to send and receive data with the client. It is then used `nsi` to send a packet containing a keep-alive message to the client with the message "ping". If an error occurs during the sending process, that is `err`, no `nil`, the client has closed the connection. The function logs the corresponding error and calls `serverInstance.ProcessWorker.Remove(port)` the method to remove the corresponding port from the work queue.

Function `time.Sleep(time.Minute)` to pause the program for one minute before the next keep-alive message is sent. This loop keeps the connection to the client alive.

### 5.10. Monitor


```go
// ListenTaskQueue 该函数的作用是监听工作队列传来的消息。
// 它通过不断检查工作队列是否有可用的连接，并将连接分配给处理函数 acceptUserRequest。
// 当工作队列未满时，会从工作队列中取出一个连接，并启动一个协程来处理该连接的用户请求。
// 函数会以很小的时间间隔进行轮询，并持续监听工作队列的新消息。
// The function listens for messages from the work queue.
// It does this by constantly checking the work queue for an available connection and assigning the connection to the handler function acceptUserRequest.
// When the work queue is not full, a connection is taken from the work queue and a coroutine is started to process user requests for that connection.
// The function polls at small intervals and continuously listens for new messages from the work queue.
func ListenTaskQueue() {
	log.Infoln("[ListenTaskQueue]", "监听工作队列传来的消息")
restLabel:
	if !serverInstance.PortIsFull() {
		conn := <-serverInstance.WorkerBuffer
		go acceptUserRequest(conn)
	}
	time.Sleep(time.Millisecond * 100)
	goto restLabel
}
```

### 5.11. A us request is receive


```go

// acceptUserRequest 接收请用户的求,该函数会首先从全局工作池中获取一个空闲的端口，然后在这个端口上监听用户的请求
// 并向客户端发送对应的信息，然后在这个端口监听用户的请求，每监听到一个请求，就向内网客户端发送一个建立通道的信号
// The function acceptUserRequest is responsible for accepting user requests.
// It first retrieves an available port from the global worker pool.
// It then listens for incoming requests on this port and sends corresponding information to the client.
// Each time a request is received, it sends a signal to establish a channel with the internal network client.
func acceptUserRequest(conn *net.TCPConn) {
	port := serverInstance.GetPort()
	userVisitAddr := "0.0.0.0:" + strconv.Itoa(int(port))
	userVisitListener, err := network.CreateTCPListener(userVisitAddr)
	if err != nil {
		log.Errorln("[CreateTCPListener]", err)
		return
	}
	defer userVisitListener.Close()
	workerInstance := NewWorker(userVisitListener, conn, port)
	serverInstance.ProcessWorker.Add(port, workerInstance)
	c := network.NewClientConnInstance(serverInstance.Counter, port)
	ready, _ := c.ToBytes()
	nsi := instance.NewSendAndReceiveInstance(conn)
	go keepAlive(conn, port)
	_, err = nsi.SendDataToClient(network.USER_AUTHENTICATION_SUCCESSFULLY, []byte{})
	_, err = nsi.SendDataToClient(network.USER_INFORMATION, ready)
	if err != nil {
		log.Infoln("[Send Client info]", err)
		return
	}
	log.Infoln("[addr]", userVisitListener.Addr().String())
	for {
		tcpConn, err := userVisitListener.AcceptTCP()
		if opErr, ok := err.(*net.OpError); ok {
			if strings.Contains(opErr.Error(), "use of closed network connection") {
				// 远程主机关闭连接，退出连接处理循环
				log.Infoln("远程客户端连接关闭")
				return
			}
		}
		if err != nil {
			log.Errorln("[userVisitListener.AcceptTCP]", err)
			continue
		}
		userConnPoolInstance.AddConnInfo(tcpConn)
		nsi := instance.NewSendAndReceiveInstance(conn)
		count, err := nsi.SendDataToClient(network.NEW_CONNECTION, []byte(network.NewConnection))
		if err != nil {
			log.Errorln("[SendData fail]", err)
			continue
		}
		log.Infoln("[SendData successfully]", count, " bytes")
	}
}

```

 `acceptUserRequest` The function is responsible for receiving the user's request. It first obtains a free port from the global worker pool, and then listens for user requests on that port. For each received request, it sends the corresponding information to the client and a signal to the internal network client to establish the channel.

The specific analysis is as follows:

- First, the function calls `serverInstance.GetPort()` to get a free port from the global worker pool.
- The function then creates a TCP listener `userVisitListener` based on the port obtained, with the listening address "0.0.0.0: port number".
- Next, the function creates an `Worker` instance `workerInstance` that includes `userVisitListener`, `conn` and the port number obtained, and adds it to `serverInstance.ProcessWorker`.
- Create an `ClientConnInstance` instance `c` to communicate with the client and convert it to a stream of bytes.
- An object `nsi` is used `SendAndReceiveInstance` to send an authentication success message and user information to the client.
- The function then starts a coroutine-calling `keepAlive` function that periodically sends keep-alive messages to keep the connection with the client active.
- Loop listens to the user's request and receives the connection request through the call `userVisitListener.AcceptTCP()`. If an error occurs, it is handled accordingly, and if the remote host closes the connection, the connection processing loop is exited.
- If the connection request is successfully received, the function adds the connection information to `userConnPoolInstance` the and uses `nsi` to signal to the client that the new connection is established.

In summary, `acceptUserRequest` the function is responsible for receiving the user's request, communicating with the client on the appropriate port, maintaining the active state of the connection, and handling the establishment and closure of the connection.

### 5.12. Receiving a client request


```go
// acceptClientRequest 该函数用于接收客户端的请求连接。它首先创建一个监听指定隧道地址的TCP监听器。如果创建监听器时发生错误，则记录错误并返回。函数执行完毕后会关闭监听器。
// 然后，函数进入一个无限循环，接受来自客户端的TCP连接请求。每当接收到一个连接请求时，会创建一个新的协程来处理该连接，即创建隧道。
// This function is responsible for accepting client connection requests.
// It first creates a TCP listener for the specified tunnel address.
// If an error occurs during the creation of the listener, the error is logged and the function returns.
// The listener is closed when the function completes.
// Then, the function enters an infinite loop to accept TCP connection requests from clients. For each incoming connection request, a new goroutine is spawned to handle the connection by creating a tunnel.
func acceptClientRequest() {
	tunnelListener, err := network.CreateTCPListener(objectConfig.TunnelAddr)
	if err != nil {
		log.Errorln("[CreateTunnelListener]" + objectConfig.TunnelAddr + err.Error())
		return
	}
	defer tunnelListener.Close()
	for {
		tcpConn, err := tunnelListener.AcceptTCP()
		if err != nil {
			log.Errorln("[TunnelAccept]", err)
			continue
		}
		// 创建隧道
		go createTunnel(tcpConn)
	}
}
```

 `acceptClientRequest` Function to receive a connection request from a client. The specific analysis is as follows:

- First, the function calls `network.CreateTCPListener(objectConfig.TunnelAddr)` to create a TCP listener `tunnelListener` for the specified tunnel address.
- If an error occurs while creating the Listener, the function logs the error and returns.
- After the function finishes executing, the listener is closed `tunnelListener` and the resources are released.
- Next, the function enters an infinite loop, receiving the client's TCP connection request by calling `tunnelListener.AcceptTCP()`.
- If the connection request is successfully received, the function creates a new coroutine and calls the `createTunnel` function to handle the connection, creating a tunnel.
- Each connection request is processed in a new coroutine, allowing multiple client connection requests to be processed simultaneously.

In summary, `acceptClientRequest` the function is responsible for creating and listening to the TCP listener for the specified tunnel address, then receiving the client's connection requests in an infinite loop, and creating a new coroutine on each request to handle the connection to create the tunnel. This enables simultaneous processing of connection requests from multiple clients and maintains the scalability of the server.

### 5.13. Tunnel


```go

// createTunnel 该函数用于创建一个隧道。
// 函数首先获取用户连接池的读锁，以保证在创建隧道期间不会有其他线程修改连接池。
// 然后它遍历用户连接池中的每个连接，找到一个可用的连接，将该连接与传入的隧道进行数据交换，然后从连接池中删除该连接。
// 如果没有找到可用连接，函数会关闭传入的隧道。最后，释放用户连接池的读锁。
// This function is used to create a tunnel.
// It first acquires the read lock of the user connection pool to ensure that no other threads modify the connection pool during the creation of the tunnel.
// Then it iterates through each connection in the user connection pool to find an available connection.
// It swaps data between the found connection and the provided tunnel, and then removes the connection from the connection pool.
// If no available connection is found, the function closes the provided tunnel. Finally, it releases the read lock of the user connection pool.
func createTunnel(tunnel *net.TCPConn) {
	userConnPoolInstance.Mutex.RLock()
	defer userConnPoolInstance.Mutex.RUnlock()

	for key, connMatch := range userConnPoolInstance.UserConnectionMap {
		if connMatch.conn != nil {
			go network.SwapConnDataEachOther(connMatch.conn, tunnel)
			delete(userConnPoolInstance.UserConnectionMap, key)
			return
		}
	}

	_ = tunnel.Close()
}
```

 `createTunnel` The tunnel () function creates a tunnel. The specific analysis is as follows:

- The function first acquires a read lock `userConnPoolInstance.Mutex.RLock()` on the user's connection pool to ensure that no other threads modify the connection pool during tunnel creation.
- The function then iterates through each connection `userConnPoolInstance.UserConnectionMap` in the user connection pool looking for available connections.
- If an available connection `connMatch.conn` is found, the function calls `network.SwapConnDataEachOther(connMatch.conn, tunnel)` Data Exchange to exchange the connection with the incoming tunnel.
- After the data exchange is complete, the function removes the connection from the connection pool `delete(userConnPoolInstance.UserConnectionMap, key)`.
- If no available connection is found in the connection pool, the function closes the incoming tunnel `tunnel`, which is a call to `tunnel.Close()` close the tunnel.
- Finally, the function releases the read lock `userConnPoolInstance.Mutex.RUnlock()` on the user's connection pool.

In summary, `createTunnel` the function is responsible for finding available connections in the user connection pool and exchanging data with the incoming tunnel to create a tunnel. If an available connection is found, data exchange takes place and the connection is removed from the connection pool; if no connection is available, the tunnel is closed. This function ensures data security of the connection pool during tunnel creation and handles cases where the connection pool is empty.

### 5.14. Clear


```go
// cleanExpireConnPool 该函数用于清理连接池中的过期连接。
// 函数会进入一个无限循环，在每次循环中，它获取连接池的互斥锁，遍历连接池中的每个连接。
// 如果某个连接的访问时间距离当前时间已经超过10秒，那么该连接会被关闭，并从连接池中删除。
// 完成遍历后，释放连接池的互斥锁。函数会每隔5秒执行一次清理操作。
// This function is responsible for cleaning up expired connections in the connection pool.
// It enters an infinite loop, and in each iteration, it acquires the mutex lock of the connection pool.
// It iterates through each connection in the pool.
// If the time elapsed since the last visit of a connection exceeds 10 seconds, the connection is closed and removed from the connection pool.
// After the iteration is complete, the mutex lock of the connection pool is released.
// The function performs the cleanup operation every 5 seconds.
func cleanExpireConnPool() {
	for {
		userConnPoolInstance.Mutex.Lock()
		for key, connMatch := range userConnPoolInstance.UserConnectionMap {
			if time.Now().Sub(connMatch.visit) > time.Second*8 {
				_ = connMatch.conn.Close()
				delete(userConnPoolInstance.UserConnectionMap, key)
			}
		}
		log.Infoln("[cleanExpireConnPool successfully]")
		userConnPoolInstance.Mutex.Unlock()
		time.Sleep(5 * time.Second)
	}
}
```

 `cleanExpireConnPool` Function is used to clean up expired connections in the connection pool. The specific analysis is as follows:

- The function enters an infinite loop.
- In each loop, the function first acquires the mutex `userConnPoolInstance.Mutex.Lock()` of the connection pool to ensure that no other threads modify the connection pool during cleanup.
- The function then walks through each connection in the connection pool, `userConnPoolInstance.UserConnectionMap` iterating through the.
- For each connection, the function checks whether the difference between the current time and the last access time is more than 10 seconds, that is `time.Now().Sub(connMatch.visit) > time.Second*8`.
- If the connection expires, the function closes the connection `connMatch.conn` and removes it from the connection pool `delete(userConnPoolInstance.UserConnectionMap, key)`.
- After the traversal is complete, the function releases the connection pool's mutex `userConnPoolInstance.Mutex.Unlock()`.
- Function performs a cleanup operation every 5 seconds. `time.Sleep(5 * time.Second)`.

In summary, `cleanExpireConnPool` the function periodically cleans up expired connections in the connection pool, closing and deleting connections that have not been accessed for a certain period of time to maintain the validity and availability of the connection pool. The function uses mutexes to ensure data security during cleanup and periodically executes cleanup operations to preserve the state of the connection pool.

### 5.15. main


```go
var log = logrus.New()

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		// 只打印帮助信息，不执行命令
		Execute()
	} else {
		Execute()
		art()
		exchange()
		printServerRelationInformation()
		go createControllerChannel()
		go ListenTaskQueue()
		go acceptClientRequest()
		go cleanExpireConnPool()
		select {}
	}

}

func init() {
	objectConfig = new(objectConfigData)
	initConfig()
	initCobra()
	initLog()
	initServer()
	initUserConnPool()
}

```

