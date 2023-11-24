package logs

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/walletera/dinopay-gateway/pkg/logs/testdata"
    "math/rand"
    "sync"
    "testing"
    "time"
)

func TestWatcher_WaitFor_LogIsAlreadyThere(t *testing.T) {
    logsWatcher := NewWatcher()
    err := logsWatcher.Start()
    require.NoError(t, err)

    time.Sleep(1 * time.Millisecond)

    fmt.Println("Hola Mundo Loco!")

    found := logsWatcher.WaitFor("Mundo", 100*time.Millisecond)

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
        time.Sleep(100 * time.Millisecond)
        fmt.Println("Mundo Loco!")
    }()

    found := logsWatcher.WaitFor("Mundo", 200*time.Millisecond)

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
        time.Sleep(200 * time.Millisecond)
        fmt.Println("Mundo Loco!")
    }()

    found := logsWatcher.WaitFor("Mundo", 100*time.Millisecond)
    assert.False(t, found)

    // Let the log appear in the console
    time.Sleep(100 * time.Millisecond)
    err = logsWatcher.Stop()
    require.NoError(t, err)
}

func TestWatcher_WaitFor_MultilineLog(t *testing.T) {
    logsWatcher := NewWatcher()
    err := logsWatcher.Start()
    require.NoError(t, err)

    go func() {
        time.Sleep(100 * time.Millisecond)
        fmt.Print(testdata.MultilineLog)
    }()

    found := logsWatcher.WaitFor("failed creating payment on dinopay", 200*time.Millisecond)

    assert.True(t, found)

    err = logsWatcher.Stop()
    require.NoError(t, err)
}

func TestWatcher_WaitFor_Concurrency(t *testing.T) {
    logsWatcher := NewWatcher()
    err := logsWatcher.Start()
    require.NoError(t, err)

    goroutinesCount := 100

    for i := 0; i < goroutinesCount; i++ {
        go func() {
            time.Sleep(time.Duration(rand.Intn(goroutinesCount)) * time.Millisecond)
            sysLogEntry := testdata.Syslog[rand.Intn(len(testdata.Syslog))]
            fmt.Println(sysLogEntry)
        }()
    }

    wg := &sync.WaitGroup{}
    wg.Add(goroutinesCount)

    for i := 0; i < goroutinesCount; i++ {
        go func() {
            time.Sleep(time.Duration(rand.Intn(goroutinesCount)) * time.Millisecond)
            keyword := testdata.SyslogSubstrs[rand.Intn(len(testdata.SyslogSubstrs))]
            found := logsWatcher.WaitFor(keyword, time.Duration(goroutinesCount+1)*time.Millisecond)
            assert.True(t, found, "keyword not found: ", keyword)
            wg.Done()
        }()
    }

    wg.Wait()

    err = logsWatcher.Stop()
    require.NoError(t, err)
}
