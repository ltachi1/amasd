package models

import (
	"fmt"
	"scrapyd-admin/core"
	"strconv"
	"strings"
	"time"
	"github.com/ltachi1/logrus"
	"sync"
)

type Task struct {
	Base        core.BaseModel `json:"-" xorm:"-"`
	Id          int            `json:"id" xorm:"pk autoincr"`
	Type        uint8          `json:"type"`
	ProjectId   int            `json:"project_id"`
	ProjectName string         `json:"project_name"`
	Version     string         `json:"version"`
	ServerId    int            `json:"server_id"`
	Host        string         `json:"host"`
	SpiderId    int            `json:"spider_id"`
	SpiderName  string         `json:"spider_name"`
	JobId       string         `json:"job_id"`
	StartTime   core.Timestamp `json:"start_time"`
	EndTime     core.Timestamp `json:"end_time"`
	Status      string         `json:"status"`
}

//Status 运行状态 error | pending | running | cancelled | finished
var (
	TaskTypeOnce        uint8 = 1 //一次性任务
	TaskTypeTiming      uint8 = 2 //定时任务
	TaskStatusError           = "error"
	TaskStatusPending         = "pending"
	TaskStatusRunning         = "running"
	TaskStatusCancelled       = "cancelled"
	TaskStatusFinished        = "finished"
)

//查询项目下是否有正在运行的爬虫，包括定时任务
func (t *Task) HaveRunning(projectId int) bool {
	count, _ := core.Db.Where("project_id = ? and (status = ? or status = ?)", projectId, TaskStatusPending, TaskStatusRunning).Table("task").Count()
	if count > 0 {
		return true
	}
	return false
}

func (t *Task) HaveRunningByServer(serverId int) bool {
	count, _ := core.Db.Where("server_id = ? and (status = ? or status = ?)", serverId, TaskStatusPending, TaskStatusRunning).Table("task").Count()
	if count > 0 {
		return true
	}
	return false
}

func (t *Task) Inert(projectId int, projectName string, version string, spiderList []string, serverList []string) (bool, []string) {
	errorTaskList := make([]string, 0)
	for _, spider := range spiderList {
		sp := strings.Split(spider, "|")
		spiderId, _ := strconv.Atoi(sp[0])
		for _, server := range serverList {
			se := strings.Split(server, "|")
			serverId, _ := strconv.Atoi(se[0])
			task := Task{
				Type:        TaskTypeOnce,
				ProjectId:   projectId,
				ProjectName: projectName,
				Version:     version,
				SpiderId:    spiderId,
				SpiderName:  sp[1],
				ServerId:    serverId,
				Host:        se[1],
				Status:      TaskStatusPending,
			}
			_, error := core.Db.InsertOne(&task)
			if error == nil {
				t.RunTask(task.Id)
			} else {
				errorTaskList = append(errorTaskList, fmt.Sprintf("%s - %s", se[1], sp[1]))
			}
		}
	}
	if len(errorTaskList) > 0 {
		return false, errorTaskList
	}
	return true, errorTaskList
}

func (t *Task) InertOne() {
	if _, error := core.Db.InsertOne(t); error == nil {
		t.RunTask(t.Id)
	}
}

//运行任务
func (t *Task) RunTask(taskId int) {
	task := core.B{}
	if ok, _ := core.Db.Select("t.*,s.auth,s.username,s.password,s.status").Table("task").Alias("t").Join("INNER", "server as s", "t.server_id = s.id").Where("t.id = ?", taskId).Limit(1).Get(&task); ok {
		serverStatus, _ := strconv.Atoi(task["status"])
		if uint8(serverStatus) == ServerStatusNormal {
			auth, _ := strconv.Atoi(task["auth"])
			scrapyd := Scrapyd{
				Host: task["host"],
			}
			if uint8(auth) == ServerAuthOpen {
				scrapyd.Username = task["username"]
				scrapyd.Password = task["password"]
			}
			if ok, jobId := scrapyd.Schedule(task["project_name"], task["version"], task["spider_name"]); ok {
				core.Db.Id(taskId).Update(&Task{Status: TaskStatusRunning, JobId: jobId, StartTime: core.Timestamp(time.Now().Unix())})
			} else {
				core.Db.Id(taskId).Update(&Task{Status: TaskStatusError, StartTime: core.Timestamp(time.Now().Unix()), EndTime: core.Timestamp(time.Now().Unix())})
			}
		} else {
			core.Db.Id(taskId).Update(&Task{Status: TaskStatusError, StartTime: core.Timestamp(time.Now().Unix()), EndTime: core.Timestamp(time.Now().Unix())})
		}
	}
}

