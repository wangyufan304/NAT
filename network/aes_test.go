package network

import (
	"fmt"
	"testing"
)

func TestEncryptData(t *testing.T) {
	data, err := EncryptData([]byte("wyfwyfwyfwyfwyfw"), []byte("wyf"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
	bytes, err := DecryptData([]byte("wyfwyfwyfwyfwyfw"), data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))
}
