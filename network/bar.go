package network

import (
	"database/sql"
)

type ControllerBar struct {
	DBController *sql.DB
}

// NewControllerBar 新建一个控制实例
func NewControllerBar(dbType, dbInfo string) *ControllerBar {
	db, err := sql.Open(dbType, dbInfo)
	if err != nil {
		panic(err)
	}
	return &ControllerBar{
		DBController: db,
	}
}

func (bar *ControllerBar) AddBar(name string, count int64) error {
	// 查询语句
	query := "SELECT COUNT(*) FROM bars WHERE username = ?"
	var existingCount int
	err := bar.DBController.QueryRow(query, name).Scan(&existingCount)
	if err != nil {
		return err
	}

	if existingCount == 0 {
		// 不存在，则执行插入操作
		insert := "INSERT INTO bars (username, bar) VALUES (?, ?)"
		_, err = bar.DBController.Exec(insert, name, count)
	} else {
		// 存在，则执行更新操作
		update := "UPDATE bars SET bar = bar + ? WHERE username = ?"
		_, err = bar.DBController.Exec(update, count, name)
	}

	if err != nil {
		return err
	}
	return nil
}
