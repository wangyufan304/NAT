package ninterfance

// ISendAndReceiveData 统一发送和接收数据的接口
type ISendAndReceiveData interface {
	SendDataToClient(uint32, []byte) (int, error)
	ReadHeadDataFromClient() (IMessage, error)
	ReadRealDataFromClient(IMessage) (IMessage, error)
}
