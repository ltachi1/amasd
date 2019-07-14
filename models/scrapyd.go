package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"scrapyd-admin/core"
	"github.com/ltachi1/logrus"
)

type Scrapyd struct {
	Host     string
	Auth     uint8
	Username string
	Password string
}

var scrapydUrls = core.C{
	"daemonStatus": core.B{"url": "http://%s/daemonstatus.json", "method": http.MethodGet},
	"addVersion":   core.B{"url": "http://%s/addversion.json", "method": http.MethodPost},
	"delProject":   core.B{"url": "http://%s/delproject.json", "method": http.MethodPost},
	"listSpiders":  core.B{"url": "http://%s/listspiders.json", "method": http.MethodGet},
	"schedule":     core.B{"url": "http://%s/schedule.json", "method": http.MethodPost},
	"cancel":       core.B{"url": "http://%s/cancel.json", "method": http.MethodPost},
	"listjobs":     core.B{"url": "http://%s/listjobs.json", "method": http.MethodGet},
}

//检查服务是否可用
func (s *Scrapyd) DaemonStatus() bool {
	var error error
	if str, error := s.send("daemonStatus", core.B{}); error == nil {
		daemonStatus := core.A{}
		if error = json.Unmarshal(core.Str2bytes(str), &daemonStatus); error == nil {
			if status, ok := daemonStatus["status"]; ok && status == "ok" {
				return true
			}
		}
	}
	if error != nil {
		core.WriteLog(core.LogTypeScrapyd, logrus.ErrorLevel, logrus.Fields{"host": s.Host}, error)
	}
	return false
}

//添加scrapy项目
func (s *Scrapyd) AddVersion(project *Project, file *multipart.FileHeader) bool {
	//构建字节缓冲区
	bodyBuffer := &bytes.Buffer{}
	//创建form表单格式的输出流
	bodyWriter := multipart.NewWriter(bodyBuffer)
	//创建header
	fw, error := bodyWriter.CreateFormFile("egg", "output.egg")
	if error != nil {
		return false
	}
	//将输入流复制到输出流
	fr, error := file.Open()
	if error != nil {
		return false
	}
	io.Copy(fw, fr)
	//创建表单其他数据
	bodyWriter.WriteField("project", project.Name)
	bodyWriter.WriteField("version", project.LastVersion)
	bodyWriter.Close()
	headers := core.B{
		"Content-Type": bodyWriter.FormDataContentType(),
	}
	if s.Username != "" && s.Password != "" {
		headers["Authorization"] = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.Username, s.Password))))
	}
	str, error := core.NewCurl().SetHeaders(headers).PostForm(fmt.Sprintf(scrapydUrls["addVersion"]["url"], s.Host), bodyBuffer)
	if error == nil {
		daemonStatus := core.A{}
		if error = json.Unmarshal(core.Str2bytes(str), &daemonStatus); error == nil {
			if status, ok := daemonStatus["status"]; ok && status == "ok" {
				return true
			}
		}
	}
	if error != nil {
		core.WriteLog(core.LogTypeScrapyd, logrus.ErrorLevel, logrus.Fields{"host": s.Host}, error)
	}
	return false
}

//删除项目
func (s *Scrapyd) DelProject(projectName string) bool {
	var error error
	if str, error := s.send("delProject", core.B{"project": projectName}); error == nil {
		result := core.A{}
		if error = json.Unmarshal(core.Str2bytes(str), &result); error == nil {
			if status, ok := result["status"]; ok && status == "ok" {
				return true
			}
		}
	}
	if error != nil {
		core.WriteLog(core.LogTypeScrapyd, logrus.ErrorLevel, logrus.Fields{"host": s.Host}, error)
	}
	return false
}

//获取项目中所包含的爬虫列表
func (s *Scrapyd) ListSpiders(project *Project) []string {
	spiders := make([]string, 0)
	var error error
	if str, error := s.send("listSpiders", core.B{"project": project.Name, "_version": project.LastVersion}); error == nil {
		result := core.A{}
		if error = json.Unmarshal(core.Str2bytes(str), &result); error == nil {
			if status, ok := result["status"]; ok && status == "ok" {
				for _, sp := range result["spiders"].([]interface{}) {
					spiders = append(spiders, sp.(string))
				}
			}
		}
	}
	if error != nil {
		core.WriteLog(core.LogTypeScrapyd, logrus.ErrorLevel, logrus.Fields{"host": s.Host}, error)
	}
	return spiders
}

func (s *Scrapyd) Schedule(projectName string, version string, spiderName string) (bool, string) {
	var error error
	if str, error := s.send("schedule", core.B{"project": projectName, "_version": version, "spider": spiderName}); error == nil {
		result := core.A{}
		if error = json.Unmarshal(core.Str2bytes(str), &result); error == nil {
			if status, ok := result["status"]; ok && status == "ok" {
				return true, result["jobid"].(string)
			}
		}
	}
	if error != nil {
		core.WriteLog(core.LogTypeScrapyd, logrus.ErrorLevel, logrus.Fields{"host": s.Host}, error)
	}
	return false, ""
}

func (s *Scrapyd) Cancel(projectName string, jobId string) bool {
	var error error
	if str, error := s.send("cancel", core.B{"project": projectName, "job": jobId}); error == nil {
		result := core.A{}
		if error = json.Unmarshal(core.Str2bytes(str), &result); error == nil {
			if status, ok := result["status"]; ok && status == "ok" {
				return true
			}
		}
	}
	if error != nil {
		core.WriteLog(core.LogTypeScrapyd, logrus.ErrorLevel, logrus.Fields{"host": s.Host}, error)
	}
	return false
}

func (s *Scrapyd) ListJobs(projectName string) (bool, map[string][]interface{}) {
	var error error
	if str, error := s.send("listjobs", core.B{"project": projectName}); error == nil {
		result := core.A{}
		if error = json.Unmarshal(core.Str2bytes(str), &result); error == nil {
			if status, ok := result["status"]; ok && status == "ok" {
				return true, map[string][]interface{}{
					"pending":  result["pending"].([]interface{}),
					"running":  result["running"].([]interface{}),
					"finished": result["finished"].([]interface{}),
				}
			}
		}
	}
	if error != nil {
		core.WriteLog(core.LogTypeScrapyd, logrus.ErrorLevel, logrus.Fields{"host": s.Host}, error)
	}
	return false, map[string][]interface{}{}
}

func (s *Scrapyd) send(key string, params core.B) (string, error) {
	headers := make(core.B, 0)
	if s.Username != "" && s.Password != "" {
		headers["Authorization"] = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.Username, s.Password))))
	}
	v := scrapydUrls[key]
	if v["method"] == http.MethodPost {
		return core.NewCurl().SetHeaders(headers).Post(fmt.Sprintf(v["url"], s.Host), params)
	} else if v["method"] == http.MethodGet {
		return core.NewCurl().SetHeaders(headers).Get(fmt.Sprintf(v["url"], s.Host), params)
	}
	return "", errors.New("请输入正确url地址")
}
