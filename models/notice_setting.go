package models

import (
	"scrapyd-admin/core"
	"fmt"
	"github.com/ltachi1/logrus"
)

type NoticeSetting struct {
	Base  core.BaseModel `json:"-" xorm:"-"`
	Id    int            `json:"id" xorm:"pk autoincr"`
	Name  string         `json:"name"`
	Desc  string         `json:"host" `
	Value string         `json:"value"`
}

var (
	NoticeSettingStatusEnabled  = "enabled"         //启用
	NoticeSettingStatusDisabled = "disabled"        //禁用
	NoticeOptionScrapyd         = "scrapyd_service" //scrapyd服务监控状态
	NoticeOptionTaskFinished    = "task_finished"   //是否监听任务完成
	NoticeOptionTaskError       = "task_error"      //是否监听任务异常
	NoticeSettingWildcards      = map[string][]string{
		NoticeOptionScrapyd:      {"{host}", "{error_time}", "{error_message}"},
		NoticeOptionTaskFinished: {"{job_id}", "{task_id}", "{host}", "{project}", "{version}", "{spider}", "{start_time}", "{end_time}", "{duration_time}"},
		NoticeOptionTaskError:    {"{job_id}", "{task_id}", "{host}", "{project}", "{version}", "{spider}", "{start_time}", "{error_time}", "{duration_time}", "{error_message}"},
	}
)

func (n *NoticeSetting) Update(fields []core.B) bool {
	if _, err := core.Db.Exec(core.JoinBatchUpdateSql("notice_setting", fields, "name")); err != nil {
		core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, nil, fmt.Sprintf("通知配置修改失败:%s", err))
		return false
	}
	core.NoticeSettings = n.Find()
	return true
}

func (n *NoticeSetting) Find() core.B {
	settings := make([]NoticeSetting, 0)
	if err := core.Db.Select("name,value").OrderBy("id asc").Find(&settings); err != nil {
		core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, nil, fmt.Sprintf("通知配置列表获取失败:%s", err))
	}
	settingMap := make(core.B)
	for _, setting := range settings {
		settingMap[setting.Name] = setting.Value
	}
	return settingMap
}
