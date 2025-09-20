# Logs Watcher
[![Go](https://github.com/walletera/logs-watcher/actions/workflows/go.yml/badge.svg)](https://github.com/walletera/logs-watcher/actions/workflows/go.yml)

**logs-watcher** is a Go library designed to help developers test and assert that their application produces the expected logs. This is particularly useful in integration tests, where it's important to verify that certain actions generate the correct log output, either through `os.Stdout`/`os.Stderr` or through structured loggers like `slog`.

## Features

- Watch logs from `Stdout` and `Stderr` during test execution.
- Watch logs produced by `slog` structured loggers.
- Wait for specific log messages or patterns to appear a set number of times.
- Set timeouts for log expectations to make your tests robust.
- Thread-safe and suitable for concurrent test scenarios.
- Restore original output streams after tests.

## Installation
```
go get github.com/walletera/logs-watcher
``` 

## Usage

### Watching Stdout/Stderr logs

The following example demonstrates how to use the library to watch for a specific message written to `os.Stdout` or `os.Stderr` in an integration test:
```
go package yourpackage
import ( "fmt" "testing" "time"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
"github.com/walletera/logs-watcher/std"
)
func Test_StdoutLogAppears(t *testing.T) { watcher, err := std.NewWatcher() require.NoError(t, err) defer watcher.Stop()
fmt.Println("Hello Integration Test!")

found := watcher.WaitFor("Integration", 100*time.Millisecond)
assert.True(t, found, "Expected log message was not found")
}
``` 

### Watching `slog` Logs

You can also use `logs-watcher` to assert logs produced by `slog` loggers:
```
go package yourpackage
import ( "log/slog" "testing" "time"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
"github.com/walletera/logs-watcher/slog"
)
func Test_SlogMessageAppears(t *testing.T) { handler := slog.NewTextHandler() // Or any slog.Handler you are using logsWatcher := slog.NewWatcher(handler) defer logsWatcher.Stop()
logger := slog.New(logsWatcher.DecoratedHandler())

logger.Info("Expected message in logs")

found := logsWatcher.WaitFor("Expected message", 100*time.Millisecond)
assert.True(t, found, "Expected slog message was not found")
}
``` 

### Waiting for a message multiple times
```
go found := watcher.WaitForNTimes("keyword", 200*time.Millisecond, 3) // Wait for "keyword" 3 times assert.True(t, found)
``` 

---

## API

- `WaitFor(keyword string, timeout time.Duration) bool`: Waits for the keyword to appear in the logs within the given timeout.
- `WaitForNTimes(keyword string, timeout time.Duration, n int) bool`: Waits for the keyword to appear `n` times within the timeout.
- `Stop() error`: Stops watching and restores original output streams (if applicable).
```
