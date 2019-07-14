package models

import (
	"scrapyd-admin/core"
	"scrapyd-admin/resource/sql"
	"bytes"
	"errors"
	"github.com/ltachi1/logrus"
	"os"
)

func init() {
	error := InitTables()
	if error != nil {
		core.WriteLog(core.LogTypeDb, logrus.PanicLevel, nil, error)
		os.Exit(1)
	}
	InitTask()
}

func InitTables() error {
	//判断表是否存在
	exist, err := core.Db.IsTableExist("admin")
	if err != nil {
		errors.New("数据库表创建失败")
	}
	if exist {
		return nil
	}
	sql, err := sql.Asset("sa.sql")
	if err != nil {
		return err
	}
	_, err = core.Db.Import(bytes.NewReader(sql))
	if err != nil {
		return err

	}
	return nil
}

func InitTask(){
	//定时检测服务器状态
	core.Cron.AddFunc("*/30 * * * *", func() {
		new(Server).DetectionStatus()
	}, "DetectionServerStatus")
	//定时检测任务状态
	core.Cron.AddFunc("*/5 * * * *", func() {
		new(Task).DetectionStatus()
	}, "DetectionTaskStatus")

	//初始化已有的计划任务
	new(SchedulesTask).InitSchedulesToCron()
}