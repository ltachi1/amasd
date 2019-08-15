package models

import (
	"scrapyd-admin/core"
	"errors"
	"github.com/go-xorm/xorm"
	"github.com/ltachi1/logrus"
)

func init() {
	error := InitTables()
	if error != nil {
		core.WriteLog(core.LogTypeDb, logrus.PanicLevel, nil, error)
	}

	//初始化定时任务相关
	InitTask()

	//初始化通知配置
	core.NoticeSettings = new(NoticeSetting).Find()
}

func InitTask() {
	//定时检测服务器状态
	core.Cron.AddFunc("*/30 * * * *", func() {
		new(Server).DetectionStatus()
	}, "DetectionServerStatus")
	//定时检测任务状态
	core.Cron.AddFunc("*/10 * * * *", func() {
		new(Task).DetectionStatus()
	}, "DetectionTaskStatus")

	//初始化已有的计划任务
	new(SchedulesTask).InitSchedulesToCron()
}

func InitTables() error {
	upgradeFuncs := []func(*xorm.Session) error{
		upgrade100,
		upgrade200,
		upgrade210,
	}
	session := core.Db.NewSession()
	defer session.Close()
	session.Begin()
	for i := 0; i < len(upgradeFuncs); i++ {
		if err := upgradeFuncs[i](session); err != nil {
			session.Rollback()
			return err
		}
	}
	return session.Commit()
}

