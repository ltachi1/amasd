package resource

import (
	"github.com/gin-gonic/gin"
	"scrapyd-admin/core"
	"scrapyd-admin/resource/views"
	"strings"
	"github.com/ltachi1/logrus"
	"html/template"
)

//加载模板文件
func LoadTemplate(r *gin.Engine) {
	//r.LoadHTMLGlob(core.RootPath + "/views/**/*")
	templateNames := views.AssetNames()
	t := template.New("").Delims("{{", "}}")
	for _, name := range templateNames {
		bytes, error := views.Asset(name)
		if error != nil {
			core.WriteLog(core.LogTypeInit, logrus.PanicLevel, nil, "模板文件初始化失败")
		}
		t.New(strings.Join(strings.Split(strings.Split(name, ".")[0], "/")[1:], "/")).Parse(string(bytes))
	}
	html := template.Must(t, nil)
	r.SetHTMLTemplate(html)

}
