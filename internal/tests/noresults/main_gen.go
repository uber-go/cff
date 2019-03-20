// +build !cff

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
)

func main() {
	h := &h{}
	ctx := context.Background()
	err := h.swallow(ctx, os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(os.Args, "is swallowed")
}

type h struct{}

func (h *h) swallow(ctx context.Context, req string) (err error) {
	err = func(ctx context.Context, v1 string) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(2)

		var err0 error
		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			err0 = func(s string) error {
				if s == "tide pods" {
					return errors.New("can not swallow")
				}
				return nil
			}(v1)
			if err0 != nil {

				once0.Do(func() {
					err = err0
				})
			}

		}()

		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			func(string) {}(v1)

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
		)

		return err
	}(ctx, req)
	return
}

func (h *h) tripleSwallow(ctx context.Context, req string) (err error) {
	err = func(ctx context.Context, v1 string) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(3)

		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			func(string) {}(v1)

		}()

		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			func(string) {}(v1)

		}()

		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			func(string) {}(v1)

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
		)

		return err
	}(ctx, req)
	return
}
