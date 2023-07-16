package network

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type ControllerUserInfo struct {
	DBController *sql.DB
	KEY          []byte
}

// NewControllerUserInfo 新建一个控制实例
func NewControllerUserInfo(key []byte, dbType, dbInfo string) *ControllerUserInfo {
	db, err := sql.Open(dbType, dbInfo)
	if err != nil {
		panic(err)
	}
	return &ControllerUserInfo{
		DBController: db,
		KEY:          key,
	}
}

// Add 插入用户信息
func (cui *ControllerUserInfo) Add(user *UserInfo) (int64, error) {
	defer cui.DBController.Close()
	// Check if username already exists
	existsQuery := "SELECT COUNT(*) FROM userInfo_userinfo WHERE username = ?"
	var count int64
	err := cui.DBController.QueryRow(existsQuery, user.UserName).Scan(&count)
	if err != nil {
		return -1, err
	}
	if count > 0 {
		// Username already exists
		fmt.Println("用户存在")
		return -1, errors.New(ProtocolMap[USER_ALREADY_EXIST].(string))
	}
	//encryptPassword, _ := EncryptData(cui.KEY, []byte(user.Password))
	encryptPassword := user.Password
	insertQuery := "INSERT INTO userInfo_userinfo  (username, password, time) VALUES (?, ?, ?)"
	result, err := cui.DBController.Exec(insertQuery, user.UserName, encryptPassword, user.ExpireTime)
	if err != nil {
		return -1, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return lastInsertID, nil
}

func (cui *ControllerUserInfo) SetExpireTime(user *UserInfo, time time.Time) error {

	modifyQuery := "UPDATE userInfo_userinfo SET time = ? WHERE username = ?"
	_, err := cui.DBController.Exec(modifyQuery, user.ExpireTime, user.UserName)
	if err != nil {
		return err
	}
	return nil
}
func (cui *ControllerUserInfo) CheckUser(user *UserInfo) error {
	defer cui.DBController.Close()
	fmt.Println("进入了")
	query := "SELECT password, time FROM userInfo_userinfo WHERE username = ?"
	result, err := cui.DBController.Query(query, user.UserName)
	defer result.Close()
	if err != nil {
		return err
	}

	var password string
	var expireTime []byte

	for result.Next() {
		if err := result.Scan(&password, &expireTime); err != nil {
			return err
		}
	}

	//decryptPassword, _ := DecryptData(cui.KEY, password)
	decryptPassword := password
	if string(decryptPassword) != user.Password {
		return errors.New(ProtocolMap[PASSWORD_INCORRET].(string))
	}

	expireTimeValue, err := time.Parse("2006-01-02 15:04:05", string(expireTime))
	if err != nil {
		return err
	}

	if expireTimeValue.Before(time.Now()) {
		return errors.New(ProtocolMap[USER_EXPIRED].(string))
	}

	return nil
}
