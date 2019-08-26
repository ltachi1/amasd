package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"amasd/core"
	"amasd/models"
	"strconv"
	"github.com/ltachi1/logrus"
	"unicode/utf8"
	"strings"
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
		username := strings.Trim(c.DefaultPostForm("username", ""), " ")
		email := strings.Trim(c.DefaultPostForm("email", ""), " ")
		displayName := strings.Trim(c.DefaultPostForm("display_name", ""), " ")
		password := strings.Trim(c.DefaultPostForm("password", ""), " ")
		confirmPassword := strings.Trim(c.DefaultPostForm("confirm_password", ""), " ")
		if username == "" {
			a.Fail(c, "system_username_error")
			return
		}
		if utf8.RuneCountInString(username) > 50 {
			a.Fail(c, "extra_long_error", "用户名", "50")
			return
		}
		if displayName == "" {
			a.Fail(c, "system_display_name_error")
			return
		}
		if utf8.RuneCountInString(displayName) > 20 {
			a.Fail(c, "extra_long_error", "昵称", "20")
			return
		}
		if utf8.RuneCountInString(email) > 50 {
			a.Fail(c, "extra_long_error", "邮箱", "50")
			return
		}
		if email != "" && !core.IsEmail(email) {
			a.Fail(c, "system_email_format_error")
			return
		}
		if password == "" {
			a.Fail(c, "system_password_error")
			return
		}
		if utf8.RuneCountInString(password) > 20 {
			a.Fail(c, "extra_long_error", "密码", "20")
			return
		}
		if password != confirmPassword {
			a.Fail(c, "system_password_not_equal_error")
			return
		}


		admin := &models.Admin{
			Username:    username,
			Email:       email,
			DisplayName: displayName,
			Password:    core.Md5(password),
			Status:      models.AdminStatusNormal,
		}
		if _, error := admin.Create(); error != "" {
			a.Fail(c, error)
			return
		}
		a.Success(c, nil)
	} else {
		c.HTML(http.StatusOK, "admin/add", gin.H{})
	}
}

func (a *Admin) Edit(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
		displayName := strings.Trim(c.DefaultPostForm("display_name", ""), " ")
		email := strings.Trim(c.DefaultPostForm("email", ""), " ")
		password := strings.Trim(c.DefaultPostForm("password", ""), " ")
		confirmPassword := strings.Trim(c.DefaultPostForm("confirm_password", ""), " ")
		if id <= 0 {
			a.Fail(c, "parameter_error")
			return
		}
		if displayName == "" {
			a.Fail(c, "system_display_name_error")
			return
		}
		if utf8.RuneCountInString(displayName) > 20 {
			a.Fail(c, "extra_long_error", "昵称", "20")
			return
		}
		if utf8.RuneCountInString(email) > 50 {
			a.Fail(c, "extra_long_error", "邮箱", "50")
			return
		}
		if email != "" && !core.IsEmail(email) {
			a.Fail(c, "system_email_format_error")
			return
		}
		if utf8.RuneCountInString(password) > 20 {
			a.Fail(c, "extra_long_error", "密码", "20")
			return
		}
		if password != confirmPassword {
			a.Fail(c, "system_password_not_equal_error")
			return
		}
		data := core.B{
			"display_name": displayName,
			"email":        email,
		}
		if password != "" {
			data["password"] = core.Md5(password)
		}
		if new(models.Admin).Update(id, data) != nil {
			a.Fail(c, "update_error")
			return
		}
		a.Success(c, nil)
	} else {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if id == 0 {
			core.Error(c, core.PromptMsg["parameter_error"])
			return
		}
		admin := new(models.Admin)
		if !admin.Get(id) {
			core.Error(c, core.PromptMsg["parameter_error"])
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
			a.Fail(c, "parameter_error")
			return
		}
		if error := new(models.Admin).Update(id, core.B{"status": status}); error == nil {
			a.Success(c, nil)
		} else {
			a.Fail(c, "update_error")
		}
	}
}

//删除用户
func (a *Admin) Del(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if id <= 0 {
			a.Fail(c, "parameter_error")
			return
		}
		if id == 1 {
			a.Fail(c, "system_admin_not_del_error")
			return
		}
		if err := new(models.Admin).Delete(id); err != nil {
			core.WriteLog(core.LogTypeAdmin, logrus.ErrorLevel, logrus.Fields{"id": id}, err)
			a.Fail(c, "del_error")
			return
		}
		a.Success(c, nil)
	}
}
