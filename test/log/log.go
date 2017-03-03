package main

import (
	"github.com/astaxie/beego/logs"
	//"github.com/admpub/log"
)

func Log(){
	log := logs.NewLogger(10)
	log.SetLogger("console","")
	log.EnableFuncCallDepth(true)
	log.SetLogFuncCallDepth(5)

	log.Trace("[update]trace")
	log.Info("info")
	log.Warn("warning")
	log.Debug("debug")
	log.Critical("critical")
	log.Error("error")
	log.Emergency("emergency")

	log.SetLogger("file",`{"filename":"log.log"}`)
	log.Trace("trace")
	log.Info("info")
	log.Warn("warning")
	log.Debug("debug")
	log.Critical("critical")
	log.Error("error")
	log.Emergency("emergency")
}

func Le1(){
	Log()
}
func Le2() {
	Le1()
}
var l *logs.BeeLogger
func init(){
	l = logs.NewLogger(10)
}
func loglog(){
	l.Info("in the log")
}

func main() {

	Le2()

	logs.Info("--------------------------------------------")
	logs.Critical("critical")
	logs.Error("error")
	l.Info("hahahah############")
	loglog()
	/*
	// 创建根记录器(root logger)
	logger := log.NewLogger()

	// 添加一个控制台标的（Console Target）和一个文件标的（File Target）
	t1 := log.NewConsoleTarget()
	t2 := log.NewFileTarget()
	t2.FileName = "app.log"
	t2.MaxLevel = log.LevelError
	logger.Targets = append(logger.Targets, t1, t2)

	logger.Open()
	defer logger.Close()

	// 调用不同的记录方法记录不同的日志信息。
	logger.Error("plain text error")
	logger.Error("error with format: %v", true)
	logger.Debug("some debug info")

	// 自定义日志类别
	l := logger.GetLogger("app.services")
	l.Info("some info")
	l.Warn("some warning")

	*/


}
