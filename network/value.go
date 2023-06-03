package network

const (
	CONTROLLER_INFO = 20001
	KEEP_ALIVE      = 20002
)

// TCPProtocol 定义传送协议
var TCPProtocol map[int]interface{}

func init() {
	TCPProtocol = make(map[int]interface{})
	TCPProtocol[1] = "NEW_CONNECTION"
	TCPProtocol[CONTROLLER_INFO] = ControllerInfo{}
	TCPProtocol[KEEP_ALIVE] = KeepAlive{}
}
