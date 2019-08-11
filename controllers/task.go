package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"scrapyd-admin/core"
	"scrapyd-admin/models"
	"strconv"
	"strings"
	"encoding/json"
)

type Task struct {
	core.BaseController
}

func (t *Task) Index(c *gin.Context) {
	if core.IsAjax(c) {
		projectId, _ := strconv.Atoi(c.DefaultPostForm("project_id", "0"))
		version := c.DefaultPostForm("version", "")
		serverId, _ := strconv.Atoi(c.DefaultPostForm("server_id", "0"))
		status := c.DefaultPostForm("status", "")
		page, _ := strconv.Atoi(c.DefaultPostForm("pagination[page]", "1"))
		pageLength, _ := strconv.Atoi(c.DefaultPostForm("pagination[perpage]", "10"))
		tasks, totalCount := new(models.Task).FindTaskPages(projectId, version, serverId, status, page, pageLength)
		c.JSON(http.StatusOK, gin.H{
			"data": tasks,
			"meta": gin.H{
				"page":    page,
				"total":   totalCount,
				"pages":   core.CalculationPages(totalCount, pageLength),
				"perpage": pageLength,
			},
		})
	} else {
		c.HTML(http.StatusOK, "task/index", gin.H{
			"projects": new(models.Project).Find(),
			"servers":  new(models.Server).Find(),
		})
	}

}

func (t *Task) Add(c *gin.Context) {
	if core.IsAjax(c) {
		project := c.DefaultPostForm("project", "")
		spiders := c.PostFormArray("spider")
		servers := c.PostFormArray("server")
		if project == "" {
			t.Fail(c, "parameter_error")
			return
		}
		projectInfo := strings.Split(project, "|")
		if len(projectInfo) != 3 {
			t.Fail(c, "parameter_error")
			return
		}
		projectId, _ := strconv.Atoi(projectInfo[0])
		if projectId <= 0 {
			t.Fail(c, "parameter_error")
			return
		}
		if ok, errorTaskList := new(models.Task).Inert(projectId, projectInfo[1], projectInfo[2], spiders, servers); ok {
			t.Success(c, nil)
		} else {
			t.Fail(c, "task_add_error", strings.Join(errorTaskList, ", "))
		}
	} else {
		project := new(models.Project)
		c.HTML(http.StatusOK, "task/add", gin.H{
			"projects": project.Find(),
		})
	}
}

func (t *Task) Cancel(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if id == 0 {
			t.Fail(c, "parameter_error")
			return
		}
		if new(models.Task).Cancel(id) {
			t.Success(c, nil)
		} else {
			t.Fail(c, "update_error")
		}
	}
}

func (t *Task) CancelMulti(c *gin.Context) {
	if core.IsAjax(c) {
		ids := c.DefaultPostForm("ids", "")
		if ids == "" {
			t.Fail(c, "parameter_error")
			return
		}
		idList := make([]string, 0)
		if err := json.Unmarshal([]byte(ids), &idList); err != nil {
			t.Fail(c, "parameter_error")
			return
		}
		if ok, failureList := new(models.Task).CancelMulti(idList); ok {
			t.Success(c, nil)
		} else {
			t.Fail(c, "task_update_error", strings.Join(failureList, ", "))
		}
	}
}

func (t *Task) CancelAll(c *gin.Context) {
	if core.IsAjax(c) {
		projectId, _ := strconv.Atoi(c.DefaultPostForm("project_id", "0"))
		version := c.DefaultPostForm("version", "")
		serverId, _ := strconv.Atoi(c.DefaultPostForm("server_id", "0"))
		status := c.DefaultPostForm("status", "")
		if ok, failureList := new(models.Task).CancelAll(projectId, version, serverId, status); ok {
			t.Success(c, nil)
		} else {
			t.Fail(c, "task_update_error", strings.Join(failureList, ", "))
		}
	}
}

func (t *Task) Del(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if id == 0 {
			t.Fail(c, "parameter_error")
			return
		}
		if new(models.Task).Del(id) {
			t.Success(c, nil)
		} else {
			t.Fail(c, "del_error")
		}
	}
}