//分页获取任务列表
func (t *Task) FindTaskPages(projectId int, version string, serverId int, status string, page int, pageSize int) ([]core.B, int) {
	tasks := make([]core.B, 0)
	countObj := core.Db.Table("task")
	selectObj := core.Db.Table("task")
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
	for i := 0; i < len(tasks); i++ {
		startTimestamp, _ := strconv.Atoi(tasks[i]["start_time"])
		if tasks[i]["start_time"] != "0" {
			tasks[i]["start_time"] = core.FormatDateByString(tasks[i]["start_time"], "2006-01-02 15:04:05")
		} else {
			tasks[i]["start_time"] = ""
		}

		if tasks[i]["end_time"] == "0" {
			endTimestamp := int(time.Now().Unix())
			tasks[i]["duration"] = core.TimeDifference(startTimestamp, endTimestamp)
		} else {
			endTimestamp, _ := strconv.Atoi(tasks[i]["end_time"])
			tasks[i]["duration"] = core.TimeDifference(startTimestamp, endTimestamp)
		}

	}
	return tasks, int(totalCount)
}

//取消单个任务
func (t *Task) Cancel(id int) bool {
	if t.ProjectName == "" || t.JobId == "" {
		if ok, _ := core.Db.Select("id,project_name,server_id,job_id").Where("id = ? and (status = ? or status = ?)", id, TaskStatusPending, TaskStatusRunning).NoAutoCondition().Get(t); !ok {
			return false
		}
	}
	if t.Id == 0 {
		return false
	}
	server := new(Server)
	if !server.Get(t.ServerId) {
		return false
	}
	if server.Id == 0 {
		return false
	}
	scrapyd := Scrapyd{
		Host:     server.Host,
		Auth:     server.Auth,
		Username: server.Username,
		Password: server.Password,
	}
	if !scrapyd.Cancel(t.ProjectName, t.JobId) {
		return false
	}
	if _, error := core.Db.Id(id).Update(&Task{Status: TaskStatusCancelled, EndTime: core.Timestamp(time.Now().Unix())}); error != nil {
		return false
	}
	return true
}

//取消多个任务
func (t *Task) CancelMulti(ids []string) (bool, []string) {
	failureList := make([]string, 0)
	var wg sync.WaitGroup
	for _, id := range ids {
		if i, err := strconv.Atoi(id); err == nil {
			wg.Add(1)
			go func(i int) {
				if !new(Task).Cancel(i) {
					failureList = append(failureList, strconv.Itoa(i))
				}
				wg.Done()
			}(i)
		}
	}
	wg.Wait()
	if len(failureList) == 0 {
		return true, failureList
	}

	return false, failureList
}

func (t *Task) CancelAll(projectId int, version string, serverId int, status string) (bool, []string) {
	failureList := make([]string, 0)
	tasks := make([]Task, 0)
	obj := core.Db.Table("task")
	if projectId > 0 {
		obj.Where("project_id = ? ", projectId)
	}
	if version != "" {
		obj.Where("version = ? ", version)
	}
	if serverId > 0 {
		obj.Where("server_id = ? ", serverId)
	}
	if status != "" {
		obj.Where("status = ? ", status)
	}
	obj.OrderBy("id asc").Find(&tasks)
	if len(tasks) > 0 {
		var wg sync.WaitGroup
		for _, task := range tasks {
			wg.Add(1)
			go func(task Task) {
				if !task.Cancel(task.Id) {
					failureList = append(failureList, strconv.Itoa(task.Id))
				}
				wg.Done()
			}(task)
		}
		wg.Wait()
		if len(failureList) > 0 {
			return false, failureList
		}
	}
	return true, failureList
}

//删除单个任务
func (t *Task) Del(id int) bool {
	if _, error := core.Db.Id(id).Delete(t); error != nil {
		return false
	}
	return true
}

//删除多个任务
func (t *Task) DelMulti(ids []string) bool {
	if _, error := core.Db.In("id", ids).Delete(t); error != nil {
		return false
	}
	return true
}

