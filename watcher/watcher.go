package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type EventType int

const (
	Created EventType = iota
	Modified
	Deleted
)

type Event struct {
	Type EventType
	Path string
}

type Watcher struct {
	root     string
	interval time.Duration

	stop     chan struct{}
	stopOnce sync.Once
}

func New(root string, interval time.Duration) *Watcher {
	return &Watcher{
		root:     root,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

func (w *Watcher) Watch(events chan<- Event) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	previous := w.scan()

	for {
		select {
		case <-ticker.C:
			current := w.scan()
			w.compare(previous, current, events)
			previous = current

		case <-w.stop:
			return
		}
	}
}

func (w *Watcher) Stop() {
	w.stopOnce.Do(func() {
		close(w.stop)
	})
}

func (w *Watcher) scan() map[string]time.Time {
	files := make(map[string]time.Time)

	filepath.Walk(w.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		files[path] = info.ModTime()

		return nil
	})

	return files
}

func (w *Watcher) compare(
	previous map[string]time.Time,
	current map[string]time.Time,
	events chan<- Event,
) {
	for path, modTime := range current {
		oldModTime, exists := previous[path]

		if !exists {
			events <- Event{
				Type: Created,
				Path: path,
			}

			continue
		}

		if !oldModTime.Equal(modTime) {
			events <- Event{
				Type: Modified,
				Path: path,
			}
		}
	}

	for path := range previous {
		if _, exists := current[path]; !exists {
			events <- Event{
				Type: Deleted,
				Path: path,
			}
		}
	}
}
