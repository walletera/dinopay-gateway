package logs

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)

type Watcher struct {
    realStdout                      *os.File
    realStderr                      *os.File
    pipeRead                        *os.File
    pipeWrite                       *os.File
    lines                           []string
    linesMutex                      sync.RWMutex
    newLinesCh                      chan string
    newLinesSubscriptionIdGenerator atomic.Int64
    newLinesSubscriptionCh          chan newLinesSubscription
    deleteNewLinesSubscriptionCh    chan newLinesSubscription
    stop                            chan bool
}

type newLinesSubscription struct {
    id         int64
    newLinesCh chan string
}

func NewWatcher() *Watcher {

    return &Watcher{
        lines:                        make([]string, 0),
        newLinesCh:                   make(chan string),
        newLinesSubscriptionCh:       make(chan newLinesSubscription),
        deleteNewLinesSubscriptionCh: make(chan newLinesSubscription),
        stop:                         make(chan bool, 1),
    }
}

func (l *Watcher) Start() error {
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

    go l.startControlLoop()
    go l.startScanner(r)

    return nil
}

func (l *Watcher) startControlLoop() {
    newLinesSubscriptions := make(map[int64]chan string)
    for {
        select {
        case newLine := <-l.newLinesCh:
            l.storeNewLine(newLine)
            l.broadcastNewLine(newLine, newLinesSubscriptions)
        case subscription := <-l.newLinesSubscriptionCh:
            newLinesSubscriptions[subscription.id] = subscription.newLinesCh
        case subscription := <-l.deleteNewLinesSubscriptionCh:
            delete(newLinesSubscriptions, subscription.id)
        case <-l.stop:
            return
        }
    }
}

func (l *Watcher) startScanner(r io.Reader) {
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        line := scanner.Text()
        l.newLinesCh <- line
    }
    if err := scanner.Err(); err != nil {
        select {
        case <-l.stop:
            return
        default:
        }
        os.Stdout = l.realStdout
        os.Stderr = l.realStderr
        fmt.Println("logs watcher failed reading standard output: ", err.Error())
    }
}

func (l *Watcher) storeNewLine(line string) {
    l.linesMutex.Lock()
    defer l.linesMutex.Unlock()
    l.lines = append(l.lines, line)
}

func (l *Watcher) broadcastNewLine(line string, newLinesSubscriptions map[int64]chan string) {
    for _, newLinesCh := range newLinesSubscriptions {
        newLinesCh <- line
    }
}

func (l *Watcher) Stop() error {
    close(l.stop)

    os.Stdout = l.realStdout
    os.Stderr = l.realStderr

    err := l.pipeRead.Close()
    if err != nil {
        return fmt.Errorf("error stop logs watcher: %w", err)
    }
    err = l.pipeWrite.Close()
    if err != nil {
        return fmt.Errorf("error stop logs watcher: %w", err)
    }
    return nil
}

func (l *Watcher) WaitFor(keyword string, timeout time.Duration) bool {
    newLinesCh := make(chan string)
    subscription := l.subscribeForNewLines(newLinesCh)
    foundLineCh := make(chan bool, 2)
    go l.searchInNewLines(keyword, newLinesCh, foundLineCh)
    go l.searchInStoredLines(keyword, foundLineCh)
    found := l.wait(foundLineCh, timeout)
    l.unsubscribeFromNewLines(subscription)
    close(newLinesCh)
    return found
}

func (l *Watcher) subscribeForNewLines(newLinesCh chan string) newLinesSubscription {
    subscriptionId := l.newLinesSubscriptionIdGenerator.Add(1)
    subscription := newLinesSubscription{
        id:         subscriptionId,
        newLinesCh: newLinesCh,
    }
    l.newLinesSubscriptionCh <- subscription
    return subscription
}

func (l *Watcher) unsubscribeFromNewLines(subscription newLinesSubscription) {
    l.deleteNewLinesSubscriptionCh <- subscription
}

func (l *Watcher) searchInNewLines(keyword string, newLinesCh chan string, foundLineCh chan bool) {
    for newLine := range newLinesCh {
        if strings.Contains(newLine, keyword) {
            foundLineCh <- true
        }
    }
}

func (l *Watcher) searchInStoredLines(keyword string, foundLine chan bool) {
    l.linesMutex.RLock()
    defer l.linesMutex.RUnlock()
    for _, line := range l.lines {
        if strings.Contains(line, keyword) {
            foundLine <- true
            return
        }
    }
}

func (l *Watcher) wait(foundLineCh chan bool, timeout time.Duration) bool {
    timeoutCh := time.After(timeout)
    found := false
    select {
    case <-foundLineCh:
        found = true
    case <-timeoutCh:
    }
    return found
}
