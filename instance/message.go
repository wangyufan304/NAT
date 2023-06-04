package instance

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
