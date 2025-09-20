package newline

import (
    "context"
    "strings"
    "sync/atomic"
    "time"
)

type Watcher struct {
    lines                           []string
    newLinesCh                      chan string
    newLinesChSize                  int
    newLinesSubscriptionIdGenerator atomic.Int64
    newLinesSubscriptionCh          chan *newLinesSubscription
    deleteNewLinesSubscriptionCh    chan *newLinesSubscription
    stop                            chan bool
}

type newLinesSubscription struct {
    id         int64
    newLinesCh chan string
}

func NewWatcher(configs ...Config) *Watcher {
    w := &Watcher{
        lines:                        make([]string, 0),
        newLinesCh:                   make(chan string),
        newLinesChSize:               100,
        newLinesSubscriptionCh:       make(chan *newLinesSubscription),
        deleteNewLinesSubscriptionCh: make(chan *newLinesSubscription),
        stop:                         make(chan bool),
    }

    for _, config := range configs {
        config(w)
    }

    go w.startControlLoop()
    return w
}

func (w *Watcher) AddLogLine(logLine string) {
    w.newLinesCh <- logLine
}

func (w *Watcher) WaitForNTimes(keyword string, timeout time.Duration, n int) bool {
    newLinesCh := make(chan string, w.newLinesChSize)
    subscription := w.subscribeForNewLines(newLinesCh)
    foundLineChSize := 1
    if n > 1 {
        foundLineChSize = n
    }
    foundLineCh := make(chan bool, foundLineChSize)
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    go w.searchInNewLines(ctx, subscription, keyword, foundLineCh)
    go w.searchInStoredLines(ctx, keyword, foundLineCh)
    found := w.wait(ctx, foundLineCh, n)
    w.unsubscribeFromNewLines(subscription)
    cancel()
    return found
}

func (w *Watcher) WaitFor(keyword string, timeout time.Duration) bool {
    return w.WaitForNTimes(keyword, timeout, 1)
}

func (w *Watcher) Stop() error {
    close(w.stop)
    return nil
}

func (w *Watcher) startControlLoop() {
    newLinesSubscriptions := make(map[int64]*newLinesSubscription)
    for {
        select {
        case newLine, ok := <-w.newLinesCh:
            if !ok {
                return
            }
            w.storeNewLine(newLine)
            w.broadcastNewLine(newLine, newLinesSubscriptions)
        case subscription := <-w.newLinesSubscriptionCh:
            newLinesSubscriptions[subscription.id] = subscription
        case subscription := <-w.deleteNewLinesSubscriptionCh:
            delete(newLinesSubscriptions, subscription.id)
            close(subscription.newLinesCh)
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
        select {
        case subscription.newLinesCh <- line:
        default:
        }

    }
}

func (w *Watcher) subscribeForNewLines(newLinesCh chan string) *newLinesSubscription {
    subscriptionId := w.newLinesSubscriptionIdGenerator.Add(1)
    subscription := &newLinesSubscription{
        id:         subscriptionId,
        newLinesCh: newLinesCh,
    }
    w.newLinesSubscriptionCh <- subscription
    return subscription
}

func (w *Watcher) unsubscribeFromNewLines(subscription *newLinesSubscription) {
    w.deleteNewLinesSubscriptionCh <- subscription
}

func (w *Watcher) searchInNewLines(ctx context.Context, subs *newLinesSubscription, keyword string, foundLineCh chan bool) {
    for {
        select {
        case newLine := <-subs.newLinesCh:
            if strings.Contains(newLine, keyword) {
                if ctx.Err() != nil {
                    // ctx is Done
                    continue
                }
                select {
                case foundLineCh <- true:
                default:
                }
            }
        case <-ctx.Done():
            return
        }
    }
}

func (w *Watcher) searchInStoredLines(ctx context.Context, keyword string, foundLineCh chan bool) {
    for _, line := range w.lines {
        if strings.Contains(line, keyword) {
            if ctx.Err() != nil {
                // ctx is Done
                return
            }
            select {
            case foundLineCh <- true:
            default:
            }
        }
    }
}

func (w *Watcher) wait(ctx context.Context, foundLineCh chan bool, n int) bool {
    for i := 0; i < n; {
        select {
        case <-foundLineCh:
            i++
        case <-ctx.Done():
            return false
        }
    }
    return true
}