func (t *Task) DelAll(projectId int, version string, serverId int, status string) bool {
	obj := core.Db.Where("1 = 1")
	if projectId > 0 {
		obj = obj.Where("project_id = ? ", projectId)
	}
	if version != "" {
		obj.Where("version = ? ", version)
	}
	if serverId > 0 {
		obj.Where("server_id = ? ", serverId)
	}
	if status != "" {
		obj.Where("status = ? ", status)
	}
	if _, error := obj.NoAutoCondition().Delete(t); error != nil {
		return false
	}
	return true
}

func (t *Task) DetectionStatus() {
	taskList := make([]core.B, 0)
	core.Db.Select("t.id,t.job_id,t.project_name,s.auth,s.username,s.password,s.status,s.host").Table("task").Alias("t").Join("INNER", "server as s", "t.server_id = s.id").Where("(t.status = ? or t.status = ?) and t.job_id<>\"\"", TaskStatusPending, TaskStatusRunning).Find(&taskList)
	serverProjectList := make(map[string]core.B, 0)
	serverProjectTaskList := make(map[string][]core.B, 0)
	for _, task := range taskList {
		key := fmt.Sprintf("%s_%s", task["host"], task["project_name"])
		if _, exists := serverProjectList[key]; exists {
			serverProjectTaskList[key] = append(serverProjectTaskList[key], core.B{"id": task["id"], "job_id": task["job_id"], "status": task["status"], "spider_name": task["spider_name"]})
		} else {
			serverProjectList[key] = core.B{"project_name": task["project_name"], "host": task["host"], "auth": task["auth"], "username": task["username"], "password": task["password"]}
			serverProjectTaskList[key] = []core.B{{"id": task["id"], "job_id": task["job_id"], "status": task["status"], "spider_name": task["spider_name"]}}
		}
	}
	updateTask := make([]core.B, 0)
	var wg sync.WaitGroup
	for k, sp := range serverProjectList {
		wg.Add(1)
		go func(k string, sp core.B) {
			auth, _ := strconv.Atoi(sp["auth"])
			scrapyd := Scrapyd{
				Host:     sp["host"],
				Auth:     uint8(auth),
				Username: sp["username"],
				Password: sp["password"],
			}
			taskStatusList := map[string]core.B{}
			if ok, result := scrapyd.ListJobs(sp["project_name"]); ok {
				for _, v := range result["pending"] {
					job := v.(map[string]interface{})
					taskStatusList[job["id"].(string)] = core.B{
						"status": TaskStatusPending,
					}
				}
				for _, v := range result["running"] {
					job := v.(map[string]interface{})
					taskStatusList[job["id"].(string)] = core.B{
						"status":     TaskStatusRunning,
						"start_time": job["start_time"].(string),
					}
				}
				for _, v := range result["finished"] {
					job := v.(map[string]interface{})
					taskStatusList[job["id"].(string)] = core.B{
						"status":     TaskStatusFinished,
						"start_time": job["start_time"].(string),
						"end_time":   job["end_time"].(string),
					}
				}
			}

			for _, task := range serverProjectTaskList[k] {
				if taskStatus, exists := taskStatusList[task["job_id"]]; exists {
					if taskStatus["status"] == TaskStatusPending {
						updateTask = append(updateTask, core.B{
							"id":     task["id"],
							"status": TaskStatusPending,
						})
					} else if taskStatus["status"] == TaskStatusRunning {
						updateTask = append(updateTask, core.B{
							"id":         task["id"],
							"status":     TaskStatusRunning,
							"start_time": strconv.Itoa(core.DateToTimestamp(taskStatus["start_time"])),
						})
					} else if taskStatus["status"] == TaskStatusFinished {
						updateTask = append(updateTask, core.B{
							"id":         task["id"],
							"status":     TaskStatusFinished,
							"start_time": strconv.Itoa(core.DateToTimestamp(taskStatus["start_time"])),
							"end_time":   strconv.Itoa(core.DateToTimestamp(taskStatus["end_time"])),
						})
					}
				} else {
					updateTask = append(updateTask, core.B{
						"id":       task["id"],
						"status":   TaskStatusError,
						"end_time": strconv.Itoa(int(time.Now().Unix())),
					})
				}
			}
			wg.Done()
		}(k, sp)
	}
	wg.Wait()
	if len(updateTask) > 0 {
		if _, error := core.Db.Exec(core.JoinBatchUpdateSql("task", updateTask, "id")); error != nil {
			core.WriteLog(core.LogTypeTask, logrus.ErrorLevel, nil, error)
		}
	}
}
