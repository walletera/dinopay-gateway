package newline

import (
    "context"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)

type Watcher struct {
    lines                           []string
    linesMutex                      sync.RWMutex
    newLinesCh                      chan string
    newLinesSubscriptionIdGenerator atomic.Int64
    newLinesSubscriptionCh          chan *newLinesSubscription
    deleteNewLinesSubscriptionCh    chan *newLinesSubscription
    stop                            chan bool
}

type newLinesSubscription struct {
    id         int64
    newLinesCh chan string
    active     bool
}

func NewWatcher() *Watcher {

    w := &Watcher{
        lines:                        make([]string, 0),
        newLinesCh:                   make(chan string),
        newLinesSubscriptionCh:       make(chan *newLinesSubscription),
        deleteNewLinesSubscriptionCh: make(chan *newLinesSubscription),
        stop:                         make(chan bool),
    }

    go w.startControlLoop()
    return w
}

func (w *Watcher) AddLogLine(logLine string) {
    w.newLinesCh <- logLine
}

func (w *Watcher) WaitFor(keyword string, timeout time.Duration) bool {
    newLinesCh := make(chan string)
    foundLineCh := make(chan bool, 2)
    subscription := w.subscribeForNewLines(newLinesCh)
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    go w.searchInNewLines(ctx, subscription, keyword, foundLineCh)
    go w.searchInStoredLines(ctx, keyword, foundLineCh)
    found := w.wait(ctx, foundLineCh)
    cancel()
    w.unsubscribeFromNewLines(subscription)
    return found
}

func (w *Watcher) Stop() error {
    close(w.stop)
    return nil
}

func (w *Watcher) startControlLoop() {
    newLinesSubscriptions := make(map[int64]*newLinesSubscription)
    for {
        select {
        case newLine := <-w.newLinesCh:
            w.storeNewLine(newLine)
            w.broadcastNewLine(newLine, newLinesSubscriptions)
        case subscription := <-w.newLinesSubscriptionCh:
            newLinesSubscriptions[subscription.id] = subscription
        case subscription := <-w.deleteNewLinesSubscriptionCh:
            subscription.active = false
        case <-w.stop:
            return
        }
    }
}

func (w *Watcher) storeNewLine(line string) {
    w.lines = append(w.lines, line)
}

func (w *Watcher) broadcastNewLine(line string, newLinesSubscriptions map[int64]*newLinesSubscription) {
    for _, subscription := range newLinesSubscriptions {
        if !subscription.active {
            close(subscription.newLinesCh)
            delete(newLinesSubscriptions, subscription.id)
            continue
        }
        subscription.newLinesCh <- line
    }
}

func (w *Watcher) subscribeForNewLines(newLinesCh chan string) *newLinesSubscription {
    subscriptionId := w.newLinesSubscriptionIdGenerator.Add(1)
    subscription := &newLinesSubscription{
        id:         subscriptionId,
        newLinesCh: newLinesCh,
        active:     true,
    }
    w.newLinesSubscriptionCh <- subscription
    return subscription
}

func (w *Watcher) unsubscribeFromNewLines(subscription *newLinesSubscription) {
    w.deleteNewLinesSubscriptionCh <- subscription
}

func (w *Watcher) searchInNewLines(ctx context.Context, subs *newLinesSubscription, keyword string, foundLineCh chan bool) {
    for newLine := range subs.newLinesCh {
        if strings.Contains(newLine, keyword) {
            if ctx.Err() != nil {
                // ctx is Done
                continue
            }
            foundLineCh <- true
        }
    }
}

func (w *Watcher) searchInStoredLines(ctx context.Context, keyword string, foundLine chan bool) {
    for _, line := range w.lines {
        if strings.Contains(line, keyword) {
            if ctx.Err() != nil {
                // ctx is Done
                return
            }
            foundLine <- true
            return
        }
    }
}

func (w *Watcher) wait(ctx context.Context, foundLineCh chan bool) bool {
    found := false
    select {
    case <-foundLineCh:
        found = true
    case <-ctx.Done():
    }
    return found
}
