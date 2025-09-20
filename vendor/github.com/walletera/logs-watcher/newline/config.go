package newline

type Config func(*Watcher)

func WithNewLinesChSize(size int) Config {
    return func(w *Watcher) {
        w.newLinesChSize = size
    }
}
