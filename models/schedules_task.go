package models

import (
	"scrapyd-admin/core"
	"strconv"
	"strings"
)

type SchedulesTask struct {
	Base        core.BaseModel `json:"-" xorm:"-"`
	Id          int            `json:"id" xorm:"pk autoincr"`
	ProjectId   int            `json:"project_id"`
	ProjectName string         `json:"project_name"`
	Version     string         `json:"version"`
	ServerId    int            `json:"server_id"`
	Host        string         `json:"host"`
	SpiderId    int            `json:"spider_id"`
	SpiderName  string         `json:"spider_name"`
	Cron        string         `json:"cronn"`
	Status      string         `json:"status"`
}

var (
	SchedulesTaskStatusEnabled  = "enabled"
	SchedulesTaskStatusDisabled = "disabled"
)

//查询项目下是有正在启用的计划任务
func (s *SchedulesTask) HaveEnabled(projectId int) bool {
	count, _ := core.DBPool.Slave().Where("project_id = ? and status = ?", projectId, SchedulesTaskStatusEnabled).Table("task").Count()
	if count > 0 {
		return true
	}
	return false
}

func (s *SchedulesTask) Inert(projectId int, projectName string, version string, cron string, spiderList []string, serverList []string) bool {
	schedulesTaskList := make([]SchedulesTask, 0)
	for _, spider := range spiderList {
		sp := strings.Split(spider, "|")
		spiderId, _ := strconv.Atoi(sp[0])
		for _, server := range serverList {
			se := strings.Split(server, "|")
			serverId, _ := strconv.Atoi(se[0])
			schedulesTaskList = append(schedulesTaskList, SchedulesTask{
				ProjectId:   projectId,
				ProjectName: projectName,
				Version:     version,
				SpiderId:    spiderId,
				SpiderName:  sp[1],
				ServerId:    serverId,
				Host:        se[1],
				Cron:        cron,
				Status:      SchedulesTaskStatusEnabled,
			})
		}
	}

	if len(schedulesTaskList) > 0 {
		_, error := core.DBPool.Master().Insert(&schedulesTaskList)
		if error != nil {
			return false
		}
		//添加定时任务
		for _, st := range schedulesTaskList{
			core.Cron.AddFunc(st.Cron, st.RunSchedules, strconv.Itoa(st.Id))
			st.RunSchedules()
		}
	}

	return true
}

//分页获取计划任务列表
func (s *SchedulesTask) FindPages(projectId int, version string, serverId int, status string, page int, pageSize int) ([]core.B, int) {
	tasks := make([]core.B, 0)
	countObj := core.DBPool.Slave().Table("schedules_task")
	selectObj := core.DBPool.Slave().Table("schedules_task")
	if projectId > 0 {
		countObj.Where("project_id = ? ", projectId)
		selectObj.Where("project_id = ? ", projectId)
	}
	if version != "" {
		countObj.Where("version = ? ", version)
		selectObj.Where("version = ? ", version)
	}
	if serverId > 0 {
		countObj.Where("server_id = ? ", serverId)
		selectObj.Where("server_id = ? ", serverId)
	}
	if status != "" {
		countObj.Where("status = ? ", status)
		selectObj.Where("status = ? ", status)
	}
	totalCount, _ := countObj.Count()
	selectObj.OrderBy("id asc").Limit(pageSize, (page-1)*pageSize).Find(&tasks)
	return tasks, int(totalCount)
}

//修改状态
func (s *SchedulesTask) UpdateStatus() bool {
	if _, error := core.DBPool.Master().Id(s.Id).Update(s); error != nil {
		return false
	}
	if ok, _ := core.DBPool.Slave().Id(s.Id).Get(s); ok {
		if s.Status == SchedulesTaskStatusEnabled {
			core.Cron.AddFunc(s.Cron, s.RunSchedules, strconv.Itoa(s.Id))
		} else if s.Status == SchedulesTaskStatusDisabled {
			core.Cron.RemoveJob(strconv.Itoa(s.Id))
		}
	}

	return true
}

//初始化计划任务
func (s *SchedulesTask) InitSchedulesToCron() bool {
	schedulesTaskList := make([]SchedulesTask, 0)
	core.DBPool.Slave().Where("status = ? ", SchedulesTaskStatusEnabled).Find(&schedulesTaskList)
	for i := 0; i < len(schedulesTaskList); i++ {
		core.Cron.AddFunc(schedulesTaskList[i].Cron, schedulesTaskList[i].RunSchedules, strconv.Itoa(schedulesTaskList[i].Id))
	}
	return true
}

//执行计划任务
func (s *SchedulesTask) RunSchedules() {
	task := Task{
		Type:        TaskTypeTiming,
		ProjectId:   s.ProjectId,
		ProjectName: s.ProjectName,
		Version:     s.Version,
		SpiderId:    s.SpiderId,
		SpiderName:  s.SpiderName,
		ServerId:    s.ServerId,
		Host:        s.Host,
		Status:      TaskStatusPending,
	}
	go task.InertOne()
}
