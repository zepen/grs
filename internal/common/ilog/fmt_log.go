package ilog

import (
	"context"
	"fmt"
)

func InitFmtLog() {
	Log = &FmtLog{}
}

type FmtLog struct {
}

func (fl *FmtLog) Info(msg string) {
	fmt.Println("info:" + msg)
}

func (fl *FmtLog) Warn(msg string) {
	fmt.Println("warn:" + msg)
}

func (fl *FmtLog) Error(msg string) {
	fmt.Println("err:" + msg)
}

func (fl *FmtLog) Alert(msg string) {
	fmt.Println("alert:" + msg)
}

func (fl *FmtLog) Infof(msg string, args ...interface{}) {
	fmt.Println("info:" + fmt.Sprintf(msg, args...))
}

func (fl *FmtLog) Warnf(msg string, args ...interface{}) {
	fmt.Println("warn:" + fmt.Sprintf(msg, args...))
}

func (fl *FmtLog) Errorf(msg string, args ...interface{}) {
	fmt.Println("error:" + fmt.Sprintf(msg, args...))
}

func (fl *FmtLog) Alertf(msg string, args ...interface{}) {
	fmt.Println("alert:" + fmt.Sprintf(msg, args...))
}

func (fl *FmtLog) FromCtx(ctx context.Context) Logger {
	return fl
}
