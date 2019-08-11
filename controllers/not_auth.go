package controllers

import (
	"github.com/gin-gonic/gin"
	"scrapyd-admin/core"
	"scrapyd-admin/models"
	"strconv"
	"strings"
)

type NotAuth struct {
	core.BaseController
}

//根据项目id获取所有版本号
func (n *NotAuth) GetVersionsByProjectId(c *gin.Context) {
	projectId, _ := strconv.Atoi(c.DefaultQuery("project_id", "0"))
	if projectId <= 0 {
		n.Fail(c, "parameter_error")
		return
	}
	projectHistory := models.ProjectHistory{
		ProjectId: projectId,
	}
	n.Success(c, core.A{
		"projectHistories": projectHistory.FindByProjectId(),
	})
}

func (n *NotAuth) GetSpidersAndServersByProjectId(c *gin.Context) {
	project := c.DefaultQuery("project", "")
	if project == "" {
		n.Fail(c, "parameter_error")
		return
	}
	projectInfo := strings.Split(project, "|")
	projectId, _ := strconv.Atoi(projectInfo[0])
	spider := new(models.Spider)
	server := new(models.Server)
	n.Success(c, core.A{
		"spiders": spider.FindByProjectIdAndVersion(projectId, projectInfo[2]),
		"servers": server.FindByProjectId(projectId),
	})
}
