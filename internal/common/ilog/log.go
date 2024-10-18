package ilog

import (
	"context"
	"github.com/opensearch-project/opensearch-go"
)

type Logger interface {
	Info(msg string)

	Warn(msg string)

	Error(msg string)

	Alert(msg string)

	Infof(msg string, args ...interface{})

	Warnf(msg string, args ...interface{})

	Errorf(msg string, args ...interface{})

	Alertf(msg string, args ...interface{})

	FromCtx(ctx context.Context) Logger
}

var Log Logger

type LogConf struct {
	Typo     string `yaml:"typo"`
	FilePath string `yaml:"file_path"`
}

func InitLog(conf *LogConf, args ...interface{}) {
	if conf.Typo == "openSearch" {
		client := args[0].(*opensearch.Client)
		index := args[1].(string)
		InitOSCHLog(client, index)
	} else if conf.Typo == "localFile" {
		InitLocalFileLog(conf.FilePath)
	} else {
		InitFmtLog()
	}
}
