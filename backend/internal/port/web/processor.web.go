package web

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

var skipLogging = map[string]bool{
	"/api/health":                     true,
	"/favicon.ico":                    true,
	"/swagger":                        true,
	"doc.json":                        true,
	"index.html":                      true,
	"swagger-ui.css":                  true,
	"swagger-ui-standalone-preset.js": true,
	"swagger-ui-bundle.js":            true,
}

type processor struct{}

func (*processor) responseBody(b []byte) any {
	var i any
	if err := json.Unmarshal(b, &i); err != nil {
		return string(b)
	} else {
		return i
	}
}

func (p *processor) requestBody(b []byte) any {
	return p.responseBody(b)
}

func (*processor) requestParams(b []byte) any {
	i := map[string]any{}
	bytes2d := bytes.Split(b, []byte{38})
	for _, kvBytes := range bytes2d {
		kv := bytes.Split(kvBytes, []byte{61})
		if (len(kv)) == 2 {
			i[string(kv[0])] = string(kv[1])
		}
	}
	return i
}

func (*processor) responseTime(b []byte) any {
	return strings.Trim(string(b), " ")
}

func (*processor) statusCode(b []byte) any {
	s := string(b)
	if i, err := strconv.ParseUint(s, 10, 16); err != nil {
		return s
	} else {
		return i
	}
}

func (*processor) requestUri(b []byte) any {
	s := strings.TrimRight(string(b), "/")
	if skipLogging[s] {
		panic(0)
	} else {
		return s
	}
}
