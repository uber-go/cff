package cff

import (
	"context"
	"sync"
)

// RunStaticTasks runs a set of statically scheduled tasks.
//
// This function is for INTERNAL USE ONLY. Do not use it directly. There is no
// guarantee of API or behavior compatibility if you use this directly.
func RunStaticTasks(ctx context.Context, schedule [][]func(context.Context) error) error {
	for _, tasks := range schedule {
		if len(tasks) == 0 {
			continue
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		if len(tasks) == 1 {
			if err := tasks[0](ctx); err != nil {
				return err
			}
			continue
		}

		var (
			wg   sync.WaitGroup
			once sync.Once
			err  error
		)

		wg.Add(len(tasks))
		for _, task := range tasks {
			go func(task func(context.Context) error) {
				defer wg.Done()
				if terr := task(ctx); terr != nil {
					once.Do(func() {
						err = terr
					})
				}
			}(task)
		}

		wg.Wait()

		if err != nil {
			return err
		}
	}

	return nil
}
