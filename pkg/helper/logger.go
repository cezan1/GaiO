package helper

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var fileLogger *log.Logger

func init() {
	// 创建日志目录和文件
	//当前目录
	curDir, _ := os.Getwd()
	logDir := filepath.Join(curDir, "GaiO_logs")
	_ = os.MkdirAll(logDir, 0755)
	logFile := filepath.Join(logDir, "gaiO_"+time.Now().Format("20060102")+"_log.txt")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		fileLogger = log.New(file, "[FILE LOG] ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// WriteLog 写入日志信息
func WriteLog(message string) {
	Logger := log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Println(message)
}

// WriteErrorLog 写入错误日志信息
func WriteErrorLog(message string) {
	errorLogger := log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger.Println(message)
}

// WriteLogToFile 单独将日志写入文件的方法
func WriteLogToFile(message string) {
	if fileLogger != nil {
		fileLogger.Println(message)
	}
}
