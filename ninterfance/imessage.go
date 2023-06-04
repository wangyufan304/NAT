package ninterfance

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
