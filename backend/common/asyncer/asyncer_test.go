package asyncer_test

import (
	asyncer "clinker-backend/common/asyncer"
	"testing"
	"time"
)

func TestMultipleAsync(t *testing.T) {
	f1 := func() string {
		time.Sleep(time.Second * 3)
		return "hello"
	}
	f2 := func(i int) int {
		time.Sleep(time.Second * 2)
		return 123 * i
	}
	f3 := func() int {
		time.Sleep(time.Second * 1)
		return 22
	}

	res, _ := asyncer.Multiple(
		func(a *any, e *error) {
			*a = f1()
		},
		func(a *any, e *error) {
			*a = f2(10)
		},
		func(a *any, e *error) {
			*a = f3()
		},
	)
	t.Log(res)

	var (
		f1r = res[0].(string)
		f2r = res[1].(int)
		f3r = res[2].(int)
	)
	t.Log(f1r, f2r, f3r)
}
