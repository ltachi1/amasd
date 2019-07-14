package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"scrapyd-admin/core"
	"scrapyd-admin/models"
	"strconv"
)

type Server struct {
	core.BaseController
}

func (s *Server) Index(c *gin.Context) {
	if core.IsAjax(c) {
		projectId, _ := strconv.Atoi(c.DefaultPostForm("project_id", "0"))
		page, _ := strconv.Atoi(c.DefaultPostForm("pagination[page]", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultPostForm("pagination[perpage]", "10"))
		servers, totalCount := new(models.Server).PageList(projectId, page, pageSize)
		c.JSON(http.StatusOK, gin.H{
			"data": servers,
			"meta": gin.H{
				"page":    page,
				"total":   totalCount,
				"pages":   core.CalculationPages(totalCount, pageSize),
				"perpage": pageSize,
			},
		})
	} else {
		c.HTML(http.StatusOK, "server/index", gin.H{
			"projects": new(models.Project).Find(),
		})
	}

}

func (s *Server) Add(c *gin.Context) {
	if core.IsAjax(c) {
		var server models.Server
		if err := c.ShouldBind(&server); err == nil {
			if server.Auth == models.ServerAuthClose {
				server.Username, server.Password = "", ""
			}
			if server.Auth == models.ServerAuthOpen && (server.Username == "" || server.Password == "") {
				s.Fail(c, core.PromptMsg["server_username_error"])
				return
			}
			if ok, error := server.InsertOne(); !ok {
				s.Fail(c, core.PromptMsg[error])
				return
			}
		} else {
			s.Fail(c, core.PromptMsg["parameter_error"])
			return
		}
		s.Success(c)
	} else {
		c.HTML(http.StatusOK, "server/add", gin.H{
			"projects": new(models.Project).Find(),
		})
	}
}

func (s *Server) Edit(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
		alias := c.DefaultPostForm("alias", "")
		auth, _ := strconv.Atoi(c.DefaultPostForm("auth", "1"))
		username := c.DefaultPostForm("username", "")
		password := c.DefaultPostForm("password", "")

		if uint8(auth) == models.ServerAuthClose {
			username, password = "", ""
		}

		if uint8(auth) == models.ServerAuthOpen && (username == "" || password == "") {
			s.Fail(c, core.PromptMsg["server_username_error"])
			return
		}

		data := core.B{
			"alias":    alias,
			"auth":     strconv.Itoa(auth),
			"username": username,
			"password": password,
		}

		if new(models.Server).Update(id, data) != nil {
			s.Fail(c, core.PromptMsg["update_error"])
			return
		}
		s.Success(c)
	} else {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		server := new(models.Server)
		if !server.Get(id) {
			core.Error(c, core.PromptMsg["parameter_error"])
			return
		}
		c.HTML(http.StatusOK, "server/edit", gin.H{
			"info": server,
		})
	}
}

func (s *Server) Del(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if ok, err := new(models.Server).Del(id); !ok {
			s.Fail(c, core.PromptMsg[err])
			return
		}
		s.Success(c)
	}
}
