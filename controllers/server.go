package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"amasd/core"
	"amasd/models"
	"strconv"
	"unicode/utf8"
	"time"
	"strings"
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
		if server, ok := s.verifyData("add", c); ok {
			server.Status = models.ServerStatusNormal
			if ok, error := server.InsertOne(); !ok {
				s.Fail(c, error)
				return
			}
			s.Success(c, nil)
		}
	} else {
		c.HTML(http.StatusOK, "server/add", gin.H{
			"projects": new(models.Project).Find(),
		})
	}
}

func (s *Server) Edit(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
		if id <= 0 {
			s.Fail(c, "parameter_error")
			return
		}
		if server, ok := s.verifyData("edit", c); ok {
			data := core.A{
				"alias":            server.Alias,
				"auth":             server.Auth,
				"username":         server.Username,
				"password":         server.Password,
				"monitor":          server.Monitor,
				"monitor_address":  server.MonitorAddress,
				"monitor_username": server.MonitorUsername,
				"monitor_password": server.MonitorPassword,
			}
			if new(models.Server).Update(id, data) != nil {
				s.Fail(c, "update_error")
				return
			}
			s.Success(c, nil)
		}
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

func (s *Server) verifyData(t string, c *gin.Context) (models.Server, bool) {
	alias := strings.Trim(c.DefaultPostForm("alias", ""), " ")
	host := strings.Trim(c.DefaultPostForm("host", ""), " ")
	authInt, _ := strconv.Atoi(c.DefaultPostForm("auth", strconv.Itoa(int(models.ServerAuthClose))))
	auth := uint8(authInt)
	username := strings.Trim(c.DefaultPostForm("username", ""), " ")
	password := strings.Trim(c.DefaultPostForm("password", ""), " ")
	monitor := strings.Trim(c.DefaultPostForm("monitor", models.ServerMonitorDisabled), " ")
	monitorAddress := strings.Trim(c.DefaultPostForm("monitorAddress", ""), " ")
	monitorUsername := strings.Trim(c.DefaultPostForm("monitorUsername", ""), " ")
	monitorPassword := strings.Trim(c.DefaultPostForm("monitorPassword", ""), " ")

	server := models.Server{}
	if utf8.RuneCountInString(alias) > 50 {
		s.Fail(c, "extra_long_error", "别名", "50")
		return server, false
	}
	if t == "add" {
		if len(host) == 0 {
			s.Fail(c, "parameter_required")
			return server, false
		}
		if utf8.RuneCountInString(host) > 50 {
			s.Fail(c, "extra_long_error", "访问地址", "50")
			return server, false
		}
	}
	if auth != models.ServerAuthClose && auth != models.ServerAuthOpen {
		s.Fail(c, "parameter_error")
		return server, false
	}
	if auth == models.ServerAuthClose {
		username, password = "", ""
	}
	if auth == models.ServerAuthOpen && (username == "" || password == "") {
		s.Fail(c, "server_username_error")
		return server, false
	}
	if utf8.RuneCountInString(username) > 20 {
		s.Fail(c, "extra_long_error", "用户名", "20")
		return server, false
	}
	if utf8.RuneCountInString(password) > 20 {
		s.Fail(c, "extra_long_error", "密码", "20")
		return server, false
	}
	if monitor != models.ServerMonitorEnabled && monitor != models.ServerMonitorDisabled {
		s.Fail(c, "parameter_error")
		return server, false
	}
	if monitor == models.ServerMonitorEnabled && monitorAddress == "" {
		s.Fail(c, "server_monitor_address_error")
		return server, false
	}
	if len(monitorAddress) > 50 {
		s.Fail(c, "extra_long_error", "监控地址", "50")
		return server, false
	}
	if len(monitorUsername) > 20 {
		s.Fail(c, "extra_long_error", "监控访问用户名", "20")
		return server, false
	}
	if len(monitorPassword) > 20 {
		s.Fail(c, "extra_long_error", "监控访问密码", "20")
		return server, false
	}
	if len(monitorUsername) > 0 && len(monitorPassword) == 0 {
		s.Fail(c, "server_monitor_password_error")
		return server, false
	}
	if t == "add" {
		server.Host = host
	}
	server.Alias = alias
	server.Auth = auth
	server.Username = username
	server.Password = password
	server.Monitor = monitor
	server.MonitorAddress = monitorAddress
	server.MonitorUsername = monitorUsername
	server.MonitorPassword = monitorPassword
	return server, true
}

func (s *Server) Del(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if ok, err := new(models.Server).Del(id); !ok {
			s.Fail(c, err)
			return
		}
		s.Success(c, nil)
	}
}

func (s *Server) Monitor(c *gin.Context) {
	if core.IsAjax(c) {
		ids := c.DefaultQuery("ids", "")
		if len(ids) == 0 {
			s.Fail(c, "parameter_error")
			return
		}
		items, _ := new(models.ServerMonitor).OverviewByIds(strings.Split(ids, ","))
		s.Success(c, core.A{
			"items": items,
		})
	} else {
		//获取所有正在监控的服务器
		server := models.Server{
			Monitor: models.ServerMonitorEnabled,
		}

		c.HTML(http.StatusOK, "server/monitor", gin.H{
			"servers": server.Find(),
		})
	}

}

func (s *Server) MonitorDetail(c *gin.Context) {
	if core.IsAjax(c) {
		serverId, _ := strconv.Atoi(c.DefaultQuery("server_id", "0"))
		lastTime, _ := strconv.Atoi(c.DefaultQuery("last_time", "0"))
		if serverId <= 0 {
			s.Fail(c, "parameter_error")
			return
		}
		items := new(models.ServerMonitor).FindByLastTime(serverId, lastTime)
		nextTime := 0
		if len(items) > 0 {
			//nextTime = int(items[len(items) - 1 ]["time"].(core.Timestamp))
			nextTime = int(time.Now().Unix())
		} else {
			nextTime = lastTime
		}
		s.Success(c, core.A{
			"items":     items,
			"next_time": nextTime,
		})
	} else {
		serverId, _ := strconv.Atoi(c.DefaultQuery("server_id", "0"))
		if serverId <= 0 {
			core.Error(c, core.PromptMsg["parameter_error"])
			return
		}
		server := models.Server{}
		if !server.Get(serverId) {
			core.Error(c, core.PromptMsg["parameter_error"])
			return
		}
		c.HTML(http.StatusOK, "server/monitor_detail", gin.H{
			"info": server,
		})
	}

}
