package network

const (
	NEW_CONNECTION   = 20001
	USER_INFORMATION = 20002
	KEEP_ALIVE       = 20003
)

var ProtocolMap map[uint32]interface{}

func init() {
	ProtocolMap = make(map[uint32]interface{})
	// 添加创建协议
	ProtocolMap[NEW_CONNECTION] = "NEW_CONNECTION"
	ProtocolMap[USER_INFORMATION] = ClientConnInfo{}
	ProtocolMap[KEEP_ALIVE] = "ping"
}
