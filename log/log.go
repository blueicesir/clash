package log

import (
	"fmt"

	"github.com/blueicesir/clash/common/observable"

	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	logCh  = make(chan interface{})
	source = observable.NewObservable(logCh)
	level  = INFO
)

func init() {
	log.SetLevel(log.DebugLevel)
	// 禁止显示日志的时间，我的命令行窗口不需要显示，否则太长
	log.SetFormatter(&log.TextFormatter{
		DisableColors:true,
		// FullTimestamp:false,
		DisableTimestamp:true,
	})
}

type Event struct {
	LogLevel LogLevel
	Payload  string
}

func (e *Event) Type() string {
	return e.LogLevel.String()
}

func Infoln(format string, v ...interface{}) {
	event := newLog(INFO, format, v...)
	logCh <- event
	print(event)
}

func Warnln(format string, v ...interface{}) {
	event := newLog(WARNING, format, v...)
	logCh <- event
	print(event)
}

func Errorln(format string, v ...interface{}) {
	event := newLog(ERROR, format, v...)
	logCh <- event
	print(event)
}

func Debugln(format string, v ...interface{}) {
	event := newLog(DEBUG, format, v...)
	logCh <- event
	print(event)
}

func Fatalln(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func Subscribe() observable.Subscription {
	sub, _ := source.Subscribe()
	return sub
}

func Level() LogLevel {
	return level
}

func SetLevel(newLevel LogLevel) {
	level = newLevel
}

func print(data *Event) {
	if data.LogLevel < level {
		return
	}

	// add by BlueICE
	if(strings.HasPrefix(data.Payload,"127.0.0.1 -->")){
		data.Payload=strings.TrimSpace(data.Payload[13:])
	}

	switch data.LogLevel {
	case INFO:
		log.Infoln(data.Payload)
	case WARNING:
		log.Warnln(data.Payload)
	case ERROR:
		log.Errorln(data.Payload)
	case DEBUG:
		log.Debugln(data.Payload)
	}
}

func newLog(logLevel LogLevel, format string, v ...interface{}) *Event {
	return &Event{
		LogLevel: logLevel,
		Payload:  fmt.Sprintf(format, v...),
	}
}
