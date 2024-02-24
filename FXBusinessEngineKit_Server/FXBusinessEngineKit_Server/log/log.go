package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type RecordType int

const (
	Err RecordType = 1
	Msg RecordType = 2
)

var errorRecorder = logrus.New()
var msgRecorder = logrus.New()

func InitLog() {
	initErrRecorder()
	initMesRecorder()
}

func initErrRecorder() {
	// 设置日志解析格式
	errorRecorder.SetFormatter(&logrus.JSONFormatter{})
	// 设置日志输出到标准输出
	errorRecorder.SetOutput(os.Stdout)
	write, err := os.OpenFile("appError.log", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("日志文件输出初始化失败", err)
		return
	}
	// 设置日志输出文件
	errorRecorder.SetOutput(write)
	// 日志记录等级
	errorRecorder.SetLevel(logrus.InfoLevel)
	// 设置日志记录当前调用方法
	errorRecorder.SetReportCaller(true)
}

func initMesRecorder() {
	// 设置日志解析格式
	msgRecorder.SetFormatter(&logrus.JSONFormatter{})
	// 设置日志输出到标准输出
	msgRecorder.SetOutput(os.Stdout)
	write, err := os.OpenFile("appMes.log", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("日志文件输出初始化失败", err)
		return
	}
	// 设置日志输出文件
	msgRecorder.SetOutput(write)
	// 日志记录等级
	msgRecorder.SetLevel(logrus.InfoLevel)
	// 设置日志记录当前调用方法
	msgRecorder.SetReportCaller(true)
}

func RecordLog(style RecordType, content string) {
	if style == Err {
		recordError(content, "")
	}
	if style == Msg {
		recordMsg(content)
	}
}

// LogErrorRecord 日志记录
func recordError(content string, local string) {
	errorRecorder.WithFields(map[string]interface{}{
		"AppError": content,
	}).Warn()
	fmt.Printf(fmt.Sprintf("\n日志输出：%s,位置：%s\n", content, local))

}

func recordMsg(content string) {
	msgRecorder.WithFields(map[string]interface{}{
		"AppMessage": content,
	}).Info()
}

// 实现一个钩子，用于在每次日志记录时记录调用API的用户id

type LogrusHook struct {
}

func (hook *LogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *LogrusHook) Fire(entry *logrus.Entry) error {
	entry.Data["钩子信息"] = ""
	return nil
}