//1.0.0数据库
func upgrade100(session *xorm.Session) error {
	//判断表是否存在
	exist, err := session.IsTableExist("admin")
	if err != nil {
		errors.New("数据库表创建失败")
	}
	if exist {
		return nil
	}
	_, err = session.Exec(`CREATE TABLE access (id INTEGER PRIMARY KEY AUTOINCREMENT, role_id INTEGER NOT NULL DEFAULT (0), app VARCHAR (50) NOT NULL DEFAULT "", controller VARCHAR (50) NOT NULL DEFAULT "", "action" VARCHAR (50) NOT NULL DEFAULT "", status TINYINT (1) NOT NULL DEFAULT (1));
CREATE TABLE admin (id INTEGER PRIMARY KEY AUTOINCREMENT, username VARCHAR (50) NOT NULL DEFAULT "", password VARCHAR (100) NOT NULL DEFAULT "", email VARCHAR (100) NOT NULL DEFAULT "", display_name VARCHAR (50) NOT NULL DEFAULT "", status VARCHAR (10) DEFAULT "" NOT NULL, create_time INT (10) DEFAULT (0) NOT NULL);
INSERT INTO "admin" VALUES (1, 'admin', '21232f297a57a5a743894a0e4a801fc3', 'admin@163.com', '管理员', 'enabled', 1562306601);
CREATE TABLE admin_role (id INTEGER PRIMARY KEY AUTOINCREMENT, admin_id INTEGER NOT NULL DEFAULT (0), role_id INTEGER NOT NULL DEFAULT (0));
INSERT INTO "admin_role" VALUES (1, 1, 1);
CREATE TABLE menu (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR (50) NOT NULL DEFAULT "", parent_id INTEGER NOT NULL DEFAULT (0), app VARCHAR (50) NOT NULL DEFAULT "", controller VARCHAR (50) NOT NULL DEFAULT "", "action" VARCHAR (50) NOT NULL DEFAULT "", parameter VARCHAR (200) NOT NULL DEFAULT "", icon VARCHAR (50) NOT NULL DEFAULT "", status TINYINT (1) NOT NULL DEFAULT (1), listing_order SMALLINT (5) NOT NULL DEFAULT (0));
INSERT INTO "menu" VALUES (1, '系统管理', 0, 'system', 'system', 'system', '', '', 1, 3);
INSERT INTO "menu" VALUES (2, '用户管理', 1, 'system', 'admin', 'index', '', 'fa fa-user', 1, 0);
INSERT INTO "menu" VALUES (3, '菜单管理', 1, 'system', 'menu', 'index', '', 'fa fa-list-ol', 2, 0);
INSERT INTO "menu" VALUES (4, '添加菜单', 3, 'system', 'menu', 'add', '', '', 2, 0);
INSERT INTO "menu" VALUES (5, '编辑菜单', 3, 'system', 'menu', 'add', '', '', 2, 0);
INSERT INTO "menu" VALUES (6, '项目', 0, 'project', 'project', 'project', '', '', 1, 0);
INSERT INTO "menu" VALUES (7, '项目管理', 6, 'project', 'project', 'index', '', 'fa fa-box', 1, 0);
INSERT INTO "menu" VALUES (8, ' 添加项目', 6, 'project', 'project', 'add', '', 'fa fa-box-open', 1, 0);
INSERT INTO "menu" VALUES (9, '爬虫', 0, 'spider', 'spider', 'spider', '', '', 1, 0);
INSERT INTO "menu" VALUES (10, '爬虫管理', 9, 'spider', 'spider', 'index', '', 'fa fa-bug', 1, 0);
INSERT INTO "menu" VALUES (12, '任务', 0, 'jobs', 'jobs', 'jobs', '', '', 1, 0);
INSERT INTO "menu" VALUES (13, '任务列表', 12, 'task', 'task', 'index', '', 'fa fa-list', 1, 0);
INSERT INTO "menu" VALUES (14, '计划任务', 12, 'task', 'task', 'schedules', '', 'fa fa-list-alt', 1, 0);
INSERT INTO "menu" VALUES (15, '服务器', 0, 'server', 'server', 'server', '', '', 1, 0);
INSERT INTO "menu" VALUES (16, '服务器管理', 15, 'server', 'server', 'index', '', 'fa fa-server', 1, 0);
INSERT INTO "menu" VALUES (17, '添加服务器', 15, 'server', 'server', 'add', '', 'fa fa-folder-plus', 1, 0);
CREATE TABLE project (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR (50) NOT NULL DEFAULT "", last_version VARCHAR (50) NOT NULL DEFAULT "", create_time INT (10) NOT NULL DEFAULT (0), update_time INT (10) NOT NULL DEFAULT (0));
CREATE TABLE project_history (id INTEGER PRIMARY KEY AUTOINCREMENT, project_id INTEGER NOT NULL DEFAULT (0), version VARCHAR (50) NOT NULL DEFAULT "", create_time INT (10) NOT NULL DEFAULT (0));
CREATE TABLE project_server (id INTEGER PRIMARY KEY AUTOINCREMENT, project_id INTEGER NOT NULL DEFAULT (0), server_id INTEGER NOT NULL DEFAULT (0), create_time INT (10) NOT NULL DEFAULT (0), update_time INT (10) NOT NULL DEFAULT (0));
CREATE TABLE role (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR (50) NOT NULL DEFAULT "", parent_id INTEGER NOT NULL DEFAULT (0), status TINYINT (1) NOT NULL DEFAULT (1), remark VARCHAR (100) NOT NULL DEFAULT "", create_time INT (10) NOT NULL DEFAULT (0));
INSERT INTO "role" VALUES (1, '超级管理员', 0, 1, '', 0);
CREATE TABLE schedules_task (id INTEGER PRIMARY KEY AUTOINCREMENT, project_id INTEGER NOT NULL DEFAULT (0), project_name VARCHAR (50) NOT NULL DEFAULT "", version VARCHAR (50) NOT NULL DEFAULT "", server_id INTEGER NOT NULL DEFAULT (0), host VARCHAR (50) NOT NULL DEFAULT "", spider_id INTEGER NOT NULL DEFAULT (0), spider_name VARCHAR (100) NOT NULL DEFAULT "", cron VARCHAR (50) NOT NULL DEFAULT "", status VARCHAR (20) NOT NULL DEFAULT "");
CREATE TABLE server (id INTEGER PRIMARY KEY AUTOINCREMENT, alias VARCHAR (50) NOT NULL DEFAULT "", host VARCHAR (50) NOT NULL DEFAULT "", username VARCHAR (50) NOT NULL DEFAULT "", password VARCHAR NOT NULL DEFAULT "", auth TINYINT (1) NOT NULL DEFAULT (1), status TINYINT (1) NOT NULL DEFAULT (1), create_time INT (10) NOT NULL DEFAULT (0), update_time INT (10) NOT NULL DEFAULT (0));
CREATE TABLE spider (id INTEGER PRIMARY KEY AUTOINCREMENT, project_id INTEGER NOT NULL DEFAULT (0), name VARCHAR (100) NOT NULL DEFAULT "", version VARCHAR (50) NOT NULL DEFAULT "", status TINYINT (1) NOT NULL DEFAULT (1), create_time INT (10) NOT NULL DEFAULT (0));
CREATE TABLE task (id INTEGER PRIMARY KEY AUTOINCREMENT, type TINYINT (1) NOT NULL DEFAULT (1), project_id INTEGER NOT NULL DEFAULT (0), project_name VARCHAR (50) NOT NULL DEFAULT "", version VARCHAR (50) NOT NULL DEFAULT "", server_id INTEGER NOT NULL DEFAULT (0), host VARCHAR (50) NOT NULL DEFAULT "", spider_id INTEGER NOT NULL DEFAULT (0), spider_name VARCHAR (100) NOT NULL DEFAULT "", job_id VARCHAR (50) NOT NULL DEFAULT "", start_time INT (10) NOT NULL DEFAULT (0), end_time INT (10) NOT NULL DEFAULT (0), status VARCHAR (20) NOT NULL DEFAULT "");
`)
	return err
}

