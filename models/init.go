package models

import (
	"scrapyd-admin/core"
	"scrapyd-admin/resource/sql"
	"bytes"
	"errors"
	"scrapyd-admin/config"
	"github.com/ltachi1/logrus"
)

func init() {
	error := InitTables()
	if error != nil {
		core.WriteLog(config.LogTypeDb, logrus.PanicLevel, nil, error)
	}
	InitTask()
}

func InitTables() error {
	//判断表是否存在
	exist, err := core.DBPool.Master().IsTableExist("admin")
	if err != nil {
		errors.New("数据库表创建失败")
	}
	if exist {
		return nil
	}
	sql, err := sql.Asset("scrapyd_admin.sql")
	if err != nil {
		return err
	}
	_, err = core.DBPool.Master().Import(bytes.NewReader(sql))
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