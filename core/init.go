// 核心包初始化
package core

import "github.com/gin-gonic/gin"

func init() {
	//初始化日志
	InitLog()
	//初始化数据库
	InitDb()
	//初始化定时任务
	InitCron()
}

//检查登录状态
func CheckLoginStatus(a Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		a.Check(c)
	}
}