//升级2.0.0数据库
func upgrade200(session *xorm.Session) error {
	exist, err := session.IsTableExist("notice_setting")
	if err != nil {
		errors.New("数据库表创建失败")
	}
	if exist {
		return nil
	}
	_, err = session.Exec(`CREATE TABLE "notice_setting" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,"name" TEXT(50) NOT NULL DEFAULT '',"value" TEXT NOT NULL DEFAULT '',"desc" TEXT NOT NULL DEFAULT '');
INSERT INTO "notice_setting" VALUES (1, 'scrapyd_service', 'disabled', 'scrapyd服务监控状态');
INSERT INTO "notice_setting" VALUES (2, 'task_finished', 'disabled', '是否监听任务完成');
INSERT INTO "notice_setting" VALUES (3, 'task_error', 'disabled', '是否监听任务异常');
INSERT INTO "notice_setting" VALUES (4, 'scrapyd_service_title', '', 'scrpayd服务异常通知标题');
INSERT INTO "notice_setting" VALUES (5, 'scrapyd_service_content', '', 'scrpayd服务异常通知内容');
INSERT INTO "notice_setting" VALUES (6, 'task_finished_title', '', '任务完成通知标题');
INSERT INTO "notice_setting" VALUES (7, 'task_finished_content', '', '任务完成通知内容');
INSERT INTO "notice_setting" VALUES (8, 'task_error_title', '', '任务异常通知标题');
INSERT INTO "notice_setting" VALUES (9, 'task_error_content', '', '任务异常通知内容');
INSERT INTO "notice_setting" VALUES (10, 'email', 'disabled', '是否开启邮件通知');
INSERT INTO "notice_setting" VALUES (11, 'email_smtp', '', '发件人smtp地址');
INSERT INTO "notice_setting" VALUES (12, 'email_sender_address', '', '发件人邮箱地址');
INSERT INTO "notice_setting" VALUES (13, 'email_sender_password', '', '发件人邮箱密码');
INSERT INTO "notice_setting" VALUES (14, 'email_addressee', '', '收件人邮箱地址');
INSERT INTO "notice_setting" VALUES (15, 'email_sender', '', '发件人别名');
INSERT INTO "notice_setting" VALUES (16, 'email_smtp_port', '', '发件箱服务器端口');
INSERT INTO "menu" VALUES (18, '通知设置', 1, 'notice', 'notice', 'setting', '', 'fa fa-bell', 1, 0);
`)
	return err
}

//升级2.1.0数据库
func upgrade210(session *xorm.Session) error {
	if count, _ := session.Where("name = ?", "dingtalk").Table("notice_setting").Count(); count > 0 {
		return nil
	}
	_, err := session.Exec(`
INSERT INTO "notice_setting" (name,value,desc) VALUES ('dingtalk', 'disabled', '是否开启钉钉通知');
INSERT INTO "notice_setting" (name,value,desc) VALUES ('dingtalk_webhook', '', '钉钉通知webhook地址');
INSERT INTO "notice_setting" (name,value,desc) VALUES ('work_weixin', 'disabled', '是否开启企业微信通知');
INSERT INTO "notice_setting" (name,value,desc) VALUES ('work_weixin_webhook', '', '企业微信通知webhook通知');
ALTER TABLE "project" ADD COLUMN "desc" VARCHAR (1000) NOT NULL DEFAULT '';
`)
	return err
}
