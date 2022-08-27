package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger(level string) Logger {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Panicf("not valid logrus debug level: %s", level)
	}
	l := Logger{e}
	l.Logger.SetLevel(lvl)
	return l
}

func init() {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: true,
		FullTimestamp: true,
	})
	l.SetReportCaller(true)

	e = logrus.NewEntry(l)
}
