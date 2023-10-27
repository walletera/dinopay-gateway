package logs

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "sync"
    "time"
)

type Watcher struct {
    realStdout  *os.File
    realStderr  *os.File
    pipeRead    *os.File
    pipeWrite   *os.File
    lines       []string
    newLinesChs []chan string

    stopping chan bool
    mutex    sync.Mutex
}

func NewWatcher() *Watcher {
    return &Watcher{
        lines:       make([]string, 0),
        newLinesChs: make([]chan string, 0),
        stopping:    make(chan bool, 1),
    }
}

func (l *Watcher) Start() error {
    fmt.Printf("logs watcher started")

    l.realStdout = os.Stdout
    l.realStderr = os.Stderr

    r, w, err := os.Pipe()
    if err != nil {
        return fmt.Errorf("failed starting logs watcher: %w", err)
    }

    l.pipeRead = r
    l.pipeWrite = w
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
            select {
            case <-l.stopping:
                fmt.Println("logs watcher stopped")
                return
            default:
            }
            os.Stdout = l.realStdout
            os.Stderr = l.realStderr
            fmt.Println("logs watcher failed reading standard output: ", err.Error())
        }
    }()

    return nil
}

func (l *Watcher) Stop() error {
    l.stopping <- true

    os.Stdout = l.realStdout
    os.Stderr = l.realStderr

    err := l.pipeRead.Close()
    if err != nil {
        return fmt.Errorf("error stopping logs watcher: %w", err)
    }
    err = l.pipeWrite.Close()
    if err != nil {
        return fmt.Errorf("error stopping logs watcher: %w", err)
    }
    return nil
}

func (l *Watcher) WaitFor(keyword string, timeout time.Duration) bool {
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
