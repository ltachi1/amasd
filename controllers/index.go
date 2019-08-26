package controllers

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"amasd/core"
	"amasd/models"
)

type Index struct {
	core.BaseController
}

func (i *Index) Index(c *gin.Context) {
	sm := new(models.Menu)
	passport := core.GetPassportInstance()
	c.HTML(http.StatusOK, "index/index", gin.H{
		"menuStr":     template.HTML(sm.GetMenuStr()),
		"displayName": passport.DisplayName,
		"id":          passport.Id,
	})
}

func (i *Index) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "index/login", gin.H{})
}

//登录
func (i *Index) DoLogin(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBind(&admin); err == nil {
		ok, code := admin.Login()
		if !ok {
			i.Fail(c, code)
			return
		}
		i.Success(c, nil)
	} else {
		i.Fail(c)
	}
}

//退出登录
func (i *Index) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func (i *Index) Error(c *gin.Context) {

}
