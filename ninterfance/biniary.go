package ninterfance

type ObjectTOByteBuffer interface {
	ToByteBuffer() ([]byte, error)
}
type BufferToObjectInterface interface {
	BufferToObject([]byte) error
}
