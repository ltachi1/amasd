package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"amasd/core"
	"amasd/models"
	"strings"
	"unicode/utf8"
)

type Notice struct {
	core.BaseController
}

func (n *Notice) Setting(c *gin.Context) {
	if core.IsAjax(c) {
		ss := strings.Trim(c.DefaultPostForm("scrapyd_service", models.NoticeSettingStatusDisabled), " ")
		tf := strings.Trim(c.DefaultPostForm("task_finished", models.NoticeSettingStatusDisabled), " ")
		te := strings.Trim(c.DefaultPostForm("task_error", models.NoticeSettingStatusDisabled), " ")
		e := strings.Trim(c.DefaultPostForm("email", models.NoticeSettingStatusDisabled), " ")
		d := strings.Trim(c.DefaultPostForm("dingtalk", models.NoticeSettingStatusDisabled), " ")
		ww := strings.Trim(c.DefaultPostForm("work_weixin", models.NoticeSettingStatusDisabled), " ")
		sst := strings.Trim(c.DefaultPostForm("scrapyd_service_title", ""), " ")
		ssc := strings.Trim(c.DefaultPostForm("scrapyd_service_content", ""), " ")
		tft := strings.Trim(c.DefaultPostForm("task_finished_title", ""), " ")
		tfc := strings.Trim(c.DefaultPostForm("task_finished_content", ""), " ")
		tet := strings.Trim(c.DefaultPostForm("task_error_title", ""), " ")
		tec := strings.Trim(c.DefaultPostForm("task_error_content", ""), " ")
		es := strings.Trim(c.DefaultPostForm("email_smtp", ""), " ")
		esp := strings.Trim(c.DefaultPostForm("email_smtp_port", ""), " ")
		esa := strings.Trim(c.DefaultPostForm("email_sender_address", ""), " ")
		esp2 := strings.Trim(c.DefaultPostForm("email_sender_password", ""), " ")
		es2 := strings.Trim(c.DefaultPostForm("email_sender", ""), " ")
		ea := strings.Trim(c.DefaultPostForm("email_addressee", ""), " ")
		dw := strings.Trim(c.DefaultPostForm("dingtalk_webhook", ""), " ")
		www := strings.Trim(c.DefaultPostForm("work_weixin_webhook", ""), " ")
		if ss == models.NoticeSettingStatusEnabled && (sst == "" || ssc == "") {
			n.Fail(c, "system_notice_scrapyd_error")
			return
		}
		if tf == models.NoticeSettingStatusEnabled && (tft == "" || tfc == "") {
			n.Fail(c, "system_notice_task_finished_error")
			return
		}
		if te == models.NoticeSettingStatusEnabled && (tet == "" || tec == "") {
			n.Fail(c, "system_notice_task_error")
			return
		}
		if utf8.RuneCountInString(sst) > 50 || utf8.RuneCountInString(tft) > 50 || utf8.RuneCountInString(tet) > 50 {
			n.Fail(c, "extra_long_error", "通知标题", "50")
			return
		}
		if utf8.RuneCountInString(ssc) > 2000 || utf8.RuneCountInString(tfc) > 2000 || utf8.RuneCountInString(tec) > 2000 {
			n.Fail(c, "extra_long_error", "通知内容", "2000")
			return
		}
		if e == models.NoticeSettingStatusEnabled && (es == "" || esp == "" || esa == "" || esp2 == "" || ea == "") {
			n.Fail(c, "system_notice_email_error")
			return
		}
		if len(es) > 0 && !core.IsDomain(es) {
			n.Fail(c, "system_notice_email_smtp_error")
			return
		}
		if len(es) > 50 {
			n.Fail(c, "extra_long_error", "发件服务器", "50")
			return
		}
		if len(esp) > 0 && !core.IsNumber(esp) {
			n.Fail(c, "system_notice_email_port_error")
			return
		}
		if len(esp) > 5 {
			n.Fail(c, "extra_long_error", "发件服务器端口", "5")
			return
		}
		if len(esa) > 0 && !core.IsEmail(esa) {
			n.Fail(c, "system_notice_email_address_error")
			return
		}
		if utf8.RuneCountInString(esa) > 50 {
			n.Fail(c, "extra_long_error", "发件人邮箱地址", "50")
			return
		}
		if utf8.RuneCountInString(es2) > 50 {
			n.Fail(c, "extra_long_error", "发件人别名", "50")
			return
		}
		if len(esp2) > 20 {
			n.Fail(c, "extra_long_error", "发件人邮箱密码", "20")
			return
		}
		if len(ea) > 0 && !core.IsEmail(ea) {
			n.Fail(c, "system_notice_email_addressee_error")
			return
		}
		if len(ea) > 50 {
			n.Fail(c, "extra_long_error", "收件人邮箱地址", "50")
			return
		}
		if d == models.NoticeSettingStatusEnabled && dw == "" {
			n.Fail(c, "system_notice_dingtalk_error")
			return
		}
		if len(dw) > 200 {
			n.Fail(c, "extra_long_error", "钉钉webhook地址", "200")
			return
		}
		if !core.IsUrl(dw) {
			n.Fail(c, "system_notice_webhook_error")
			return
		}
		if ww == models.NoticeSettingStatusEnabled && www == "" {
			n.Fail(c, "system_notice_work_weixin_error")
			return
		}
		if len(www) > 200 {
			n.Fail(c, "extra_long_error", "企业微信webhook地址", "200")
			return
		}
		if !core.IsUrl(www) {
			n.Fail(c, "system_notice_webhook_error")
			return
		}

		fields := []core.B{
			{"name": "\"scrapyd_service\"", "value": ss},
			{"name": "\"task_finished\"", "value": tf},
			{"name": "\"task_error\"", "value": te},
			{"name": "\"email\"", "value": e},
			{"name": "\"dingtalk\"", "value": d},
			{"name": "\"work_weixin\"", "value": ww},
			{"name": "\"scrapyd_service_title\"", "value": sst},
			{"name": "\"scrapyd_service_content\"", "value": ssc},
			{"name": "\"task_finished_title\"", "value": tft},
			{"name": "\"task_finished_content\"", "value": tfc},
			{"name": "\"task_error_title\"", "value": tet},
			{"name": "\"task_error_content\"", "value": tec},
			{"name": "\"email_smtp\"", "value": es},
			{"name": "\"email_smtp_port\"", "value": esp},
			{"name": "\"email_sender_address\"", "value": esa},
			{"name": "\"email_sender_password\"", "value": esp2},
			{"name": "\"email_sender\"", "value": es2},
			{"name": "\"email_addressee\"", "value": ea},
			{"name": "\"dingtalk_webhook\"", "value": dw},
			{"name": "\"work_weixin_webhook\"", "value": www},
		}
		if new(models.NoticeSetting).Update(fields) {
			n.Success(c, nil)
		} else {
			n.Fail(c, "update_error")
		}
	} else {
		c.HTML(http.StatusOK, "notice/setting", gin.H{
			"settings": new(models.NoticeSetting).Find(),
		})
	}

}
