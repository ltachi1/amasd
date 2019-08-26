package models

import (
	"amasd/core"
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
	MemUsed         int64          `json:"mem_used"`
	MemUsedPercent  string         `json:"mem_used_percent"`
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

func (s *ServerMonitor) FindByLastTime(serverId int, lastTime int) []core.A {
	sm := make([]ServerMonitor, 0)
	core.Db.Where("server_id = ? and create_time > ?", serverId, lastTime).Find(&sm)
	items := make([]core.A, 0)
	for _, s := range sm {
		cpuLoad1, _ := strconv.ParseFloat(s.CpuLoad1, 32)
		cpuLoad5, _ := strconv.ParseFloat(s.CpuLoad5, 32)
		cpuLoad15, _ := strconv.ParseFloat(s.CpuLoad15, 32)
		cpuLoad1Percent := (float32(cpuLoad1)/float32(s.CpuCoreCount)) * 100
		cpuLoad5Percent := (float32(cpuLoad5)/float32(s.CpuCoreCount)) * 100
		cpuLoad15Percent := (float32(cpuLoad15)/float32(s.CpuCoreCount)) * 100
		if cpuLoad1Percent > 100 {
			cpuLoad1Percent = 100
		}
		if cpuLoad5Percent > 100 {
			cpuLoad5Percent = 100
		}
		if cpuLoad15Percent > 100 {
			cpuLoad15Percent = 100
		}

		items = append(items, core.A{
			"mem_used_percent":  s.MemUsedPercent,
			"mem_used":          fmt.Sprintf("%.2f", float32(s.MemUsed)/1073741824),
			"mem_total":         fmt.Sprintf("%.2f", float32(s.MemTotal)/1073741824),
			"cpu_percent":       s.CpuPercent,
			"cpu_core_count":    s.CpuCoreCount,
			"cpu_load1":         fmt.Sprintf("%.2f", cpuLoad1Percent),
			"cpu_load5":         fmt.Sprintf("%.2f", cpuLoad5Percent),
			"cpu_load15":        fmt.Sprintf("%.2f", cpuLoad15Percent),
			"process_count":     s.ProcessCount,
			"net_send_speed":    fmt.Sprintf("%.2f", float32(s.NetSendSpeed)/1048576),
			"net_receive_speed": fmt.Sprintf("%.2f", float32(s.NetReceiveSpeed)/1048576),
			"time":              s.CreateTime,
		})
	}

	return items
}

//获取服务器性能指标概览
func (s *ServerMonitor) OverviewByIds(ids []string) ([]map[string]string, error) {
	//获取服务器状态
	servers := new(Server).FindByIds(core.StringArrayToInt(ids))
	statusList := map[string]string{}
	for _, server := range servers {
		statusList[strconv.Itoa(server.Id)] = strconv.Itoa(int(server.AgentStatus))
	}
	monitorList, err := core.Db.QueryString("select id,server_id,mem_total,mem_available,mem_used,mem_used_percent,cpu_percent,create_time,cpu_core_count,cpu_load15 " +
		"from server_monitor where id in(select max(id) from server_monitor group by server_id order by create_time)")
	if err == nil {
		for i := 0; i < len(monitorList); i++ {
			memAvailable, _ := strconv.ParseFloat(monitorList[i]["mem_available"], 32)
			memUsed, _ := strconv.ParseFloat(monitorList[i]["mem_used"], 32)
			memTotal, _ := strconv.ParseFloat(monitorList[i]["mem_total"], 32)
			cpuLoad15, _ := strconv.ParseFloat(monitorList[i]["cpu_load15"], 32)
			cpuCoreCount, _ := strconv.ParseFloat(monitorList[i]["cpu_core_count"], 32)
			memAvailable = memAvailable / 1073741824
			memUsed = memUsed / 1073741824
			memTotal = memTotal / 1073741824
			cpuLoad := (cpuLoad15/cpuCoreCount) * 100
			if cpuLoad > 100 {
				cpuLoad = 100
			}

			monitorList[i]["mem_available"] = fmt.Sprintf("%.1f", memAvailable)
			monitorList[i]["mem_used"] = fmt.Sprintf("%.1f", memUsed)
			monitorList[i]["mem_total"] = fmt.Sprintf("%.1f", memTotal)
			monitorList[i]["cpu_load"] = fmt.Sprintf("%.2f", cpuLoad)
			monitorList[i]["agent_status"] = statusList[monitorList[i]["server_id"]]
		}
	}
	if monitorList == nil {
		monitorList = make([]map[string]string, 0)
	}
	return monitorList, err
}
