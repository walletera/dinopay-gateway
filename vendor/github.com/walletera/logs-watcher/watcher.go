package logs

import "time"

type IWatcher interface {
    WaitForNTimes(keyword string, timeout time.Duration, n int) bool
    WaitFor(keyword string, timeout time.Duration) bool
    Stop() error
}
