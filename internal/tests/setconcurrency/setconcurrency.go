// +build cff

package setconcurrency

import (
	"bufio"
	"bytes"
	"context"
	"runtime"
	"strings"
	"time"

	"go.uber.org/cff"
)

// This must be updated if scheduler.worker is renamed.
const _workerFunction = "go.uber.org/cff/scheduler.worker"

// NumWorkers runs a CFF flow with the provided concurrency, and reports the
// number of workers from within the flow.
func NumWorkers(conc int) (int, error) {
	var numGoroutines int

	err := cff.Flow(
		context.Background(),
		cff.Concurrency(conc),
		cff.Results(&numGoroutines),

		// Workers may run this task while other workers are still
		// spinning up. To work around this, we wait for the number of
		// workers to stabilize before returning.
		cff.Task(func() (int, error) {
			return numWorkersStable(10, time.Millisecond)
		}),
	)

	return numGoroutines, err
}

// numWorkersStable waits for the number of workers reported by numWorkers to
// stabilize for n ticks before reporting it.
func numWorkersStable(n int, tick time.Duration) (int, error) {
	numw, err := numWorkers()
	if err != nil {
		return 0, err
	}

	for remaining := n; remaining > 0; {
		time.Sleep(tick)
		next, err := numWorkers()
		if err != nil {
			return 0, err
		}

		if numw == next {
			remaining--
		} else {
			numw = next
			remaining = n
		}
	}

	return numw, nil
}

// numWorkers reports the number of goroutines currently running the CFF2
// scheduler's worker function.
func numWorkers() (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(getStack()))

	var (
		workers int
		inStack bool
	)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case !inStack && strings.HasPrefix(line, "goroutine "):
			// goroutine 42 [running]:
			inStack = true
		case inStack:
			// path/to/package.function(...)
			if strings.HasPrefix(line, _workerFunction+"(") {
				workers++
			}
		case len(line) == 0:
			inStack = false
		}
	}

	return workers, scanner.Err()
}

// getStack retrieves a stack trace for all running goroutines using
// runtime.Stack.
func getStack() []byte {
	const bufferSize = 64 * 1024 // 64kb

	// runtime.Stack reports the number of bytes actually written to the
	// buffer. If the buffer wasn't large enough, it stops writing. To
	// make sure we have the full stack, we'll double the buffer until we
	// have one large enough to hold the full stack trace.
	for size := bufferSize; ; size *= 2 {
		buf := make([]byte, size)
		if n := runtime.Stack(buf, true); n < size {
			return buf[:n]
		}
	}
}
