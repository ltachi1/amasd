package resource

import (
	"github.com/gin-gonic/gin"
	"scrapyd-admin/config"
	"strings"
	"html/template"
	"scrapyd-admin/resource/views"
	"github.com/ltachi1/logrus"
	"scrapyd-admin/core"
)

//加载模板文件
func LoadTemplate(r *gin.Engine) {
	//r.LoadHTMLGlob(config.RootPath + "/views/**/*")
	templateNames := views.AssetNames()
	t := template.New("").Delims("{{", "}}")
	for _, name := range templateNames {
		bytes, error := views.Asset(name)
		if error != nil {
			core.WriteLog(config.LogTypeInit, logrus.PanicLevel, nil, "模板文件初始化失败")
		}
		t.New(strings.Join(strings.Split(strings.Split(name, ".")[0], "/")[1:], "/")).Parse(string(bytes))
	}
	html := template.Must(t, nil)
	r.SetHTMLTemplate(html)

}
