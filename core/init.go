// 核心包初始化
package core

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
	"os"
	"log"
	"strings"
	"flag"
	"fmt"
)

func init() {
	RootPath = getCurrentDirectory()
	args := parseArgs()
	Env = args["env"]
	HttpPort = args["port"]

	if Env == EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	//初始化日志
	InitLog(args["logPath"])
	//初始化数据库
	InitDb(args["dbPath"])
	//初始化定时任务
	InitCron()
}

//获取当前运行的绝对路径
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//解析命令行参数
func parseArgs() B {
	args := make(B, 0)
	var (
		port    string
		env     string
		logPath string
		dbPath  string
	)
	//加载命令行参数
	flag.StringVar(&port, "p", "8000", "监听端口")
	flag.StringVar(&env, "e", EnvDev, fmt.Sprintf("运行环境:%s|%s|%s", EnvDev, EnvTesting, EnvProduction))
	flag.StringVar(&logPath, "log", "", "日志库存放路径，默认当前目录")
	flag.StringVar(&dbPath, "db", "", "数据库存放路径，默认当前目录")
	flag.Parse()
	if !(env == EnvDev || env == EnvTesting || env == EnvProduction) {
		flag.Usage()
		os.Exit(1)
	}
	args["port"] = port
	args["env"] = env
	args["logPath"] = logPath
	args["dbPath"] = dbPath
	return args
}