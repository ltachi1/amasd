package core

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"encoding/json"
)

// Request构造类
type Curl struct {
	cli             *http.Client
	req             *http.Request
	timeout         time.Duration
	response        *http.Response
	ResponseHeaders map[string]string
	ResponseBody    string
	headers         B
	cookies         B
	queries         B
	postData        A
	postFormData    *bytes.Buffer
}

// 创建一个Curl实例
func NewCurl() *Curl {
	return &Curl{timeout: time.Second * time.Duration(30)}
}

// 设置请求头
func (c *Curl) SetHeaders(headers B) *Curl {
	c.headers = headers
	return c
}

// 将用户自定义请求头添加到http.Request实例上
func (c *Curl) setHeaders() *Curl {
	for k, v := range c.headers {
		c.req.Header.Set(k, v)
	}
	if c.req.Header.Get("Content-Type") == "" {
		c.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return c
}

// 设置请求cookies
func (c *Curl) SetCookies(cookies B) *Curl {
	c.cookies = cookies
	return c
}

// 将用户自定义cookies添加到http.Request实例上
func (c *Curl) setCookies() *Curl {
	for k, v := range c.cookies {
		c.req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	return c
}

// 将用户自定义url查询参数添加到http.Request上
func (c *Curl) setQueries() *Curl {
	q := c.req.URL.Query()
	for k, v := range c.queries {
		q.Add(k, v)
	}
	c.req.URL.RawQuery = q.Encode()
	return c
}

// 发起get请求
func (c *Curl) Get(url string, params B) (string, error) {
	c.queries = params
	return c.execute(url, http.MethodGet)
}

// 发起post请求
func (c *Curl) Post(url string, params A) (string, error) {
	c.postData = params
	return c.execute(url, http.MethodPost)
}

func (c *Curl) PostForm(url string, postData *bytes.Buffer) (string, error) {
	c.postFormData = postData
	return c.execute(url, http.MethodPost)
}

// 发起Delete请求
func (c *Curl) Delete(url string) (string, error) {
	return c.execute(url, http.MethodDelete)
}

// 发起Put请求
func (c *Curl) Put(url string) (string, error) {
	return c.execute(url, http.MethodPut)
}

// 发起put请求
func (c *Curl) PATCH(url string) (string, error) {
	return c.execute(url, http.MethodPatch)
}

//设置超时时间
func (c *Curl) SetTimeOut(TimeOutSecond int) *Curl {
	c.timeout = time.Duration(TimeOutSecond)
	return c
}

//解析headers
func (c *Curl) parseHeaders() error {
	headers := map[string]string{}
	for k, v := range c.response.Header {
		headers[k] = v[0]
	}
	c.ResponseHeaders = headers
	return nil
}

//解析内容
func (c *Curl) parseBody() error {
	if body, err := ioutil.ReadAll(c.response.Body); err != nil {
		panic(err)
	} else {
		c.ResponseBody = string(body)
	}
	return nil
}

// 发起请求
func (c *Curl) execute(u string, method string) (string, error) {
	// 检测请求url是否填了
	if u == "" {
		return "", errors.New("请输入url地址")
	}
	if !IsUrl(u) {
		return "", errors.New("请输入正确url地址")
	}
	// 初始化http.Client对象
	c.cli = &http.Client{
		//这里不使用Transport(更高级的用法)，只使用Timeout设置整体超时
		//设置超时时间
		Timeout: c.timeout,
	}
	var reader io.Reader
	contentType := c.headers["Content-Type"]
	if strings.Index(contentType, "multipart/form-data") > -1 {
		reader = c.postFormData
	} else if strings.Index(contentType, "application/json") > -1 {
		bytesData, err := json.Marshal(c.postData)
		if err != nil {
			return "", errors.New("参数错误")
		}
		reader = bytes.NewReader(bytesData)
	} else {
		var payload url.Values
		if c.postData != nil {
			payload = map[string][]string{}
			for k, v := range c.postData {
				payload.Add(k, v.(string))
			}
		}
		reader = strings.NewReader(payload.Encode())
	}
	// 加载用户自定义的post数据到http.Request
	if req, err := http.NewRequest(method, u, reader); err != nil {
		return "", err
	} else {
		c.req = req
	}
	c.setHeaders().setCookies().setQueries()
	if resp, err := c.cli.Do(c.req); err != nil {
		return "", err
	} else {
		c.response = resp
	}
	defer c.response.Body.Close()
	c.parseHeaders()
	c.parseBody()
	if c.response.StatusCode != http.StatusOK {
		return "", errors.New(c.ResponseBody)
	}
	return c.ResponseBody, nil
}
