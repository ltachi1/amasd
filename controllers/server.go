package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"scrapyd-admin/core"
	"scrapyd-admin/models"
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
		var server models.Server
		if err := c.ShouldBind(&server); err == nil {
			if utf8.RuneCountInString(server.Alias) > 50 {
				s.Fail(c, "extra_long_error", "别名", "50")
				return
			}
			if utf8.RuneCountInString(server.Host) > 50 {
				s.Fail(c, "extra_long_error", "访问地址", "50")
				return
			}
			if server.Auth == models.ServerAuthClose {
				server.Username, server.Password = "", ""
			}
			if server.Auth == models.ServerAuthOpen && (server.Username == "" || server.Password == "") {
				s.Fail(c, "server_username_error")
				return
			}
			if utf8.RuneCountInString(server.Username) > 20 {
				s.Fail(c, "extra_long_error", "用户名", "20")
				return
			}
			if utf8.RuneCountInString(server.Password) > 20 {
				s.Fail(c, "extra_long_error", "密码", "20")
				return
			}
			if server.Monitor == models.ServerMonitorEnabled && server.MonitorAddress == "" {
				s.Fail(c, "server_monitor_address_error")
				return
			}
			if len(server.MonitorAddress) > 50 {
				s.Fail(c, "extra_long_error", "监控地址", "50")
				return
			}
			if len(server.MonitorUsername) > 20 {
				s.Fail(c, "extra_long_error", "监控地址用户名", "20")
				return
			}
			if len(server.MonitorPassword) > 20 {
				s.Fail(c, "extra_long_error", "监控地址密码", "20")
				return
			}

			if ok, error := server.InsertOne(); !ok {
				s.Fail(c, error)
				return
			}
		} else {
			s.Fail(c, "parameter_error")
			return
		}
		s.Success(c, nil)
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
		monitor := c.DefaultPostForm("monitor", models.ServerMonitorDisabled)
		monitorAddress := c.DefaultPostForm("monitor_address", "")
		monitorUsername := c.DefaultPostForm("monitor_username", "")
		monitorPassword := c.DefaultPostForm("monitor_password", "")
		if id <= 0 {
			s.Fail(c, "parameter_error")
			return
		}
		if utf8.RuneCountInString(alias) > 50 {
			s.Fail(c, "extra_long_error", "别名", "50")
			return
		}
		if uint8(auth) == models.ServerAuthClose {
			username, password = "", ""
		}
		if uint8(auth) == models.ServerAuthOpen && (username == "" || password == "") {
			s.Fail(c, "server_username_error")
			return
		}
		if utf8.RuneCountInString(username) > 20 {
			s.Fail(c, "extra_long_error", "用户名", "20")
			return
		}
		if utf8.RuneCountInString(password) > 20 {
			s.Fail(c, "extra_long_error", "密码", "20")
			return
		}
		if monitor == models.ServerMonitorEnabled && monitorAddress == "" {
			s.Fail(c, "server_monitor_address_error")
			return
		}
		if len(monitorAddress) > 50 {
			s.Fail(c, "extra_long_error", "监控地址", "50")
			return
		}
		if len(monitorUsername) > 20 {
			s.Fail(c, "extra_long_error", "监控地址用户名", "20")
			return
		}
		if len(monitorPassword) > 20 {
			s.Fail(c, "extra_long_error", "监控地址密码", "20")
			return
		}

		data := core.A{
			"alias":            alias,
			"auth":             auth,
			"username":         username,
			"password":         password,
			"monitor":          monitor,
			"monitor_address":  monitorAddress,
			"monitor_username": monitorUsername,
			"monitor_password": monitorPassword,
		}

		if new(models.Server).Update(id, data) != nil {
			s.Fail(c, "update_error")
			return
		}
		s.Success(c, nil)
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