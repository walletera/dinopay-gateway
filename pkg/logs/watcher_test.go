package logs

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "testing"
    "time"
)

func TestWatcher_WaitFor_LogIsAlreadyThere(t *testing.T) {
    logsWatcher := NewWatcher()
    err := logsWatcher.Start()
    require.NoError(t, err)

    time.Sleep(1 * time.Millisecond)

    fmt.Println("Hola Mundo Loco!")

    found := logsWatcher.WaitFor("Mundo", 1*time.Second)

    assert.True(t, found)

    err = logsWatcher.Stop()
    require.NoError(t, err)
}

func TestWatcher_WaitFor_LogAppearsAfterTheCallToWaitFor(t *testing.T) {
    logsWatcher := NewWatcher()
    err := logsWatcher.Start()
    require.NoError(t, err)

    fmt.Println("Hola")

    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("Mundo Loco!")
    }()

    found := logsWatcher.WaitFor("Mundo", 2*time.Second)

    assert.True(t, found)

    err = logsWatcher.Stop()
    require.NoError(t, err)
}

func TestWatcher_WaitFor_LogAppearsTooLate(t *testing.T) {
    logsWatcher := NewWatcher()
    err := logsWatcher.Start()
    require.NoError(t, err)

    fmt.Println("Hola")

    go func() {
        time.Sleep(2 * time.Second)
        fmt.Println("Mundo Loco!")
    }()

    found := logsWatcher.WaitFor("Mundo", 1*time.Second)

    assert.False(t, found)

    err = logsWatcher.Stop()
    require.NoError(t, err)
}
