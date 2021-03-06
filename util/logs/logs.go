package logs

import (
	"context"
	"fmt"
	"log"
	"os"
)

type logLevel int

const (
	Debug logLevel = iota
	Info
	Warn
	Error
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	level       logLevel
)

func init() {
	debugLogger = log.New(os.Stdout, "DEBUG:", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	infoLogger = log.New(os.Stdout, "INFO:", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	warnLogger = log.New(os.Stderr, "WARN:", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	errorLogger = log.New(os.Stderr, "ERROR:", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	level = Debug
}

// 设置全局日志级别
func SetGlobalLogLevel(l logLevel) {
	level = l
}

// 生成单次请求唯一后缀
func getSuffix(ctx context.Context) string {
	ret := "["
	reqId, ok := ctx.Value("req_id").(string)
	if ok {
		ret += "reqId=" + reqId + " "
	}
	email, ok := ctx.Value("email").(string)
	if ok {
		ret += "email=" + email + " "
	}
	return ret + "]"
}

func CtxDebug(ctx context.Context, format string, v ...interface{}) {
	if level > Debug {
		return
	}
	msg := fmt.Sprintf(format, v...)
	msg += getSuffix(ctx)
	_ = debugLogger.Output(2, msg)
}

func CtxInfo(ctx context.Context, format string, v ...interface{}) {
	if level > Info {
		return
	}
	msg := fmt.Sprintf(format, v...)
	msg += getSuffix(ctx)
	_ = infoLogger.Output(2, msg)
}

func CtxWarn(ctx context.Context, format string, v ...interface{}) {
	if level > Warn {
		return
	}
	msg := fmt.Sprintf(format, v...)
	msg += getSuffix(ctx)
	_ = warnLogger.Output(2, msg)
}

func CtxError(ctx context.Context, format string, v ...interface{}) {
	if level > Error {
		return
	}
	msg := fmt.Sprintf(format, v...)
	msg += getSuffix(ctx)
	_ = errorLogger.Output(2, msg)
}
