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
	//增加默认异常捕捉程序
	//e.Use(core.RecoveryWithWriter())
	//设置session有效期以及存储路径
	//store, _ := sessions.NewRedisStore(10, "tcp", fmt.Sprintf("%s:%s", "127.0.0.1", config.Conf.Redis.Master.Port), "", []byte("secret"))
	store := sessions.NewCookieStore([]byte("secret"))
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

	e.Run(fmt.Sprintf(":%s", core.HttpPort))

}
