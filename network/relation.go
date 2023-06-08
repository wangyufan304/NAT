package network

import (
	"bytes"
	"encoding/binary"
	"time"
)

// UserInfo 用户信息模块
type UserInfo struct {
	// UserName
	UserName string
	// Password
	Password string
	// ExpireTime
	ExpireTime time.Time
}

// NewUserInfoInstance 新建一个实体
func NewUserInfoInstance(username, password string) *UserInfo {
	return &UserInfo{
		UserName: username,
		Password: password,
	}
}

func (info *UserInfo) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 将用户名长度编码到字节流中
	userNameLen := len(info.UserName)
	if err := binary.Write(buf, binary.BigEndian, int32(userNameLen)); err != nil {
		return nil, err
	}

	// 将用户名内容编码到字节流中
	if err := binary.Write(buf, binary.BigEndian, []byte(info.UserName)); err != nil {
		return nil, err
	}

	// 将密码长度编码到字节流中
	passwordLen := len(info.Password)
	if err := binary.Write(buf, binary.BigEndian, int32(passwordLen)); err != nil {
		return nil, err
	}

	// 将密码内容编码到字节流中
	if err := binary.Write(buf, binary.BigEndian, []byte(info.Password)); err != nil {
		return nil, err
	}

	// 将过期时间编码到字节流中
	if err := binary.Write(buf, binary.BigEndian, info.ExpireTime.Unix()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (info *UserInfo) FromBytes(data []byte) error {
	buf := bytes.NewReader(data)

	// 从字节流中读取用户名长度
	var userNameLen int32
	if err := binary.Read(buf, binary.BigEndian, &userNameLen); err != nil {
		return err
	}

	// 从字节流中读取用户名内容，并根据长度进行截取
	userNameBytes := make([]byte, userNameLen)
	if err := binary.Read(buf, binary.BigEndian, userNameBytes); err != nil {
		return err
	}
	info.UserName = string(userNameBytes)

	// 从字节流中读取密码长度
	var passwordLen int32
	if err := binary.Read(buf, binary.BigEndian, &passwordLen); err != nil {
		return err
	}

	// 从字节流中读取密码内容，并根据长度进行截取
	passwordBytes := make([]byte, passwordLen)
	if err := binary.Read(buf, binary.BigEndian, passwordBytes); err != nil {
		return err
	}
	info.Password = string(passwordBytes)

	// 从字节流中读取过期时间
	var expireTimeUnix int64
	if err := binary.Read(buf, binary.BigEndian, &expireTimeUnix); err != nil {
		return err
	}
	info.ExpireTime = time.Unix(expireTimeUnix, 0)

	return nil
}
