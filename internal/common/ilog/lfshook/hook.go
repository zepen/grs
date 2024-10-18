package lfshook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

// 这里代码示从 https://github.com/rifflock/lfshook/blob/master/lfshook.go copy过来

var defaultFormatter = &logrus.TextFormatter{DisableColors: true}

type PathMap map[logrus.Level]string

type WriterMap map[logrus.Level]io.Writer

type LfsHook struct {
	paths            PathMap
	writers          WriterMap
	levels           []logrus.Level
	lock             *sync.Mutex
	formatter        logrus.Formatter
	defaultPath      string
	defaultWriter    io.Writer
	hasDefaultPath   bool
	hasDefaultWriter bool
}

func NewHook(output interface{}, formatter logrus.Formatter) *LfsHook {
	hook := &LfsHook{
		lock: new(sync.Mutex),
	}

	hook.SetFormatter(formatter)

	switch output.(type) {
	case string:
		hook.SetDefaultPath(output.(string))
		break
	case io.Writer:
		hook.SetDefaultWriter(output.(io.Writer))
		break
	case PathMap:
		hook.paths = output.(PathMap)
		for level := range output.(PathMap) {
			hook.levels = append(hook.levels, level)
		}
		break
	case WriterMap:
		hook.writers = output.(WriterMap)
		for level := range output.(WriterMap) {
			hook.levels = append(hook.levels, level)
		}
		break
	default:
		panic(fmt.Sprintf("unsupported level map type: %v", reflect.TypeOf(output)))
	}

	return hook
}

func (hook *LfsHook) SetFormatter(formatter logrus.Formatter) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if formatter == nil {
		formatter = defaultFormatter
	} else {
		switch formatter.(type) {
		case *logrus.TextFormatter:
			textFormatter := formatter.(*logrus.TextFormatter)
			textFormatter.DisableColors = true
		}
	}

	hook.formatter = formatter
}

func (hook *LfsHook) SetDefaultPath(defaultPath string) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.defaultPath = defaultPath
	hook.hasDefaultPath = true
}

func (hook *LfsHook) SetDefaultWriter(defaultWriter io.Writer) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.defaultWriter = defaultWriter
	hook.hasDefaultWriter = true
}

func (hook *LfsHook) Fire(entry *logrus.Entry) error {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if hook.writers != nil || hook.hasDefaultWriter {
		return hook.ioWrite(entry)
	} else if hook.paths != nil || hook.hasDefaultPath {
		return hook.fileWrite(entry)
	}

	return nil
}

func (hook *LfsHook) ioWrite(entry *logrus.Entry) error {
	var (
		writer io.Writer
		msg    []byte
		err    error
		ok     bool
	)

	if writer, ok = hook.writers[entry.Level]; !ok {
		if hook.hasDefaultWriter {
			writer = hook.defaultWriter
		} else {
			return nil
		}
	}
	// use our formatter instead of entry.String()
	msg, err = hook.formatter.Format(entry)
	if err != nil {
		log.Println("failed to generate string for entry:", err)
		return err
	}
	_, err = writer.Write(msg)
	return err
}

func (hook *LfsHook) fileWrite(entry *logrus.Entry) error {
	var (
		fd   *os.File
		path string
		msg  []byte
		err  error
		ok   bool
	)

	if path, ok = hook.paths[entry.Level]; !ok {
		if hook.hasDefaultPath {
			path = hook.defaultPath
		} else {
			return nil
		}
	}

	dir := filepath.Dir(path)
	os.MkdirAll(dir, os.ModePerm)

	fd, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("failed to open logfile:", path, err)
		return err
	}
	defer fd.Close()

	// use our formatter instead of entry.String()
	msg, err = hook.formatter.Format(entry)

	if err != nil {
		log.Println("failed to generate string for entry:", err)
		return err
	}
	fd.Write(msg)
	return nil
}

// Levels returns configured log levels.
func (hook *LfsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
