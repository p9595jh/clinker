package util

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
)

const DateFormat = "2006-01-02T15:04:05.999999-07:00"

// strings.Repeat("0", 64)
const DefaultTxHash = "0000000000000000000000000000000000000000000000000000000000000000"

// strings.Repeat("0", 40)
const DefaultAddress = "0000000000000000000000000000000000000000"

func RandHex(l int) string {
	rand.Seed(time.Now().Unix())
	var buf bytes.Buffer
	for i := 0; i < l; i++ {
		buf.WriteString(strconv.FormatInt(int64(rand.Intn(16)), 16))
	}
	return buf.String()
}
