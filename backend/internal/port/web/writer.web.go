package web

import (
	"bytes"
	"clinker-backend/common/logger"
)

type writer struct {
	keys       []string
	processors map[string]func([]byte) any
	id         string
}

func (w *writer) Write(p []byte) (n int, err error) {
	defer func() {
		recover()
	}()

	logItem := logger.Info(ctx)
	bs := bytes.Split(p, sep_bytes)
	for i, b := range bs {
		k := w.keys[i]
		var v any
		if processor, ok := w.processors[k]; ok {
			v = processor(b)
		} else {
			v = string(b)
		}
		logItem = logItem.D(k, v)
	}
	if w.id != "" {
		logItem = logItem.D("user_id", w.id)
	}
	logItem.W()

	return len(p), nil
}
