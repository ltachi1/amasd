// session相关
package core

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var session sessions.Session

//设置session
func SetSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session = sessions.Default(c)
	}
}

//获取session
func GetSession() sessions.Session {
	return session
}
