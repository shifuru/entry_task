package database

import (
	"context"
	"entry_task/common/log"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	//输出前缀
	infoStr  = "%s : [requestId]: %s | [info] "
	warnStr  = "%s : [requestId]: %s | [warn] "
	errStr   = "%s : [requestId]: %s | [error] "
	traceStr = "%s : [requestId]: %s | [%.3fms] [rows:%v] %s"
	//traceWarnStr = "%s %s : [requestId]: %s | [%.3fms] [rows:%v] %s"
	traceErrStr = "%s %s : [requestId]: %s | [%.3fms] [rows:%v] %s"

	//gorm目录前缀，用于输出sql所在行
	gormSourceDir     string
	gormSourceDirOnce sync.Once
	//非用户请求时的requestId
	defaultRequestId = "noRequest"
)

type gormLogger struct {
	logger.Writer
	requestId     *string
	slowThreshold time.Duration
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	ret := *l
	return &ret
}
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Printf(infoStr+msg, append([]interface{}{fileWithLineNum(), *l.requestId}, data...)...)
}
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Printf(warnStr+msg, append([]interface{}{fileWithLineNum(), *l.requestId}, data...)...)
}
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Printf(errStr+msg, append([]interface{}{fileWithLineNum(), *l.requestId}, data...)...)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(traceErrStr, fileWithLineNum(), *l.requestId, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(traceErrStr, fileWithLineNum(), *l.requestId, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(traceStr, fileWithLineNum(), *l.requestId, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(traceStr, fileWithLineNum(), *l.requestId, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		if l.slowThreshold > 0 && elapsed > time.Millisecond*l.slowThreshold {
			log.Error(l.requestId, "slow sql error: "+traceStr, fileWithLineNum(), *l.requestId, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		}
	}
}

/**
 * @author
 * @description 获取调用行
 * @date
 * @param
 * @return
 **/
func fileWithLineNum() string {
	//只有gorm会调用到这里，获取所在目录
	gormSourceDirOnce.Do(func() {
		for i := 1; i < 15; i++ {
			_, file, _, ok := runtime.Caller(i)
			if ok && strings.Contains(file, "gorm") {
				gormSourceDir = file[0:strings.LastIndex(file, "/")]
				return
			}
		}
		//获取失败，写空
		gormSourceDir = ""
	})

	// 只有trace中调用，获取trace外层，直接从2开始
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, gormSourceDir) || strings.HasSuffix(file, "_test.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}

/**
 * @author
 * @description 获取gorm的log
 * @date
 * @param
 * @return
 **/
func newGormLog(requestId *string, slowThreshold *int64) *gormLogger {
	//var l dataLogger
	//l.Interface = logger.New(mylog.GetDataLog(), configs)
	//l.requestId = requestId
	//return &l
	var l gormLogger
	l.Writer = log.GetDataLog()
	if nil != requestId {
		l.requestId = requestId
	} else {
		l.requestId = &defaultRequestId
	}
	//慢sql默认为50ms
	if nil == slowThreshold || *slowThreshold <= 0 {
		l.slowThreshold = time.Duration(50)
	} else {
		l.slowThreshold = time.Duration(*slowThreshold)
	}
	return &l
}
