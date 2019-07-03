// 验证token中间件
package core

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"scrapyd-admin/config"
	"time"
)

type tokenData struct {
	Id         int64 `json:"id"`
	Iat        int64 `json:"iat"`
	Exp        int64 `json:"exp"`
	LongestExp int64 `json:"longestExp"`
}

//初始验证，主要用于web的登录和api的token
type Auth interface {
	Check(c *gin.Context) //校验具体的信息
}

//web验证
type WebAuth struct{}

//api验证
type ApiAuth struct{}

func (w *WebAuth) Check(c *gin.Context) {
	session := GetSession()
	userInfo := session.Get(SessionUserInfoKey)
	if userInfo == nil {
		//如果没有登录信息则直接跳转到登录页面
		//c.Redirect(http.StatusFound, "/login")
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<script>top.location.href="/login";</script>`)
		c.Abort()
	} else {
		//重新更新session有效时间
		session.Set(SessionUserInfoKey, userInfo)
		session.Save()
	}
}

func (a *ApiAuth) Check(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		//解析不成功,或者其他原因
		c.JSON(http.StatusForbidden, config.PromptMsg["token_valid"])
		c.Abort()
		return
	}
	td := &tokenData{}
	if err := json.Unmarshal([]byte(AesDecrypt(token)), td); err != nil {
		c.JSON(http.StatusForbidden, config.PromptMsg["token_valid"])
		c.Abort()
		return
	}
	if td.LongestExp < time.Now().Unix() {
		//token已过最长使用期限，需要重新登录
		c.JSON(200, config.PromptMsg["no_login"])
		c.Abort()
		return
	} else {
		if td.Exp < time.Now().Unix() {
			//token过期，重新颁发token
			c.Header("Authorization", GenerateToken(td.Id))
		}
	}
}

//验证Token
func AuthValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		//解析不成功,或者其他原因
		c.JSON(http.StatusForbidden, config.PromptMsg["token_valid"])
		c.Abort()
		return
	}
	td := &tokenData{}
	if err := json.Unmarshal([]byte(AesDecrypt(token)), td); err != nil {
		c.JSON(http.StatusForbidden, config.PromptMsg["token_valid"])
		c.Abort()
		return
	}
	if td.LongestExp < time.Now().Unix() {
		//token已过最长使用期限，需要重新登录
		c.JSON(200, config.PromptMsg["no_login"])
		c.Abort()
		return
	} else {
		if td.Exp < time.Now().Unix() {
			//token过期，重新颁发token
			c.Header("Authorization", GenerateToken(td.Id))
		}
	}
}

//生成Token
func GenerateToken(id int64) string {
	timeNow := time.Now()
	info, _ := json.Marshal(tokenData{
		Id:         id,
		Iat:        timeNow.Unix(),                                      //签发时间
		Exp:        timeNow.Add(time.Minute * time.Duration(15)).Unix(), // 有效期15分钟
		LongestExp: timeNow.Add(time.Hour * time.Duration(2)).Unix(),    // 两个小时之内可以重新更新令牌，超过则需要重新登录
	})
	return AesEncrypt(string(info))
}
