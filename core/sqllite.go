// 数据库操作
package core

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/ltachi1/logrus"
	"os"
	"fmt"
)

var Db *xorm.Engine

func InitDb(dbPath string) {
	if dbPath == "" {
		dbPath = RootPath + "/db"
	}
	dbPath = SupplementDir(dbPath)
	var err error
	err = os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		fmt.Println("数据库存放路径输入不正确或者没有写入权限")
		os.Exit(1)
	}
	//初始化数据库链接
	Db, err = xorm.NewEngine("sqlite3", dbPath + "sa.db")
	if err != nil {
		WriteLog(LogTypeAdmin, logrus.FatalLevel, nil, "数据库配置初始化失败")
		os.Exit(1)
	}
	//设置最大连接数和空闲连接数
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(5)
	if err != nil {
		WriteLog(LogTypeAdmin, logrus.FatalLevel, nil, "数据库配置初始化失败")
		os.Exit(1)
	}
	//设置表前缀
	Db.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, ""))
	//Db.ShowSQL(true)
}
