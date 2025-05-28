package log

import (
	"entry_task/common/utils"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	accessLog Logger
	errorLog  Logger
	infoLog   Logger
	dataLog   Logger
)

const (
	logSizeMax = 9999
	logDir     = "./logs"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var (
	//error文件日志前缀
	errorFilePrefix = "ERROR :"
	//error屏幕日志前缀
	errorScreenPrefix = red + errorFilePrefix + reset
	//info文件日志前缀
	infoFilePrefix = "INFO :"
	//info屏幕日志前缀
	infoScreenPrefix = yellow + infoFilePrefix + reset
	//access文件日志前缀
	accessFilePrefix = "ACCESS :"
	//access屏幕日志前缀
	accessScreenPrefix = green + accessFilePrefix + reset
	//data文件日志前缀
	dataFilePrefix = "DATA :"
	//access屏幕日志前缀
	dataScreenPrefix = blue + dataFilePrefix + reset

	//log输出flag
	LogFlag     = log.Llongfile | log.Ldate | log.Lmicroseconds
	DataLogFlag = log.Ldate | log.Lmicroseconds
)

type Logger struct {
	*log.Logger
	toFile *log.Logger
	buffer io.Writer
}

func (l *Logger) Printf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if len(s) > logSizeMax {
		s = s[0:logSizeMax]
	}

	if nil != l.toFile {
		err := l.toFile.Output(3, s)
		if nil != err {
			errorLog.Logger.Printf("toFileError: %s", err.Error())
		}
	}
}

func init() {
	//error日志配置
	{
		errorLog.Logger = log.New(os.Stderr, errorScreenPrefix, LogFlag)
		errorLogFileUri := logDir + "/error.log"
		file, err := utils.CreateOrOpenFile(errorLogFileUri)
		file, err = utils.CreateOrOpenFile(errorLogFileUri)
		if nil != err {
			defer errorLog.Logger.Printf("log init error, can't open %s, err: %s", errorLogFileUri, err.Error())
		} else {
			errorLog.buffer = file
			errorLog.toFile = log.New(errorLog.buffer, errorFilePrefix, LogFlag)
		}
	}
	//info日志配置
	{
		infoLog.Logger = log.New(os.Stderr, infoScreenPrefix, LogFlag)
		infoLogFileUri := logDir + "/info.log"
		file, err := utils.CreateOrOpenFile(infoLogFileUri)
		if nil != err {
			defer errorLog.Logger.Printf("log init error, can't open %s, err: %s", infoLogFileUri, err.Error())
		} else {
			infoLog.buffer = file
			infoLog.toFile = log.New(infoLog.buffer, infoFilePrefix, LogFlag)
		}
	}
	//access日志配置
	{
		accessLog.Logger = log.New(os.Stderr, accessScreenPrefix, LogFlag)
		accessLogFileUri := logDir + "/access.log"
		file, err := utils.CreateOrOpenFile(accessLogFileUri)
		if nil != err {
			defer errorLog.Logger.Printf("log init error, can't open %s, err: %s", accessLogFileUri, err.Error())
		} else {
			accessLog.buffer = file
			accessLog.toFile = log.New(accessLog.buffer, accessFilePrefix, LogFlag)
		}
	}
	//data日志配置
	{
		dataLog.Logger = log.New(os.Stderr, dataScreenPrefix, DataLogFlag)
		dataLogFileUri := logDir + "/data.log"
		file, err := utils.CreateOrOpenFile(dataLogFileUri)
		if nil != err {
			defer errorLog.Logger.Printf("log init error, can't open %s, err: %s", dataLogFileUri, err.Error())
		} else {
			dataLog.buffer = file
			dataLog.toFile = log.New(dataLog.buffer, dataFilePrefix, DataLogFlag)
		}
	}
}

func Error(requestId *string, format string, v ...interface{}) {
	if requestId == nil {
		errorLog.Printf(format, v...)
	} else {
		format = "request_id=" + *requestId + " | " + format
		errorLog.Printf(format, v...)
	}
}

func Info(requestId *string, format string, v ...interface{}) {
	if requestId == nil {
		infoLog.Printf(format, v...)
	} else {
		format = "request_id=" + *requestId + " | " + format
		infoLog.Printf(format, v...)
	}
}

func Access(requestId *string, format string, v ...interface{}) {
	if requestId == nil {
		accessLog.Printf(format, v...)
	} else {
		format = "request_id=" + *requestId + " | " + format
		accessLog.Printf(format, v...)
	}
}

func GetDataLog() *Logger {
	return &dataLog
}
