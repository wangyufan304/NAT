package ninterfance

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
