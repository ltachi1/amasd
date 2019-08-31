package core

const (
	//长度16、24、32必须是这三种中的一种，否则加密解密是报错
	AesSalt           = "amasd-salt123456"
	SessionCookieName = "AMASD_SESSION_ID"
	SessionExpires    = 1800
	EnvDev            = "dev"
	EnvTesting        = "testing"
	EnvProduction     = "production"
	ErrorDir          = "error"
	LogTypeDb         = "db"
	LogTypeProject    = "project"
	LogTypeSpider     = "spiders"
	LogTypeTask       = "task"
	LogTypeServer     = "server"
	LogTypeAdmin      = "admin"
	LogTypeScrapyd    = "scrapyd"
	LogTypeInit       = "init"
	LogTypeSystem     = "system"
)

var (
	//返回码
	PromptMsg = map[string]string{
		"success":  "成功",   //通用成功
		"fail":     "未知错误", //通用未知错误码，如系统错误等
		"no_login": "重新登录", //需要登录统一码

		"parameter_error":    "参数错误",
		"parameter_required": "请输入必填项",
		"add_error":          "添加失败",
		"update_error":       "更新失败",
		"del_error":          "删除失败",
		"extra_long_error":   "%s最多可输入%s字符",
		//token相关
		"token_valid":   "Token 无效",
		"token_expired": "Token 过期",
		//登录相关
		"login_password_error": "用户名或密码错误",
		"login_user_disable":   "此用户不允许登录",
		//服务器相关
		"scrapyd_server_error":              "请检查爬虫服务器是否可用",
		"file_upload_error":                 "请确定上传正确的文件",
		"host_error":                        "请输入正确的服务器访问地址",
		"host_repeat_error":                 "此服务器已存在不能重复添加",
		"server_info_error":                 "服务器信息获取错误",
		"server_username_error":             "请输入服务器用户名和密码",
		"server_cutback_task_running_error": "所减少的服务器有正在运行的爬虫或者定时任务，请先删除或者停止",
		"server_del_task_running_error":     "此服务器有正在运行的爬虫或者定时任务，请先删除或者停止",
		"server_monitor_address_error":      "请输入监控地址",
		"server_monitor_password_error":     "请输入监控访问密码",
		//项目相关
		"project_name_repeat":               "项目名称重复，请重新输入",
		"project_version_repeat":            "项目版本号重复，请重新输入",
		"project_info_error":                "项目信息获取错误",
		"project_server_error":              "请选择已有服务器",
		"project_server_relation_error":     "部分服务器关联失败,请重试，包括: %s",
		"project_update_version_error":      "部分服务器版本更新失败，请重新操作",
		"project_no_server":                 "当前项目无可用服务器,请先关联服务器",
		"project_server_relation_all_error": "服务器关联失败，请重新操作",
		"task_running_error":                "当前项目有正在运行的爬虫不允许更新，请手动关闭或等待执行完",
		"project_spider_number_error":       "获取不到项目下爬虫信息，请重新尝试",
		"project_spider_update_error":       "爬虫信息更新失败，请重新尝试",
		"project_del_error":                 "项目删除失败，请重试(某些服务器上的项目文件可能已经删除)",
		"project_task_running_error":        "当前项目有正在运行的爬虫或者定时任务，请先删除或者停止",

		//任务相关
		"task_add_error":    "部分任务添加失败，包括: %s",
		"task_update_error": "部分任务停止失败，包括: %s",

		//系统相关
		"system_username_error":               "请输入用户名",
		"system_username_repeat_error":        "用户名重复",
		"system_email_format_error":           "请输入正确的邮箱",
		"system_email_repeat_error":           "邮箱重复",
		"system_display_name_error":           "请输入昵称",
		"system_password_error":               "请输入密码",
		"system_password_not_equal_error":     "密码输入不一致",
		"system_admin_not_del_error":          "此用户不允许删除",
		"system_menu_error":                   "密码输入不一致",
		"system_menu_name_error":              "请输入菜单名称",
		"system_menu_app_error":               "输入模块名称",
		"system_menu_controller_error":        "输入控制器名称",
		"system_menu_action_error":            "输入方法名称",
		"system_menu_status_error":            "状态类型错误",
		"system_menu_info_error":              "菜单信息获取失败",
		"system_notice_scrapyd_error":         "请输入scrapyd服务异常通知标题和内容",
		"system_notice_task_finished_error":   "请输入任务运行结束通知标题和内容",
		"system_notice_task_error":            "请输入任务运行异常通知标题和内容",
		"system_notice_email_smtp_error":      "请输入正确的发件服务器",
		"system_notice_email_error":           "请输入邮箱设置相关信息",
		"system_notice_email_port_error":      "请输入正确的发件服务器端口",
		"system_notice_email_address_error":   "请输入正确的发件人邮箱地址",
		"system_notice_email_addressee_error": "请输入正确的收件人邮箱地址",
		"system_notice_dingtalk_error":        "请输入钉钉群机器人webhook地址",
		"system_notice_work_weixin_error":     "请输入企业微信群机器人webhook地址",
		"system_notice_webhook_error":         "请输入正确的的webhook地址",
	}
	//运行文件所在路径
	RootPath = ""
	//当前运行环境
	Env = ""
	//监听端口
	HttpPort       = "8000"
	NoticeSettings = B{}
)
