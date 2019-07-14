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
	count, _ := core.Db.Where("project_id = ? and status = ?", projectId, SchedulesTaskStatusEnabled).Table("task").Count()
	if count > 0 {
		return true
	}
	return false
}

func (s *SchedulesTask) HaveEnabledByServer(serverId int) bool {
	count, _ := core.Db.Where("server_id = ? and status = ?", serverId, SchedulesTaskStatusEnabled).Table("task").Count()
	if count > 0 {
		return true
	}
	return false
}

func (s *SchedulesTask) Inert(projectId int, projectName string, version string, cron string, spiderList []string, serverList []string) bool {
	schedulesTaskList := make([]SchedulesTask, 0)
	session := core.Db.NewSession()
	defer session.Close()
	session.Begin()
	for _, spider := range spiderList {
		sp := strings.Split(spider, "|")
		spiderId, _ := strconv.Atoi(sp[0])
		for _, server := range serverList {
			se := strings.Split(server, "|")
			serverId, _ := strconv.Atoi(se[0])
			st := SchedulesTask{
				ProjectId:   projectId,
				ProjectName: projectName,
				Version:     version,
				SpiderId:    spiderId,
				SpiderName:  sp[1],
				ServerId:    serverId,
				Host:        se[1],
				Cron:        cron,
				Status:      SchedulesTaskStatusEnabled,
			}
			if _, err :=session.InsertOne(&st); err != nil {
				session.Rollback()
				return false
			}
			schedulesTaskList = append(schedulesTaskList, st)
		}
	}

	if len(schedulesTaskList) > 0 {
		session.Commit()
		//添加定时任务
		for i := 0; i < len(schedulesTaskList); i++ {
			core.Cron.AddFunc(schedulesTaskList[i].Cron, schedulesTaskList[i].RunSchedules, strconv.Itoa(schedulesTaskList[i].Id))
		}
	} else {
		session.Rollback()
	}

	return true
}

//分页获取计划任务列表
func (s *SchedulesTask) FindPages(projectId int, version string, serverId int, status string, page int, pageSize int) ([]core.B, int) {
	tasks := make([]core.B, 0)
	countObj := core.Db.Table("schedules_task")
	selectObj := core.Db.Table("schedules_task")
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
	if _, error := core.Db.Id(s.Id).Update(s); error != nil {
		return false
	}
	if ok, _ := core.Db.Id(s.Id).Get(s); ok {
		if s.Status == SchedulesTaskStatusEnabled {
			core.Cron.AddFunc(s.Cron, s.RunSchedules, strconv.Itoa(s.Id))
		} else if s.Status == SchedulesTaskStatusDisabled {
			core.Cron.RemoveJob(strconv.Itoa(s.Id))
		}
	}

	return true
}

func (s *SchedulesTask) Del(id int) bool {
	if _, err := core.Db.Where("id = ?", id).NoAutoCondition().Delete(&SchedulesTask{}); err != nil {
		return false
	}
	core.Cron.RemoveJob(strconv.Itoa(s.Id))
	return true
}

//初始化计划任务
func (s *SchedulesTask) InitSchedulesToCron() bool {
	schedulesTaskList := make([]SchedulesTask, 0)
	core.Db.Where("status = ? ", SchedulesTaskStatusEnabled).Find(&schedulesTaskList)
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
