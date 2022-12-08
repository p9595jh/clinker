package asyncer

import (
	"reflect"
	"sync"
)

func Multiple(fs ...func(good, bad *any)) ([]any, []any) {
	var (
		wg   sync.WaitGroup
		ress = make([]any, len(fs))
		errs = make([]any, len(fs))
	)
	for i, f := range fs {
		wg.Add(1)
		go func(i int, f func(good, bad *any)) {
			defer wg.Done()
			f(&ress[i], &errs[i])
			if errs[i] != nil {
				if reflect.ValueOf(errs[i]).IsNil() {
					errs[i] = nil
				}
			}
		}(i, f)
	}
	wg.Wait()
	return ress, errs
}

func Race(fs ...func(good, bad *any)) (any, any) {
	var (
		ch   = make(chan int)
		ress = make([]any, len(fs))
		errs = make([]any, len(fs))
	)
	for i, f := range fs {
		go func(i int, f func(good, bad *any)) {
			f(&ress[i], &errs[i])
			ch <- i
		}(i, f)
	}
	idx := <-ch
	return ress[idx], errs[idx]
}
