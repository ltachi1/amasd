package models

import (
	"scrapyd-admin/core"
	"time"
)

type ServerMonitor struct {
	Base           core.BaseModel `json:"-" xorm:"-"`
	Id             int            `json:"id" xorm:"pk autoincr"`
	ServerId       int            `json:"server_id"`
	MemTotal       int64          `json:"mem_total"`
	MemAvailable   int64          `json:"mem_available"`
	MemUsedPercent int            `json:"mem_used_percent"`
	CpuPercent     int            `json:"cpu_percent"`
	CreateTime     core.Timestamp `json:"create_time"`
}

//添加服务器监控数据
func (s *ServerMonitor) InsertOne() {
	core.Db.Table(s).InsertOne(s)
}

//删除一小时之前的监控数据
func (s *ServerMonitor) DelAnHourAgo() {
	core.Db.Where("create_time < ?", time.Now().Unix()-3600).NoAutoCondition().Delete(s)
}

func (s *ServerMonitor) FindByLastTime(lastTime int) []core.A {
	server := Server{
		Monitor: ServerMonitorEnabled,
	}
	serverList := server.Find()
	serverMonitorItems := make([]core.A, 0)
	for _, se := range serverList {
		sm := make([]ServerMonitor, 0)
		core.Db.Where("server_id = ? and create_time > ?", se.Id, lastTime).Find(&sm)
		items := make([]core.A, 0)
		for _, s := range sm {
			items = append(items, core.A{
				"mem_used_percent": s.MemUsedPercent,
				"cpu_percent":      s.CpuPercent,
				"time":             s.CreateTime,
				"id":               s.Id,
			})
		}
		serverMonitorItems = append(serverMonitorItems, core.A{"name": se.MonitorAddress, "timeline": items})
	}

	return serverMonitorItems
}
