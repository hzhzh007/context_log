package context_log

import (
	"container/list"
	"fmt"
	"github.com/op/go-logging"
	"os"
	"time"
)

const (
	LOG_NAME = "log"
)

var (
	Log               = logging.MustGetLogger(LOG_NAME)
	fileList          = list.New()
	confGetLogPathFn  func() string
	confGetLogLevelFn func() string
)

/*
type LogConfig interface {
	GetLogPath() string
	GetLogLevel() string
}*/

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
//     %{id}        Sequence number for log message (uint64).
//     %{pid}       Process id (int)
//     %{time}      Time when log occurred (time.Time)
//     %{level}     Log level (Level)
//     %{module}    Module (string)
//     %{program}   Basename of os.Args[0] (string)
//     %{message}   Message (string)
//     %{longfile}  Full file name and line number: /a/b/c/d.go:23
//     %{shortfile} Final file name element and line number: d.go:23
//     %{color}     ANSI color based on log level
//     %{longpkg}   Full package path, eg. github.com/go-logging
//     %{shortpkg}  Base package path, eg. go-logging
//     %{longfunc}  Full function name, eg. littleEndian.PutUint32
//     %{shortfunc} Base function name, eg. PutUint32
var stdFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfile} >%{level:.5s}%{color:reset} - %{message}",
)

var fileFormat = logging.MustStringFormatter(
	"%{time:15:04:05.000} %{shortfile} >%{level:.5s} - %{message}",
)

var auditFormat = logging.MustStringFormatter(
	"%{time:15:04:05.000} - %{message}",
)

//关闭旧log 打开的文件
//newFile 本次是否打开了新文件
func closeOldLogFd(newFile bool) {
	expectedFdNum := 0
	if newFile {
		expectedFdNum++
	}
	if fileList.Len() > expectedFdNum {
		element := fileList.Front()
		if element == nil {
			return
		}
		if fp, ok := element.Value.(*os.File); ok {
			fileList.Remove(element)
			time.Sleep(time.Second * 5)
			Log.Notice("start close old log file")
			fp.Close()
		} else {
			Log.Error("fd type error")
		}
	}
}

// 如果 path 路径不为空，则使用文件记录日志
// 否则使用std out 输出日志
// change by zzh 20151130
// SetBackend  可重复调用
func initLog(path string, level logging.Level) error {

	if len(path) > 0 {
		fp, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		fileList.PushBack(fp)
		fileBackend := logging.NewLogBackend(fp, "", 1)
		fileFormatter := logging.NewBackendFormatter(fileBackend, fileFormat)

		fileB := logging.AddModuleLevel(fileFormatter)
		fileB.SetLevel(level, "")
		logging.SetBackend(fileB)
	} else {
		stdBackend := logging.NewLogBackend(os.Stdout, "", 1)
		stdFormatter := logging.NewBackendFormatter(stdBackend, stdFormat)
		stdB := logging.AddModuleLevel(stdFormatter)
		stdB.SetLevel(level, "")
		// Set the backends to be used.
		logging.SetBackend(stdB)
	}
	go closeOldLogFd(len(path) > 0)
	return nil
}

func InitLog(
	fn_get_log_path func() string,
	fn_get_log_level func() string,
) {
	confGetLogPathFn = fn_get_log_path
	confGetLogLevelFn = fn_get_log_level
	reloadLog()
}
func reloadLog() {
	logPath := confGetLogPathFn()
	logLevel := confGetLogLevelFn()
	level, err := logging.LogLevel(logLevel)
	if err != nil {
		level = logging.INFO
	}
	if err := initLog(logPath, level); err != nil {
		fmt.Printf("init log  %s error, %v\n", logPath, err)
		os.Exit(1)
	}
	Log = logging.MustGetLogger(LOG_NAME)
}
