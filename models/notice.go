package models

import (
	"scrapyd-admin/core"
	"gopkg.in/gomail.v2"
	"fmt"
	"strconv"
	"github.com/ltachi1/logrus"
	"strings"
	"encoding/json"
)

type Notice interface {
	send(title string, content string)
}

type EmailNotice struct{}

func (e *EmailNotice) send(title string, content string) {
	//校验邮箱相关配置
	if !core.IsEmail(core.NoticeSettings["email_sender_address"]) ||
		!core.IsEmail(core.NoticeSettings["email_addressee"]) ||
		core.NoticeSettings["email_smtp"] == "" ||
		core.NoticeSettings["email_sender_password"] == "" ||
		core.NoticeSettings["email_smtp_port"] == "" {
		core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{}, "邮件发送失败,邮件配置错误")
		return
	}
	port, err := strconv.Atoi(core.NoticeSettings["email_smtp_port"])
	if err != nil {
		core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{"port": core.NoticeSettings["email_smtp_port"]}, fmt.Sprintf("邮件发送失败,端口设置错误:%s", err))
		return
	}
	m := gomail.NewMessage()
	//设置发件人别名
	if core.NoticeSettings["email_sender"] == "" {
		m.SetHeader("From", core.NoticeSettings["email_sender_address"])
	} else {
		m.SetHeader("From", m.FormatAddress(core.NoticeSettings["email_sender_address"], core.NoticeSettings["email_sender"]))
	}
	m.SetHeader("To", core.NoticeSettings["email_addressee"])
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)
	d := gomail.NewDialer(core.NoticeSettings["email_smtp"], port, core.NoticeSettings["email_sender_address"], core.NoticeSettings["email_sender_password"])
	if err := d.DialAndSend(m); err != nil {
		core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{"title": title, "content": content}, fmt.Sprintf("邮件发送失败:%s", err))
	}
}

type DingtalkNotice struct {
	Errcode int
	Errmsg  string
}

func (d *DingtalkNotice) send(title string, content string) {
	if core.NoticeSettings["dingtalk_webhook"] == "" || !core.IsUrl(core.NoticeSettings["dingtalk_webhook"]) {
		core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{}, "钉钉通知失败,配置错误")
		return
	}
	param := core.A{"msgtype": "text", "text": core.B{"content": fmt.Sprintf("%s\n%s", title, content)}}
	if result, err := core.NewCurl().SetHeaders(core.B{"Content-Type": "application/json"}).Post(core.NoticeSettings["dingtalk_webhook"], param); err != nil {
		core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{"title": title, "content": content}, fmt.Sprintf("钉钉通知失败:%s", err))
	} else {
		if len(result) > 0 {
			if err := json.Unmarshal([]byte(result), &d); err != nil {
				core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{"title": title, "content": content}, fmt.Sprintf("钉钉通知返回结果解析失败:%s", err))
				return
			}
			if d.Errcode != 0 {
				core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{"title": title, "content": content}, fmt.Sprintf("钉钉通知失败:%s", d.Errmsg))
			}
		}
	}
}

type WorkWeixinNotice struct {
	Errcode int
	Errmsg  string
}

func (w *WorkWeixinNotice) send(title string, content string) {
	if core.NoticeSettings["work_weixin_webhook"] == "" || !core.IsUrl(core.NoticeSettings["work_weixin_webhook"]) {
		core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{}, "企业微信通知失败,配置错误")
		return
	}
	param := core.A{"msgtype": "text", "text": core.B{"content": fmt.Sprintf("%s\n%s", title, content)}}
	if result, err := core.NewCurl().SetHeaders(core.B{"Content-Type": "application/json"}).Post(core.NoticeSettings["work_weixin_webhook"], param); err != nil {
		core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{"title": title, "content": content}, fmt.Sprintf("企业微信通知失败:%s", err))
	} else {
		if len(result) > 0 {
			if err := json.Unmarshal([]byte(result), &w); err != nil {
				core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{"title": title, "content": content}, fmt.Sprintf("企业微信通知返回结果解析失败:%s", err))
				return
			}
			if w.Errcode != 0 {
				core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{"title": title, "content": content}, fmt.Sprintf("企业微信通知失败:%s", w.Errmsg))
			}
		}
	}
}

//通知类型调度
func noticeTypeDispatch(title string, content string) {
	if core.NoticeSettings["email"] == NoticeSettingStatusEnabled {
		//开启邮件通知
		send(new(EmailNotice), title, content)
	}
	if core.NoticeSettings["dingtalk"] == NoticeSettingStatusEnabled {
		//开启钉钉通知
		send(new(DingtalkNotice), title, content)
	}
	if core.NoticeSettings["work_weixin"] == NoticeSettingStatusEnabled {
		//开启企业微信通知
		send(new(WorkWeixinNotice), title, content)
	}
	//后续增加其他通知类型
}

//通知项调度
func noticeOptionsDispatch(option string, b core.B) {
	title, content := "", ""
	if core.NoticeSettings[option] == NoticeSettingStatusEnabled {
		if option == NoticeOptionScrapyd {
			title, content = core.NoticeSettings["scrapyd_service_title"], core.NoticeSettings["scrapyd_service_content"]
		} else if option == NoticeOptionTaskFinished {
			title, content = core.NoticeSettings["task_finished_title"], core.NoticeSettings["task_finished_content"]
		} else if option == NoticeOptionTaskError {
			title, content = core.NoticeSettings["task_error_title"], core.NoticeSettings["task_error_content"]
		}
	}

	if len(title) == 0 || len(content) == 0 {
		return
	}
	for _, w := range NoticeSettingWildcards[option] {
		if strings.Index(content, w) > -1 {
			content = strings.Replace(content, w, b[strings.Trim(w, "{}")], -1)
		}
	}
	noticeTypeDispatch(title, content)
}

func send(n Notice, title string, content string) {
	n.send(title, content)
}
