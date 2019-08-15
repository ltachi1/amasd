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
	"daemonStatus": core.B{"url": "%s/daemonstatus.json", "method": http.MethodGet},
	"addVersion":   core.B{"url": "%s/addversion.json", "method": http.MethodPost},
	"delProject":   core.B{"url": "%s/delproject.json", "method": http.MethodPost},
	"listSpiders":  core.B{"url": "%s/listspiders.json", "method": http.MethodGet},
	"schedule":     core.B{"url": "%s/schedule.json", "method": http.MethodPost},
	"cancel":       core.B{"url": "%s/cancel.json", "method": http.MethodPost},
	"listjobs":     core.B{"url": "%s/listjobs.json", "method": http.MethodGet},
}

//检查服务是否可用
func (s *Scrapyd) DaemonStatus() error {
	var (
		err error
		str string
	)
	if str, err = s.send("daemonStatus", core.A{}); err == nil {
		daemonStatus := core.A{}
		if err = json.Unmarshal(core.Str2bytes(str), &daemonStatus); err == nil {
			if status, ok := daemonStatus["status"]; ok && status == "ok" {
				return nil
			}
		}
	}
	return err
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
	if s.Auth == ServerAuthOpen {
		headers["Authorization"] = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.Username, s.Password))))
	}
	str, error := core.NewCurl().SetHeaders(headers).PostForm(core.CompletionUrl(fmt.Sprintf(scrapydUrls["addVersion"]["url"], s.Host)), bodyBuffer)
	if error == nil {
		daemonStatus := core.A{}
		if error = json.Unmarshal(core.Str2bytes(str), &daemonStatus); error == nil {
			if status, ok := daemonStatus["status"]; ok && status == "ok" {
				return true
			}
		}
	}
	if error != nil {
		core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"host": s.Host, "project_name": project.Name, "version": project.LastVersion}, fmt.Sprintf("scrpayd项目添加失败:%s", error))
	}
	return false
}

//删除项目
func (s *Scrapyd) DelProject(projectName string) bool {
	var error error
	if str, error := s.send("delProject", core.A{"project": projectName}); error == nil {
		result := core.A{}
		if error = json.Unmarshal(core.Str2bytes(str), &result); error == nil {
			if status, ok := result["status"]; ok && status == "ok" {
				return true
			}
		}
	}
	if error != nil {
		core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"host": s.Host, "project_name": projectName}, fmt.Sprintf("scrpayd项目删除失败:%s", error))
	}
	return false
}

//获取项目中所包含的爬虫列表
func (s *Scrapyd) ListSpiders(project *Project) (error, []string) {
	var (
		err error
		str string
	)
	spiders := make([]string, 0)
	if str, err = s.send("listSpiders", core.A{"project": project.Name, "_version": project.LastVersion}); err == nil {
		result := core.A{}
		if err = json.Unmarshal(core.Str2bytes(str), &result); err == nil {
			if status, ok := result["status"]; ok && status == "ok" {
				for _, sp := range result["spiders"].([]interface{}) {
					spiders = append(spiders, sp.(string))
				}
				return nil, spiders
			}
		}
	}
	return err, spiders
}

//投递任务
func (s *Scrapyd) Schedule(projectName string, version string, spiderName string) (error, string) {
	var (
		err error
		str string
	)
	if str, err = s.send("schedule", core.A{"project": projectName, "_version": version, "spider": spiderName}); err == nil {
		result := core.A{}
		if err = json.Unmarshal(core.Str2bytes(str), &result); err == nil {
			if status, ok := result["status"]; ok && status == "ok" {
				return nil, result["jobid"].(string)
			}
		}
	}
	return err, ""
}

//取消任务
func (s *Scrapyd) Cancel(projectName string, jobId string) error {
	var (
		err error
		str string
	)
	if str, err = s.send("cancel", core.A{"project": projectName, "job": jobId}); err == nil {
		result := core.A{}
		if err = json.Unmarshal(core.Str2bytes(str), &result); err == nil {
			if status, ok := result["status"]; ok && status == "ok" {
				return nil
			}
		}
	}
	return err
}

//获取制定服务器项目下的任务列表
func (s *Scrapyd) ListJobs(projectName string) (error, map[string][]interface{}) {
	var (
		err error
		str string
	)
	if str, err = s.send("listjobs", core.A{"project": projectName}); err == nil {
		result := core.A{}
		if err = json.Unmarshal(core.Str2bytes(str), &result); err == nil {
			if status, ok := result["status"]; ok && status == "ok" {
				return nil, map[string][]interface{}{
					"pending":  result["pending"].([]interface{}),
					"running":  result["running"].([]interface{}),
					"finished": result["finished"].([]interface{}),
				}
			}
		}
	}
	return err, map[string][]interface{}{}
}

func (s *Scrapyd) send(key string, params core.A) (string, error) {
	headers := make(core.B, 0)
	if s.Auth == ServerAuthOpen {
		headers["Authorization"] = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.Username, s.Password))))
	}
	v := scrapydUrls[key]
	url := core.CompletionUrl(fmt.Sprintf(v["url"], s.Host))
	if v["method"] == http.MethodPost {
		return core.NewCurl().SetHeaders(headers).Post(url, params)
	} else if v["method"] == http.MethodGet {
		b := core.B{}
		for k, v := range params {
			b[k] = v.(string)
		}
		return core.NewCurl().SetHeaders(headers).Get(url, b)
	}
	return "", errors.New("请输入正确url地址")
}
