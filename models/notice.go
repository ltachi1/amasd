package models

import (
	"scrapyd-admin/core"
	"gopkg.in/gomail.v2"
	"fmt"
	"strconv"
	"github.com/ltachi1/logrus"
	"strings"
)

type Notice interface {
	send(b core.B)
}

type EmailNotice struct {
	Smtp           string
	SenderEmail    string
	Sender         string
	Password       string
	AddresseeEmail string
}

func (e *EmailNotice) send(b core.B) {
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
	m.SetHeader("Subject", b["title"])
	m.SetBody("text/html", b["content"])
	d := gomail.NewDialer(core.NoticeSettings["email_smtp"], port, core.NoticeSettings["email_sender_address"], core.NoticeSettings["email_sender_password"])
	if err := d.DialAndSend(m); err != nil {
		core.WriteLog(core.LogTypeSystem, logrus.ErrorLevel, logrus.Fields{"title": b["title"], "content": b["content"]}, fmt.Sprintf("邮件发送失败:%s", err))
	}
}

//通知类型调度
func noticeTypeDispatch(b core.B) {
	if core.NoticeSettings["email"] == NoticeSettingStatusEnabled {
		//开启邮件通知
		send(new(EmailNotice), b)
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
	noticeTypeDispatch(core.B{"title": title, "content": content})
}

func send(n Notice, b core.B) {
	n.send(b)
}
