package asyncer

import "sync"

func Multiple(fs ...func(a *any, e *error)) ([]any, []error) {
	var (
		wg   sync.WaitGroup
		res  = make([]any, len(fs))
		errs = make([]error, len(fs))
	)
	for i, f := range fs {
		wg.Add(1)
		go func(i int, f func(a *any, e *error)) {
			defer wg.Done()
			f(&res[i], &errs[i])
		}(i, f)
	}
	wg.Wait()
	return res, errs
}
