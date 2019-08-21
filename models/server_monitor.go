package models

import (
	"scrapyd-admin/core"
	"time"
	"fmt"
	"strconv"
)

type ServerMonitor struct {
	Base            core.BaseModel `json:"-" xorm:"-"`
	Id              int            `json:"id" xorm:"pk autoincr"`
	ServerId        int            `json:"server_id"`
	MemTotal        int64          `json:"mem_total"`
	MemAvailable    int64          `json:"mem_available"`
	MemUsedPercent  int            `json:"mem_used_percent"`
	CpuPercent      string         `json:"cpu_percent"`
	CpuCoreCount    int            `json:"cpu_core_count"`
	CpuLoad1        string         `json:"cpu_load1"`
	CpuLoad5        string         `json:"cpu_load5"`
	CpuLoad15       string         `json:"cpu_load15"`
	ProcessCount    int            `json:"process_count"`
	NetSendSpeed    int            `json:"net_send_speed"`
	NetReceiveSpeed int            `json:"net_receive_speed"`
	CreateTime      core.Timestamp `json:"create_time"`
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
	for i, se := range serverList {
		sm := make([]ServerMonitor, 0)
		core.Db.Where("server_id = ? and create_time > ?", se.Id, lastTime).Find(&sm)
		items := make([]core.A, 0)
		for _, s := range sm {
			items = append(items, core.A{
				"mem_used_percent":  s.MemUsedPercent,
				"cpu_percent":       s.CpuPercent,
				"cpu_core_count":    s.CpuCoreCount,
				"cpu_load5":         s.CpuLoad5,
				"process_count":     s.ProcessCount,
				"net_send_speed":    fmt.Sprintf("%.2f", float32(s.NetSendSpeed)/1048576),
				"net_receive_speed": fmt.Sprintf("%.2f", float32(s.NetReceiveSpeed)/1048576),
				"time":              s.CreateTime,
			})
		}
		serverMonitorItems = append(serverMonitorItems, core.A{"name": fmt.Sprintf("服务器%d", i), "timeline": items})
	}

	return serverMonitorItems
}

//获取服务器性能指标概览
func (s *ServerMonitor) OverviewByIds(ids []string) ([]map[string]string, error) {
	monitorList, err := core.Db.QueryString("select id,server_id,mem_total,mem_available,mem_used_percent,cpu_percent,create_time,cpu_core_count,cpu_load15 " +
		"from server_monitor where id in(select max(id) from server_monitor group by server_id order by create_time)")
	if err == nil {
		for i := 0; i < len(monitorList); i++ {
			memAvailable, _ := strconv.ParseFloat(monitorList[i]["mem_available"], 32)
			memTotal, _ := strconv.ParseFloat(monitorList[i]["mem_total"], 32)
			monitorList[i]["mem_available"] = fmt.Sprintf("%.2f", memAvailable/1048576)
			monitorList[i]["mem_total"] = fmt.Sprintf("%.2f", memTotal/1048576)
		}
	}
	return monitorList, err
}
