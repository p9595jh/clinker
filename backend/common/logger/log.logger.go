package logger

import (
	"strings"

	"github.com/rs/zerolog"
)

type logItem struct {
	e *zerolog.Event
}

var logger *zerolog.Logger = nil

func Inject(_logger *zerolog.Logger) {
	logger = _logger
}

func newLogItem(ctx []string, e *zerolog.Event) LogItem {
	return &logItem{
		e: e.Str("context", strings.Join(ctx, "::")),
	}
}

func Info(ctx ...string) LogItem {
	return newLogItem(ctx, logger.Info())
}

func Debug(ctx ...string) LogItem {
	return newLogItem(ctx, logger.Debug())
}

func Warn(ctx ...string) LogItem {
	return newLogItem(ctx, logger.Warn())
}

func Error(ctx ...string) LogItem {
	return newLogItem(ctx, logger.Error())
}

// data
func (l *logItem) D(key string, data interface{}) LogItem {
	l.e = l.e.Interface(key, data)
	return l
}

// error
func (l *logItem) E(err error) LogItem {
	l.e = l.e.Err(err)
	return l
}

// write
func (l *logItem) W(messages ...string) {
	if len(messages) == 0 {
		l.e.Send()
	} else {
		l.e.Msg(strings.Join(messages, " "))
	}
}

func (l *logItem) Wf(message string, a ...interface{}) {
	l.e.Msgf(message, a...)
}
