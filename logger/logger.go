package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

type MyFormatter struct{}

var levelList = []string{
	"PANIC",
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

func (mf *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	level := levelList[int(entry.Level)]
	strList := strings.Split(entry.Caller.File, "/")
	fileName := strList[len(strList)-1]
	b.WriteString(fmt.Sprintf("%s %s %s:%d | %s\n",
		entry.Time.Format("2006-01-02 15:04:05.678"), level, fileName,
		entry.Caller.Line, entry.Message))
	return b.Bytes(), nil
}

func MakeLogger(filename string, display bool) *logrus.Logger {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err.Error())
	}
	logger := logrus.New()
	if display {
		logger.SetOutput(io.MultiWriter(os.Stdout, f))
	} else {
		logger.SetOutput(io.MultiWriter(f))
	}
	logger.SetReportCaller(true)
	logger.SetFormatter(&MyFormatter{})
	return logger
}

func LoggerInit(lvl string, filename string, display bool) {
	Log = MakeLogger(filename, display)
	level, _ := logrus.ParseLevel(lvl)
	Log.SetLevel(level)
}
