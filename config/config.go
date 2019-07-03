// 加载所有配置
package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/gin-gonic/gin"
)

var Conf *config

type config struct {
	Http struct {
		Port string `yaml:"port"`
	} `yaml:"http"`
	Db                db     `yaml:"db"`
	BusinessLogPath   string `yaml:"businessLogPath"`
	SessionExpires    int    `yaml:"sessionExpires"`
	SessionCookieName string `yaml:"sessionCookieName"`
	Env               string `yaml:"env"`
}

type db struct {
	TablePrefix string  `yaml:"table_prefix"`
	Master      dbField `yaml:"master"`
	Slave struct {
		MaxIdleConns int       `yaml:"max_idle_conns"`
		MaxOpenConns int       `yaml:"max_open_conns"`
		List         []dbField `yaml:"list"`
	} `yaml:"slave"`
}

//对应数据库配置文件
type dbField struct {
	Host         string `yaml:"host"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	Charset      string `yaml:"charset"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type redis struct {
	Master redisField `yaml:"master"`
	Slave struct {
		MaxIdleConns int          `yaml:"max_idle_conns"`
		MaxOpenConns int          `yaml:"max_open_conns"`
		List         []redisField `yaml:"list"`
	} `yaml:"slave"`
}

//对应redis配置文件
type redisField struct {
	Host         string `yaml:"host"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	Port         string `yaml:"port"`
	Timeout      string `yaml:"timeout"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

//加载配置文件
func (c *config) loadConf(filePath string) *config {
	//应该是 绝对地址
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

//获取当前运行的绝对路径
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//判断文件是否存在
func fileExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err == nil {
		if fileInfo.IsDir() {
			return false
		}
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func init() {
	RootPath = getCurrentDirectory()
	configFile := RootPath + "/config.yml"
	if !fileExists(configFile) {
		fmt.Println("请确保当前路径下有config.yml文件")
		os.Exit(1)
	}
	Conf = Conf.loadConf(configFile)
	if !(Conf.Env == EnvDev || Conf.Env == EnvTesting || Conf.Env == EnvProduction) {
		fmt.Println("配置文件中env设置错误，只能是dev | testing | production其中的一种")
		os.Exit(1)
	}
	if Conf.Env == EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}
}
