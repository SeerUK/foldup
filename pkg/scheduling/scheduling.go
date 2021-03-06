package scheduling

import (
	"time"

	"github.com/gorhill/cronexpr"
)

type expression interface {
	Next(time.Time) time.Time
}

// For testing
var parseExpr = parseCronExpr
var timeNow = time.Now

func parseCronExpr(expr string) (expression, error) {
	return cronexpr.Parse(expr)
}

// ScheduleFunc takes a cron-like expression, and a callback function to execute on a schedule. It
// will run indefinitely, or until there is an error. If there is an error, it will be returned.
// This function is synchronous, and will block.
//
// The `expr` parameter is the cron-like expression (we're using github.com/gorhill/cronexpr).
// The `quit` parameter is a channel that can be sent any int that will break the loop.
// The `fn` parameter is a callback function that will be called each scheduled interval.
func ScheduleFunc(done <-chan int, expr string, fn func() error) error {
	prev := timeNow()

	cexpr, err := parseExpr(expr)
	if err != nil {
		return err
	}
	errs := make(chan error, 1)

	go func() {
		for {
			next := cexpr.Next(prev)
			dur := next.Sub(prev)

			// Keep the last time, we'll use it in the next loop for better accuracy.
			prev = next

			timer := time.NewTimer(dur)

			select {
			case <-timer.C:
				// Call the function, if we error, bail and return it.
				err = fn()
				if err != nil {
					errs <- err
					return
				}
			case <-done:
				errs <- nil
				return
			}
		}
	}()

	return <-errs
}