func (t *Task) DelMulti(c *gin.Context) {
	if core.IsAjax(c) {
		ids := c.DefaultPostForm("ids", "")
		if ids == "" {
			t.Fail(c, "parameter_error")
			return
		}
		idList := make([]string, 0)
		if err := json.Unmarshal([]byte(ids), &idList); err != nil {
			t.Fail(c, "parameter_error")
			return
		}
		if new(models.Task).DelMulti(idList) {
			t.Success(c, nil)
		} else {
			t.Fail(c, "del_error")
		}
	}
}

func (t *Task) DelAll(c *gin.Context) {
	if core.IsAjax(c) {
		projectId, _ := strconv.Atoi(c.DefaultPostForm("project_id", "0"))
		version := c.DefaultPostForm("version", "")
		serverId, _ := strconv.Atoi(c.DefaultPostForm("server_id", "0"))
		status := c.DefaultPostForm("status", "")
		if new(models.Task).DelAll(projectId, version, serverId, status) {
			t.Success(c, nil)
		} else {
			t.Fail(c, "del_error")
		}
	}
}


//计划列表
func (t *Task) Schedules(c *gin.Context) {
	if core.IsAjax(c) {
		projectId, _ := strconv.Atoi(c.DefaultPostForm("project_id", "0"))
		version := c.DefaultPostForm("version", "")
		serverId, _ := strconv.Atoi(c.DefaultPostForm("server_id", "0"))
		status := c.DefaultPostForm("status", "")
		page, _ := strconv.Atoi(c.DefaultPostForm("pagination[page]", "1"))
		pageLength, _ := strconv.Atoi(c.DefaultPostForm("pagination[perpage]", "10"))
		tasks, totalCount := new(models.SchedulesTask).FindPages(projectId, version, serverId, status, page, pageLength)
		c.JSON(http.StatusOK, gin.H{
			"data": tasks,
			"meta": gin.H{
				"page":    page,
				"total":   totalCount,
				"pages":   core.CalculationPages(totalCount, pageLength),
				"perpage": pageLength,
			},
		})
	} else {
		c.HTML(http.StatusOK, "task/schedules", gin.H{
			"projects": new(models.Project).Find(),
			"servers":  new(models.Server).Find(),
		})
	}

}

func (t *Task) AddSchedules(c *gin.Context) {
	if core.IsAjax(c) {
		project := c.DefaultPostForm("project", "")
		spiders := c.PostFormArray("spider")
		servers := c.PostFormArray("server")
		cron := c.DefaultPostForm("cron", "")
		if project == "" {
			t.Fail(c, "parameter_error")
			return
		}
		if cron == "" {
			t.Fail(c, "parameter_error")
			return
		}
		projectInfo := strings.Split(project, "|")
		if len(projectInfo) != 3 {
			t.Fail(c, "parameter_error")
			return
		}
		projectId, _ := strconv.Atoi(projectInfo[0])
		if projectId <= 0 {
			t.Fail(c, "parameter_error")
			return
		}
		if new(models.SchedulesTask).Inert(projectId, projectInfo[1], projectInfo[2], cron, spiders, servers) {
			t.Success(c, nil)
		} else {
			t.Fail(c, "add_error")
		}
	} else {
		project := new(models.Project)
		c.HTML(http.StatusOK, "task/add_schedules", gin.H{
			"projects": project.Find(),
		})
	}
}

func (t *Task) UpdateSchedulesStatus(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
		status := c.DefaultPostForm("status", "")
		if id == 0 {
			t.Fail(c, "parameter_error")
			return
		}
		if status != models.SchedulesTaskStatusEnabled && status != models.SchedulesTaskStatusDisabled {
			t.Fail(c, "parameter_error")
			return
		}

		schedulesTask := models.SchedulesTask{
			Id:     id,
			Status: status,
		}
		if schedulesTask.UpdateStatus() {
			t.Success(c, nil)
		} else {
			t.Fail(c, "update_error")
		}
	}
}

func (t *Task) DelSchedules(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if new(models.SchedulesTask).Del(id) {
			t.Success(c, nil)
		} else {
			t.Fail(c, "del_error")
		}
	}
}
