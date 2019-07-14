// 日志处理，使用logrus日志框架
package core

import (
	"bufio"
	"fmt"
	"github.com/ltachi1/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/ltachi1/logrus"
	"os"
	"strings"
	"sync"
	"time"
)

var Log = logrus.New()

//初始化日志
func InitLog(logPath string) {
	if logPath == "" {
		logPath = RootPath + "/logs"
	}
	logPath = SupplementDir(logPath)
	err := os.MkdirAll(logPath, os.ModePerm)
	if err != nil {
		fmt.Println("日志路径输入不正确或者没有写入权限")
		os.Exit(1)
	}

	//增加日志中间件，将日志按天按类型输出到不同的文件，error以上信息单独输出目录
	writer, err := rotatelogs.New(
		logPath + "{dir}%Y%m%d.log",
		//rotatelogs.WithLinkName(./logs),     // 生成软链，指向最新日志文件
		rotatelogs.WithRotationTime(time.Hour*time.Duration(24)), // 日志切割时间间隔
	)
	if err != nil {
		Log.Errorf("钩子初始化失败. %v", errors.WithStack(err))
	}
	setNull()
	//根据不同的环境设置不同的日志等级
	switch Env {
	case EnvDev:
		Log.SetLevel(logrus.TraceLevel)
	case EnvTesting:
		Log.SetLevel(logrus.DebugLevel)
	case EnvProduction:
		Log.SetLevel(logrus.InfoLevel)
	default:
		Log.SetLevel(logrus.ErrorLevel)
	}
	lfHook := NewHook(WriterMap{
		logrus.TraceLevel: Stdout,     //跟踪日志,跟踪日志不会输出到文件
		logrus.DebugLevel: RotateLogs, //调试日志
		logrus.InfoLevel:  RotateLogs, //代表业务日志
		logrus.WarnLevel:  RotateLogs, //警告日志
		logrus.ErrorLevel: RotateLogs, //错误日志
		logrus.FatalLevel: RotateLogs, //致命等级日志
		logrus.PanicLevel: RotateLogs, //恐慌(导致程序终端)等级日志
	}, &logrus.TextFormatter{DisableColors: true}, writer)
	Log.AddHook(lfHook)
}

//丢弃logurs的默认日志输出
func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	Log.SetOutput(writer)
}

//一下是logurs钩子实现

//所有输出类型
const (
	Stdout     int32 = iota
	Stderr
	RotateLogs
)

// 默认日志格式化方式
var defaultFormatter = &logrus.TextFormatter{DisableColors: true}

// 不同日志等级的写入类型
type WriterMap map[logrus.Level]int32

// 用来向特定路径下写入日志的钩子
type writePathHook struct {
	writers   WriterMap
	rl        *rotatelogs.RotateLogs
	levels    []logrus.Level
	lock      *sync.Mutex
	formatter logrus.Formatter
}

// 生成一个新的钩子
func NewHook(output WriterMap, formatter logrus.Formatter, rl *rotatelogs.RotateLogs) *writePathHook {
	hook := &writePathHook{
		lock: new(sync.Mutex),
	}
	hook.SetFormatter(formatter)
	hook.rl = rl
	hook.writers = output
	for level := range output {
		hook.levels = append(hook.levels, level)
	}

	return hook
}

// 设置日志格式化方式
func (hook *writePathHook) SetFormatter(formatter logrus.Formatter) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if formatter == nil {
		formatter = defaultFormatter
	} else {
		switch formatter.(type) {
		case *logrus.TextFormatter:
			textFormatter := formatter.(*logrus.TextFormatter)
			textFormatter.DisableColors = true
		}
	}
	hook.formatter = formatter
}

// 实现rotatelogs Hook接口
func (hook *writePathHook) Fire(entry *logrus.Entry) error {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if hook.writers != nil {
		return hook.ioWrite(entry)
	}
	return nil
}

// 实现rotatelogs Hook接口
func (hook *writePathHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// 写日志
func (hook *writePathHook) ioWrite(entry *logrus.Entry) error {
	var (
		writerType int32
		msg        []byte
		err        error
		ok         bool
	)
	if writerType, ok = hook.writers[entry.Level]; !ok {
		return nil
	}
	// 格式化输出内容
	msg, err = hook.formatter.Format(entry)
	if err != nil {
		Log.Println("输出字符串格式化错误", err)
		return err
	}
	if writerType == Stdout {
		_, err = os.Stdout.Write(msg)
	} else if writerType == Stderr {
		_, err = os.Stderr.Write(msg)
	} else if writerType == RotateLogs {
		_, err = hook.rl.Write(msg, hook.formatPath(entry.Path))
	}
	return err
}

//格式化路径字符串，传递的格式是x.x.x需要格式化成/x/x/x/
func (hook *writePathHook) formatPath(path string) string {
	if path != "" {
		path = strings.Replace(path, ".", "/", -1)
		path += "/"
	}
	return path
}

//记录日志
func WriteLog(logType string, logLevel logrus.Level, fields logrus.Fields, log ...interface{}) {
	go func() {
		logPath := ""
		if logLevel >= logrus.WarnLevel {
			logPath = logType
		} else {
			//错误日志单独输出到一个文件夹下,类型用你logType区分
			logPath = ErrorDir
			if fields != nil {
				fields["error_type"] = logType
			} else {
				fields = logrus.Fields{"error_type": logType}
			}
		}
		e := Log.WithPath(logPath)
		if fields != nil {
			e = e.WithFields(fields)
		}
		switch logLevel {
		case logrus.TraceLevel:
			e.Trace(log...)
		case logrus.DebugLevel:
			e.Debug(log...)
		case logrus.InfoLevel:
			e.Info(log...)
		case logrus.WarnLevel:
			e.Warn(log...)
		case logrus.ErrorLevel:
			e.Error(log...)
		case logrus.FatalLevel:
			e.Fatal(log...)
		case logrus.PanicLevel:
			e.Panic(log...)
		}
	}()

}
