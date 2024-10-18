package ilog

import (
	"context"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/opensearch-project/opensearch-go"
	"github.com/sirupsen/logrus"
	"recommend-server/internal/common/ilog/eshook"
	"recommend-server/internal/common/ilog/lfshook"
	"time"
)

type LogrusLog struct {
	Logger *logrus.Logger
}

func (lrsl *LogrusLog) Info(msg string) {
	lrsl.Logger.Info(msg)
}

func (lrsl *LogrusLog) Warn(msg string) {
	lrsl.Logger.Warn(msg)
}

func (lrsl *LogrusLog) Error(msg string) {
	lrsl.Logger.Error(msg)
}

func (lrsl *LogrusLog) Alert(msg string) {

}

func (lrsl *LogrusLog) Infof(msg string, args ...interface{}) {
	lrsl.Logger.Infof(msg, args...)
}

func (lrsl *LogrusLog) Warnf(msg string, args ...interface{}) {
	lrsl.Logger.Warnf(msg, args...)
}

func (lrsl *LogrusLog) Errorf(msg string, args ...interface{}) {
	lrsl.Logger.Errorf(msg, args...)
}

func (lrsl *LogrusLog) Alertf(msg string, args ...interface{}) {

}

func (lrsl *LogrusLog) FromCtx(ctx context.Context) Logger {
	return lrsl
}

func InitOSCHLog(client *opensearch.Client, logIndex string) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Hooks.Add(&eshook.OpenSearchHook{Client: client, Index: logIndex})
	Log = &LogrusLog{Logger: logger}
}

func InitLocalFileLog(filePath string) {
	logger := logrus.New()
	writer, _ := rotatelogs.New(
		filePath+".%Y-%m-%d",
		rotatelogs.WithLinkName(filePath),
		rotatelogs.WithMaxAge(time.Duration(604800)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.FatalLevel: writer,
		logrus.DebugLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.PanicLevel: writer,
	}
	// logger.SetReportCaller(true)
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	logger.AddHook(lfHook)
	Log = &LogrusLog{Logger: logger}
}
