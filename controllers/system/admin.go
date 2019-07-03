package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"scrapyd-admin/config"
	"scrapyd-admin/core"
	"scrapyd-admin/models"
	"strconv"
	"github.com/ltachi1/logrus"
)

type Admin struct {
	core.BaseController
}

//用户列表
func (a *Admin) Index(c *gin.Context) {
	if core.IsAjax(c) {
		page, _ := strconv.Atoi(c.DefaultPostForm("pagination[page]", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultPostForm("pagination[perpage]", "10"))
		items, totalCount := new(models.Admin).PageList(page, pageSize)
		c.JSON(http.StatusOK, core.PageResponse(items, page, pageSize, totalCount))
	} else {
		c.HTML(http.StatusOK, "admin/index", gin.H{})
	}
}

//添加用户
func (a *Admin) Add(c *gin.Context) {
	if core.IsAjax(c) {
		email := c.DefaultPostForm("email", "")
		realName := c.DefaultPostForm("real_name", "")
		password := c.DefaultPostForm("password", "")
		confirmPassword := c.DefaultPostForm("confirm_password", "")
		if email == "" {
			a.Fail(c, config.PromptMsg["system_email_error"])
			return
		}
		if !core.IsEmail(email) {
			a.Fail(c, config.PromptMsg["system_email_format_error"])
			return
		}
		if realName == "" {
			a.Fail(c, config.PromptMsg["system_real_name_error"])
			return
		}
		if password == "" {
			a.Fail(c, config.PromptMsg["system_password_error"])
			return
		}
		if password != confirmPassword {
			a.Fail(c, config.PromptMsg["system_password_not_equal_error"])
			return
		}
		admin := &models.Admin{
			Email:    email,
			RealName: realName,
			Password: core.Md5(password),
			Status:   models.AdminStatusNormal,
		}
		if _, error := admin.Create(); error != "" {
			a.Fail(c, config.PromptMsg[error])
			return
		}
		a.Success(c)
	} else {
		c.HTML(http.StatusOK, "admin/add", gin.H{})
	}
}

func (a *Admin) Edit(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
		realName := c.DefaultPostForm("real_name", "")
		password := c.DefaultPostForm("password", "")
		confirmPassword := c.DefaultPostForm("confirm_password", "")
		if realName == "" {
			a.Fail(c, config.PromptMsg["system_real_name_error"])
			return
		}
		if password != confirmPassword {
			a.Fail(c, config.PromptMsg["system_password_not_equal_error"])
			return
		}
		data := core.B{
			"real_name": realName,
		}
		if password != "" {
			data["password"] = core.Md5(password)
		}
		if new(models.Admin).Update(id, data) != nil {
			a.Fail(c, config.PromptMsg["update_error"])
			return
		}
		a.Success(c)
	} else {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if id == 0 {
			core.Error(c, config.PromptMsg["parameter_error"])
			return
		}
		admin := new(models.Admin)
		if !admin.Get(id) {
			core.Error(c, config.PromptMsg["parameter_error"])
			return
		}

		c.HTML(http.StatusOK, "admin/edit", gin.H{
			"info": admin,
		})
	}
}

//修改状态
func (a *Admin) EditStatus(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		status := c.DefaultQuery("status", "")
		if !(status == models.AdminStatusNormal || status == models.AdminStatusDisabled) {
			a.Fail(c, config.PromptMsg["parameter_error"])
			return
		}
		if error := new(models.Admin).Update(id, core.B{"status": status}); error == nil {
			a.Success(c)
		} else {
			a.Fail(c, config.PromptMsg["update_error"])
		}
	}
}

//删除用户
func (a *Admin) Del(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if err := new(models.Admin).Delete(id); err != nil {
			core.WriteLog(config.LogTypeAdmin, logrus.ErrorLevel, logrus.Fields{"id": id}, err)
			a.Fail(c, config.PromptMsg["del_error"])
			return
		}
		a.Success(c)
	}
}
