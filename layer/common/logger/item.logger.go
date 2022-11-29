package logger

type LogItem interface {
	// data
	D(key string, data interface{}) LogItem

	// error
	E(err error) LogItem

	// write
	W(messages ...string)

	// writef
	Wf(message string, a ...interface{})
}
