package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct{}

func (b *BaseController) Success(c *gin.Context, data ...gin.H) {
	msg := gin.H{}
	if len(data) == 0 {
		msg = PromptMsg["success"]
	} else {
		msg = data[0]
		if _, exists := msg["code"]; !exists {
			msg["code"] = PromptMsg["success"]["code"]
		}
	}
	c.JSON(http.StatusOK, msg)
}

func (b *BaseController) Fail(c *gin.Context, data ...gin.H) {
	msg := gin.H{}
	if len(data) == 0 {
		msg = PromptMsg["fail"]
	} else {
		msg = data[0]
		if _, exists := msg["code"]; !exists {
			msg["code"] = PromptMsg["fail"]["code"]
		}
	}
	c.JSON(http.StatusOK, msg)
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