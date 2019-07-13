package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/elazarl/go-bindata-assetfs"
	"scrapyd-admin/resource/assets"
	"scrapyd-admin/core"
	"scrapyd-admin/controllers/system"
)

func Register(e *gin.Engine) {
	//静态资源
	e.GET("/assets/*path", func(c *gin.Context) {
		http.FileServer(
			&assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, AssetInfo: assets.AssetInfo, Prefix: ""}).ServeHTTP(c.Writer, c.Request)

	})
	e.GET("/", new(Index).Index)
	e.GET("/login", new(Index).Login)
	e.POST("/login", new(Index).DoLogin)
	e.Use(core.CheckLoginStatus(new(core.WebAuth)))
	e.GET("/logout", new(Index).Logout)
	e.GET("/index", new(Index).Index)
	e.GET("/index/main", new(Index).Main)

	//不需要权限验证的
	e.GET("/not_auth/getVersionsByProjectId", new(NotAuth).GetVersionsByProjectId)
	e.GET("/not_auth/getSpidersAndServersByProjectId", new(NotAuth).GetSpidersAndServersByProjectId)

	//系统管理相关
	postAndGet(e, "/system/admin/index", new(system.Admin).Index)
	postAndGet(e, "/system/admin/add", new(system.Admin).Add)
	postAndGet(e, "/system/admin/edit", new(system.Admin).Edit)
	postAndGet(e, "/system/admin/editStatus", new(system.Admin).EditStatus)
	postAndGet(e, "/system/admin/del", new(system.Admin).Del)
	postAndGet(e, "/system/menu/index", new(system.Menu).Index)
	postAndGet(e, "/system/menu/add", new(system.Menu).Add)
	postAndGet(e, "/system/menu/edit", new(system.Menu).Edit)
	postAndGet(e, "/system/menu/editStatus", new(system.Menu).EditStatus)
	e.GET("/system/menu/del", new(system.Menu).Del)

	//项目相关
	postAndGet(e, "/project/project/index", new(Project).Index)
	postAndGet(e, "/project/project/add", new(Project).Add)
	postAndGet(e, "/project/project/editVersion", new(Project).EditVersion)
	postAndGet(e, "/project/project/editServers", new(Project).EditServers)
	e.GET("/project/project/del", new(Project).Del)

	//服务器相关
	postAndGet(e, "/server/server/index", new(Server).Index)
	postAndGet(e, "/server/server/add", new(Server).Add)
	postAndGet(e, "/server/server/edit", new(Server).Edit)
	e.GET("/server/server/del", new(Server).Del)

	//爬虫管理
	postAndGet(e, "/spider/spider/index", new(Spider).Index)

	//任务管理
	postAndGet(e, "/task/task/index", new(Task).Index)
	postAndGet(e, "/task/task/add", new(Task).Add)
	postAndGet(e, "/task/task/cancel", new(Task).Cancel)
	postAndGet(e, "/task/task/cancelMulti", new(Task).CancelMulti)
	postAndGet(e, "/task/task/cancelAll", new(Task).CancelAll)
	postAndGet(e, "/task/task/del", new(Task).Del)
	postAndGet(e, "/task/task/delMulti", new(Task).DelMulti)
	postAndGet(e, "/task/task/delAll", new(Task).DelAll)
	postAndGet(e, "/task/task/schedules", new(Task).Schedules)
	postAndGet(e, "/task/task/addSchedules", new(Task).AddSchedules)
	postAndGet(e, "/task/task/updateSchedulesStatus", new(Task).UpdateSchedulesStatus)
	postAndGet(e, "/task/task/delSchedules", new(Task).DelSchedules)
}

//注册post和get方法
func postAndGet(e *gin.Engine, relativePath string, handlers gin.HandlerFunc) {
	e.GET(relativePath, handlers)
	e.POST(relativePath, handlers)
}
