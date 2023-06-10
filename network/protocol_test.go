package network

import (
	"fmt"
	"testing"
)

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
