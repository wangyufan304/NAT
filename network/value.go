package network

const (
	KEY = "WYFFYWYYTT123456"
)

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
	PASSWORD_INCORRET  = 30004
	AUTH_FAIL          = 30005
)

const (
	USER_REQUEST_AUTH = 60001
)

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
