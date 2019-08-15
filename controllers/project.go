package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"scrapyd-admin/core"
	"scrapyd-admin/models"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Project struct {
	core.BaseController
}

func (p *Project) Index(c *gin.Context) {
	if core.IsAjax(c) {
		serverId, _ := strconv.Atoi(c.DefaultQuery("server_id", "0"))
		page, _ := strconv.Atoi(c.DefaultPostForm("pagination[page]", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultPostForm("pagination[perpage]", "10"))
		project := new(models.Project)
		projects, totalCount := project.GetPageProjects(serverId, page, pageSize)
		c.JSON(http.StatusOK, gin.H{
			"data": projects,
			"meta": gin.H{
				"page":    page,
				"total":   totalCount,
				"pages":   core.CalculationPages(totalCount, pageSize),
				"perpage": pageSize,
			},
		})

	} else {
		server := new(models.Server)
		c.HTML(http.StatusOK, "project/index", gin.H{
			"servers": server.Find(),
		})
	}

}

func (p *Project) Add(c *gin.Context) {
	if core.IsAjax(c) {
		name := strings.Trim(c.DefaultPostForm("name", ""), " ")
		desc := strings.Trim(c.DefaultPostForm("desc", ""), " ")
		lastVersion := strings.Trim(c.DefaultPostForm("lastVersion", ""), " ")
		relation := strings.Trim(c.DefaultPostForm("relation", "no"), "")
		serverIds, _ := c.GetPostFormArray("serverIds")
		file, _ := c.FormFile("customFile")
		if name == "" || relation == "" {
			p.Fail(c, "parameter_error")
			return
		}
		if utf8.RuneCountInString(name) > 20 {
			p.Fail(c, "extra_long_error", "项目名称", "20")
			return
		}
		if utf8.RuneCountInString(name) > 500 {
			p.Fail(c, "extra_long_error", "项目描述", "500")
			return
		}
		if utf8.RuneCountInString(lastVersion) > 20 {
			p.Fail(c, "extra_long_error", "版本号", "20")
			return
		}
		if relation == "yes" && len(serverIds) == 0 {
			p.Fail(c, "project_server_error")
			return
		}
		if relation == "yes" && file == nil {
			p.Fail(c, "file_upload_error")
			return
		}
		project := models.Project{
			Name:        name,
			Desc:        desc,
			LastVersion: lastVersion,
		}
		if ok, str, errorServerList := project.InsertOne(relation, core.StringArrayToInt(serverIds), file); ok {
			p.Success(c, nil)
		} else {
			if len(errorServerList) > 0 {
				p.Fail(c, str, strings.Join(errorServerList, ", "))
			} else {
				p.Fail(c, str)
			}
		}
	} else {
		c.HTML(http.StatusOK, "project/add", gin.H{
			"servers": new(models.Server).Find(),
		})
	}
}

func (p *Project) EditDesc(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
		desc := strings.Trim(c.DefaultPostForm("desc", ""), " ")
		if id <= 0 {
			p.Fail(c, "parameter_error")
			return
		}
		if utf8.RuneCountInString(desc) > 500 {
			p.Fail(c, "extra_long_error", "项目描述", "500")
			return
		}

		if err := new(models.Project).Update(id, core.A{"desc": desc}); err == nil {
			p.Success(c, nil)
		} else {
			p.Fail(c, "update_error")
		}
	}
}

//更新项目文件
func (p *Project) EditVersion(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
		useHistoryVersion := strings.Trim(c.DefaultPostForm("useHistoryVersion", ""), " ")
		version := strings.Trim(c.DefaultPostForm("version", ""), " ")
		file, _ := c.FormFile("customFile")
		if id <= 0 {
			p.Fail(c, "parameter_error")
			return
		}
		if utf8.RuneCountInString(version) > 20 {
			p.Fail(c, "extra_long_error", "版本号", "20")
			return
		}
		if file == nil {
			p.Fail(c, "file_upload_error")
			return
		}

		project := models.Project{
			Id: id,
		}
		if ok, err := project.UpdateVersion(useHistoryVersion, version, file); ok {
			p.Success(c, nil)
		} else {
			p.Fail(c, err)
		}
	} else {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if id <= 0 {
			core.Error(c, core.PromptMsg["parameter_error"])
			return
		}
		project := new(models.Project)
		if !project.Get(id) {
			core.Error(c, core.PromptMsg["parameter_error"])
			return
		}
		projectHistory := &models.ProjectHistory{
			ProjectId: id,
		}
		c.HTML(http.StatusOK, "project/edit", gin.H{
			"info":        project,
			"historyList": projectHistory.FindByProjectId(),
		})
	}
}

//更新关联服务器
func (p *Project) EditServers(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
		serverIds, _ := c.GetPostFormArray("serverIds")
		file, _ := c.FormFile("customFile")
		if id <= 0 {
			p.Fail(c, "parameter_error")
			return
		}
		project := new(models.Project)
		if !project.Get(id) {
			p.Fail(c, "parameter_error")
			return
		}
		if ok, err := project.UpdateServers(core.StringArrayToInt(serverIds), file); ok {
			p.Success(c, nil)
		} else {
			p.Fail(c, err)
		}
	} else {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if id <= 0 {
			core.Error(c, core.PromptMsg["parameter_error"])
			return
		}
		project := new(models.Project)
		if !project.Get(id) {
			core.Error(c, core.PromptMsg["parameter_error"])
			return
		}
		servers := make(map[int]map[string]string)
		server := new(models.Server)
		//获取所有服务器
		allServers := server.Find()
		//获取已关联服务器
		projectServers := server.FindByProjectId(id)
		for i := 0; i < len(allServers); i++ {
			servers[allServers[i].Id] = map[string]string{
				"id":       strconv.Itoa(allServers[i].Id),
				"host":     allServers[i].Host,
				"alias":    allServers[i].Alias,
				"selected": "0",
			}
		}
		for i := 0; i < len(projectServers); i++ {
			servers[projectServers[i].Id]["selected"] = "1"
		}
		c.HTML(http.StatusOK, "project/edit_servers", gin.H{
			"info":    project,
			"servers": servers,
		})
	}
}

func (p *Project) Del(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if id <= 0 {
			p.Fail(c, "parameter_error")
			return
		}
		if ok, err := new(models.Project).Del(id); !ok {
			p.Fail(c, err)
			return
		}
		p.Success(c, nil)
	}
}