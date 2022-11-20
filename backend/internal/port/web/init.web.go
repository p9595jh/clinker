package web

import (
	"strings"

	"github.com/p9595jh/fpgo"
)

var (
	reqLogFormat string

	sep_string = "\r\n"
	sep_bytes  = []byte(sep_string)
	ctx        = "HTTP"
)

func init() {
	// https://docs.gofiber.io/api/middleware/logger
	tags := fpgo.Pipe[[]string, []string](
		[]string{"status", "method", "path", "queryParams", "body", "resBody", "latency", "ip"},
		fpgo.Map(func(i int, s *string) *string {
			*s = "${" + *s + "}"
			return s
		}),
	)
	reqLogFormat = strings.Join(tags, sep_string)
}
