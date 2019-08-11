package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

type BaseController struct{}

func (b *BaseController) Success(c *gin.Context, data A, message ...interface{}) {
	msg := ""
	if len(message) == 0 {
		msg = PromptMsg["success"]
	} else {
		msg = PromptMsg[message[0].(string)]
		if len(message) > 1 {
			msg = fmt.Sprintf(msg, message[1:]...)
		}
	}
	c.JSON(http.StatusOK, A{"code": 0, "msg": msg, "data": data})
}

func (b *BaseController) Fail(c *gin.Context, data ...interface{}) {
	msg := ""
	if len(data) == 0 {
		msg = PromptMsg["fail"]
	} else {
		msg = PromptMsg[data[0].(string)]
		if len(data) > 1 {
			msg = fmt.Sprintf(msg, data[1:]...)
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": msg})
}


type BaseModel struct {
}

func (b *BaseModel) Insert(obj interface{}) bool {
	_, error := Db.InsertOne(obj)
	if error != nil {
		return false
	}
	return true
}

func (b *BaseModel) Delete(obj interface{}) bool {
	_, error := Db.Delete(obj)
	if error != nil {
		return false
	}
	return true
}