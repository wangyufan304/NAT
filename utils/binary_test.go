package utils

import (
	"fmt"
	"github.com/byteYuFan/NAT/network"
	"testing"
)

func TestObjectToBufferStream(t *testing.T) {
	ci := &network.ControllerInfo{
		ID:   1,
		Port: 8080,
	}
	buf, err := ObjectToBufferStream(ci)
	if err != nil {
		fmt.Println("[To]", err)
		return
	}
	fmt.Println("[Result]", buf)
}
