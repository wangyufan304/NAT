package network

import (
	"fmt"
	"testing"
)

func TestControllerInfo_ToByteBuffer(t *testing.T) {
	ci := &ControllerInfo{
		ID:   1,
		Port: 80001,
	}
	buf, err := ci.ToByteBuffer()
	if err != nil {
		fmt.Println("[To]", err)
		return
	}
	fmt.Println("[Result]", buf)
}

func TestKeepAlive_ToByteBuffer(t *testing.T) {
	ka := &KeepAlive{
		ID:  1,
		Msg: "ping",
	}
	buf, err := ka.ToByteBuffer()
	if err != nil {
		fmt.Println("[To]", err)
		return
	}
	fmt.Println("[Result]", buf)
}
