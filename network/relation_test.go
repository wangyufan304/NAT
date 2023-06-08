package network

import (
	"fmt"
	"testing"
	"time"
)

func TestNewUserInfoInstance(t *testing.T) {
	ui := NewUserInfoInstance("wyf", "123456")
	ui.ExpireTime = time.Now()
	d, err := ui.ToBytes()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d)
	un := new(UserInfo)
	un.FromBytes(d)
	fmt.Println(un)
}
