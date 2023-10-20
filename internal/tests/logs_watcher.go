package tests

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "sync"
    "time"
)

type LogsWatcher struct {
    realStdout  *os.File
    lines       []string
    newLinesChs []chan string

    mutex sync.Mutex
}

func NewLogsWatcher() *LogsWatcher {
    return &LogsWatcher{
        lines:       make([]string, 0),
        newLinesChs: make([]chan string, 0),
    }
}

func (l *LogsWatcher) Start() {
    l.realStdout = os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    os.Stderr = w

    go func() {
        scanner := bufio.NewScanner(r)
        for scanner.Scan() {
            line := scanner.Text()
            fmt.Fprintln(l.realStdout, line)
            l.mutex.Lock()
            l.lines = append(l.lines, line)
            for _, newLinesCh := range l.newLinesChs {
                select {
                case newLinesCh <- line:
                default:
                }
            }
            l.mutex.Unlock()
        }
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "error reading standard output:", err)
        }
    }()
}

func (l *LogsWatcher) WaitFor(keyword string, timeout time.Duration) bool {
    newLinesCh := make(chan string)
    l.mutex.Lock()
    l.newLinesChs = append(l.newLinesChs, newLinesCh)
    l.mutex.Unlock()
    foundLine := make(chan bool)
    go func() {
        for {
            newLine := <-newLinesCh
            if strings.Contains(newLine, keyword) {
                foundLine <- true
            }
        }
    }()
    go func() {
        l.mutex.Lock()
        for _, line := range l.lines {
            if strings.Contains(line, keyword) {
                foundLine <- true
            }
        }
        l.mutex.Unlock()
    }()
    timeoutCh := time.After(timeout)
    select {
    case <-foundLine:
        return true
    case <-timeoutCh:
        return false
    }
}
