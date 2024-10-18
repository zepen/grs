package ilog

import (
	"testing"
)

func TestFmtLog(t *testing.T) {
	lc := &LogConf{}
	lc.Typo = "std"
	InitLog(lc)
	Log.Infof("hello!")
	Log.Warnf("what?")
	Log.Errorf("No!")
}

func TestLogrusLog(t *testing.T) {
	lc := &LogConf{FilePath: "app.log"}
	lc.Typo = "localFile"
	InitLog(lc)
	Log.Infof("hello!")
	Log.Warnf("what?")
	Log.Errorf("No!")
}
