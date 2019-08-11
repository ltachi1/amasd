package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"scrapyd-admin/core"
	"scrapyd-admin/models"
	"strconv"
	"unicode/utf8"
	"strings"
)

type Menu struct {
	core.BaseController
}

func (a *Menu) Index(c *gin.Context) {
	if core.IsAjax(c) {
		c.JSON(http.StatusOK, gin.H{
			"data": new(models.Menu).TreeMenus(),
		})
	} else {
		c.HTML(http.StatusOK, "menu/index", gin.H{})
	}
}

//添加菜单
func (a *Menu) Add(c *gin.Context) {
	if core.IsAjax(c) {
		var menu models.Menu
		if err := c.ShouldBind(&menu); err == nil {
			if menu.Name == "" {
				a.Fail(c, "system_menu_name_error")
				return
			}
			if utf8.RuneCountInString(menu.Name) > 50 {
				a.Fail(c, "extra_long_error", "菜单名称", "50")
				return
			}
			if menu.App == "" {
				a.Fail(c, "system_menu_app_error")
				return
			}
			if utf8.RuneCountInString(menu.App) > 20 {
				a.Fail(c, "extra_long_error", "模块", "20")
				return
			}
			if menu.Controller == "" {
				a.Fail(c, "system_menu_controller_error")
				return
			}
			if utf8.RuneCountInString(menu.Controller) > 20 {
				a.Fail(c, "extra_long_error", "控制器", "20")
				return
			}
			if menu.Action == "" {
				a.Fail(c, "system_menu_action_error")
				return
			}
			if utf8.RuneCountInString(menu.Action) > 20 {
				a.Fail(c, "extra_long_error", "方法", "20")
				return
			}
			if utf8.RuneCountInString(menu.Parameter) > 50 {
				a.Fail(c, "extra_long_error", "附加参数", "50")
				return
			}
			if utf8.RuneCountInString(menu.Parameter) > 50 {
				a.Fail(c, "extra_long_error", "图标", "50")
				return
			}
			if menu.Insert() {
				a.Success(c, nil)
			} else {
				a.Fail(c, "add_error")
			}
		} else {
			a.Fail(c, "add_error")
		}
	} else {
		parentId, _ := strconv.Atoi(c.DefaultQuery("parent_id", "0"))
		c.HTML(http.StatusOK, "menu/add", gin.H{
			"parentId":  parentId,
			"treeMenus": new(models.Menu).TreeMenus(),
		})
	}
}

//编辑菜单
func (a *Menu) Edit(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
		parentId, _ := strconv.Atoi(c.DefaultPostForm("parent_id", ""))
		name := strings.Trim(c.DefaultPostForm("name", ""), " ")
		app := strings.Trim(c.DefaultPostForm("app", ""), " ")
		controller := strings.Trim(c.DefaultPostForm("controller", ""), " ")
		action := strings.Trim(c.DefaultPostForm("action", ""), " ")
		parameter := strings.Trim(c.DefaultPostForm("parameter", ""), " ")
		icon := strings.Trim(c.DefaultPostForm("icon", ""), " ")
		status, _ := strconv.Atoi(c.DefaultPostForm("status", ""))
		if id <= 0 || parentId <= 0 {
			a.Fail(c, "parameter_error")
			return
		}
		if name == "" {
			a.Fail(c, "system_menu_name_error")
			return
		}
		if utf8.RuneCountInString(name) > 50 {
			a.Fail(c, "extra_long_error", "菜单名称", "50")
			return
		}
		if app == "" {
			a.Fail(c, "system_menu_app_error")
			return
		}
		if utf8.RuneCountInString(app) > 20 {
			a.Fail(c, "extra_long_error", "模块", "20")
			return
		}
		if controller == "" {
			a.Fail(c, "system_menu_controller_error")
			return
		}
		if utf8.RuneCountInString(controller) > 20 {
			a.Fail(c, "extra_long_error", "控制器", "20")
			return
		}
		if action == "" {
			a.Fail(c, "system_menu_action_error")
			return
		}
		if utf8.RuneCountInString(action) > 20 {
			a.Fail(c, "extra_long_error", "方法", "20")
			return
		}
		if utf8.RuneCountInString(parameter) > 50 {
			a.Fail(c, "extra_long_error", "附加参数", "50")
			return
		}
		if utf8.RuneCountInString(icon) > 50 {
			a.Fail(c, "extra_long_error", "图标", "50")
			return
		}
		if !(status == models.MenuStatusEnable || status == models.MenuStatusDisable) {
			a.Fail(c, "system_menu_status_error")
			return
		}

		err := new(models.Menu).Update(id, core.A{
			"parent_id":  parentId,
			"name":       name,
			"app":        app,
			"controller": controller,
			"action":     action,
			"parameter":  parameter,
			"icon":       icon,
			"status":     status,
		})
		if err == nil {
			a.Success(c, nil)
		} else {
			a.Fail(c, "update_error")
		}
	} else {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		if id == 0 {
			core.Error(c, core.PromptMsg["parameter_error"])
			return
		}
		sm := new(models.Menu)
		if ok := sm.Get(id); ok && sm.Id == id {
			c.HTML(http.StatusOK, "menu/edit", gin.H{
				"info":      sm,
				"treeMenus": sm.TreeMenus(),
			})
		} else {
			core.Error(c, core.PromptMsg["system_menu_info_error"])
		}

	}
}

//修改状态
func (a *Menu) EditStatus(c *gin.Context) {
	if core.IsAjax(c) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
		if id <= 0 {
			a.Fail(c, "parameter_error")
			return
		}
		if !(status == models.MenuStatusEnable || status == models.MenuStatusDisable) {
			a.Fail(c, "parameter_error")
			return
		}
		if error := new(models.Menu).Update(id, core.A{"status": status}); error == nil {
			a.Success(c, nil)
		} else {
			a.Fail(c, "update_error")
		}
	}
}

//删除菜单
func (a *Menu) Del(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	if id <= 0 {
		a.Fail(c, "parameter_error")
		return
	}
	sm := new(models.Menu)
	sm.Id = id
	if ok := sm.DeleteById(); ok {
		a.Success(c, nil)
	} else {
		a.Fail(c, "del_error")
	}
}
