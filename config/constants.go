package config

import (
	"github.com/gin-gonic/gin"
)

const (
	//长度16、24、32必须是这三种中的一种，否则加密解密是报错
	AesSalt        = "scrapyd-admin-salt123456"
	EnvDev         = "dev"
	EnvTesting     = "testing"
	EnvProduction  = "production"
	ErrorDir       = "error"
	LogTypeDb      = "db"
	LogTypeProject = "project"
	LogTypeSpider  = "spiders"
	LogTypeTask    = "task"
	LogTypeServer  = "server"
	LogTypeAdmin   = "admin"
	LogTypeScrapyd = "scrapyd"
	LogTypeInit    = "init"
	LogTypeSystem  = "system"
)

var (
	//返回码
	PromptMsg = map[string]map[string]interface{}{
		"success":  gin.H{"code": 0, "msg": "成功"},   //通用成功
		"fail":     gin.H{"code": 1, "msg": "未知错误"}, //通用未知错误码，如系统错误等
		"no_login": gin.H{"code": 2, "msg": "重新登录"}, //需要登录统一码

		"parameter_error": gin.H{"msg": "参数错误"},
		"add_error":       gin.H{"msg": "添加失败"},
		"update_error":    gin.H{"msg": "更新失败"},
		"del_error":       gin.H{"msg": "删除失败"},
		//token相关
		"token_valid":   gin.H{"msg": "Token 无效"},
		"token_expired": gin.H{"msg": "Token 过期"},
		//登录相关
		"login_password_error": gin.H{"msg": "用户名或密码错误"},
		"login_user_disable":   gin.H{"msg": "此用户不允许登录"},
		//服务器相关
		"scrapyd_server_error":  gin.H{"msg": "请检查爬虫服务器是否可用"},
		"file_upload_error":     gin.H{"msg": "请确定上传正确的文件"},
		"host_error":            gin.H{"msg": "请输入正确的服务器访问地址"},
		"host_repeat_error":     gin.H{"msg": "此服务器已存在不能重复添加"},
		"server_info_error":     gin.H{"msg": "服务器信息获取错误"},
		"server_username_error": gin.H{"msg": "请输入服务器用户名和密码"},
		//项目相关
		"project_name_repeat":               gin.H{"msg": "项目名称重复，请重新输入"},
		"project_version_repeat":            gin.H{"msg": "项目版本号重复，请重新输入"},
		"project_info_error":                gin.H{"msg": "项目信息获取错误"},
		"project_server_error":              gin.H{"msg": "请选择已有服务器"},
		"project_server_relation_error":     gin.H{"code": 4001, "msg": "部分服务器关联失败,请重试，包括: "},
		"project_update_version_error":      gin.H{"msg": "部分服务器版本更新失败，请重新操作"},
		"project_server_relation_all_error": gin.H{"msg": "服务器关联失败，请重新操作"},
		"task_running_error":                gin.H{"msg": "当前项目有正在运行的爬虫不允许更新，请手动关闭或等待执行完"},
		"project_spider_number_error":       gin.H{"msg": "获取不到项目下爬虫信息，请重新尝试"},
		"project_spider_update_error":       gin.H{"msg": "爬虫信息更新失败，请重新尝试"},

		//任务相关
		"task_add_error":    gin.H{"code": 5001, "msg": "部分任务添加失败，包括: "},
		"task_update_error": gin.H{"code": 5001, "msg": "部分任务停止失败，包括: "},

		//系统相关
		"system_email_error":              gin.H{"msg": "请输入邮箱"},
		"system_email_format_error":       gin.H{"msg": "请输入正确的邮箱"},
		"system_email_repeat_error":       gin.H{"msg": "邮箱重复"},
		"system_real_name_error":          gin.H{"msg": "请输入姓名"},
		"system_password_error":           gin.H{"msg": "请输入密码"},
		"system_password_not_equal_error": gin.H{"msg": "密码输入不一致"},
		"system_menu_error":               gin.H{"msg": "密码输入不一致"},
		"system_menu_name_error":          gin.H{"msg": "请输入菜单名称"},
		"system_menu_app_error":           gin.H{"msg": "输入模块名称"},
		"system_menu_controller_error":    gin.H{"msg": "输入控制器名称"},
		"system_menu_action_error":        gin.H{"msg": "输入方法名称"},
		"system_menu_status_error":        gin.H{"msg": "状态类型错误"},
		"system_menu_info_error":          gin.H{"msg": "菜单信息获取失败"},
	}
	//运行文件所在路径
	RootPath = ""
)
