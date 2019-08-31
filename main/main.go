package main

import (
	"amasd/core"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
	"amasd/resource"
	"amasd/controllers"
	"fmt"
)

func main() {
	e := gin.Default()
	store := sessions.NewCookieStore([]byte(core.AesSalt))
	store.Options(sessions.Options{
		MaxAge: core.SessionExpires,
		Path:   "/",
	})
	e.Use(sessions.Sessions(core.SessionCookieName, store))
	//加载模板文件
	resource.LoadTemplate(e)

	//设置session
	e.Use(core.SetSession())

	//初始化路由
	controllers.Register(e)

	//增加默认异常捕捉程序
	e.Use(core.RecoveryWithWriter())

	e.Run(fmt.Sprintf(":%s", core.HttpPort))

}
