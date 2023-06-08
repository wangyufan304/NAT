package network

import (
	"fmt"
	"testing"
	"time"
)

func TestNewControllerUserInfo(t *testing.T) {
	cui := NewControllerUserInfo([]byte("WYFFYWYYTT123456"), "mysql", "root:123456@tcp(pogf.com.cn:3309)/NAT")
	user := UserInfo{
		UserName:   "admin",
		Password:   "123456",
		ExpireTime: time.Now().Add(time.Hour * 1314521),
	}
	id, err := cui.Add(&user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(id)
}

func TestControllerUserInfo_CheckUser(t *testing.T) {
	cui := NewControllerUserInfo([]byte("WYFFYWYYTT123456"), "mysql", "root:123456@tcp(pogf.com.cn:3309)/NAT")
	user := UserInfo{
		UserName:   "admin",
		Password:   "123456",
		ExpireTime: time.Now().Add(time.Hour * 24),
	}
	err := cui.CheckUser(&user)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestControllerUserInfo_SetExpireTime(t *testing.T) {
	cui := NewControllerUserInfo([]byte("wyf1234567891011"), "mysql", "root:wyf20040305...@tcp(docker:3306)/NAT")
	user := UserInfo{
		UserName:   "admin",
		Password:   "123456",
		ExpireTime: time.Now().Add(time.Hour * 24),
	}
	err := cui.SetExpireTime(&user, time.Now().Add(time.Hour*48))
	if err != nil {
		fmt.Println(err)
		return
	}
}
