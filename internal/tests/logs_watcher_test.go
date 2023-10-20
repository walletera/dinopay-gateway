package tests

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

func TestLogsWatcher_WaitFor_LogIsAlreadyThere(t *testing.T) {
    logsWatcher := NewLogsWatcher()
    logsWatcher.Start()
    time.Sleep(1 * time.Millisecond)

    fmt.Println("Hola Mundo Loco!")

    found := logsWatcher.WaitFor("Mundo", 1*time.Second)

    assert.True(t, found)
}

func TestLogsWatcher_WaitFor_LogAppearsAfterTheCallToWaitFor(t *testing.T) {
    logsWatcher := NewLogsWatcher()
    logsWatcher.Start()

    fmt.Println("Hola")

    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("Mundo Loco!")
    }()

    found := logsWatcher.WaitFor("Mundo", 2*time.Second)

    assert.True(t, found)
}

func TestLogsWatcher_WaitFor_LogAppearsTooLate(t *testing.T) {
    logsWatcher := NewLogsWatcher()
    logsWatcher.Start()

    fmt.Println("Hola")

    go func() {
        time.Sleep(2 * time.Second)
        fmt.Println("Mundo Loco!")
    }()

    found := logsWatcher.WaitFor("Mundo", 1*time.Second)

    assert.False(t, found)
}
