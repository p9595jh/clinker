package logger

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	consoleLogger, fileLogger zerolog.Logger
	format                    = "./log/2006-01-02.log"
)

type logItem struct {
	c, f *zerolog.Event
}

func init() {
	consoleLogger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
}

func FileLoggerInit() {
	if _, err := os.Stat("./log"); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir("./log", 0755); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	var (
		mutex           = &sync.Mutex{}
		date            = time.Now().Format(format)
		writer *os.File = nil
	)

	fileLoggerInit := func(date string) {
		file, err := os.OpenFile(date, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		fileLogger = zerolog.New(file).With().Timestamp().Logger()
		if writer != nil {
			writer.Close()
		}
		writer = file
		mutex.Unlock()
	}
	fileLoggerInit(date)

	go func(date string) {
		for t := range time.NewTimer(time.Second * 30).C {
			newDate := t.Format(format)
			if date != newDate {
				date = newDate
				fileLoggerInit(newDate)
			}
		}
	}(date)
}

func newLogItem(ctx []string, f, c *zerolog.Event) LogItem {
	s := strings.Join(ctx, "::")
	return &logItem{
		c: c.Str("context", s),
		f: f.Str("context", s),
	}
}

func Info(ctx ...string) LogItem {
	return newLogItem(ctx, consoleLogger.Info(), fileLogger.Info())
}

func Debug(ctx ...string) LogItem {
	return newLogItem(ctx, consoleLogger.Debug(), fileLogger.Debug())
}

func Warn(ctx ...string) LogItem {
	return newLogItem(ctx, consoleLogger.Warn(), fileLogger.Warn())
}

func Error(ctx ...string) LogItem {
	return newLogItem(ctx, consoleLogger.Error(), fileLogger.Error())
}

// data
func (l *logItem) D(key string, data interface{}) LogItem {
	l.f = l.f.Interface(key, data)
	l.c = l.c.Interface(key, data)
	return l
}

// error
func (l *logItem) E(err error) LogItem {
	l.c = l.c.Err(err)
	l.f = l.f.Err(err)
	return l
}

// write
func (l *logItem) W(messages ...string) {
	if len(messages) == 0 {
		l.c.Send()
		l.f.Send()
	} else {
		msg := strings.Join(messages, " ")
		l.c.Msg(msg)
		l.f.Msg(msg)
	}
}

// write format
func (l *logItem) Wf(message string, a ...interface{}) {
	l.c.Msgf(message, a...)
	l.f.Msgf(message, a...)
}
