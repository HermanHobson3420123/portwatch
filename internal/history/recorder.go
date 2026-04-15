package history

import (
	"log"
	"time"

	"portwatch/internal/scanner"
)

// ChangeEvent carries a port and whether it was opened or closed.
type ChangeEvent struct {
	Type EventType
	Port scanner.Port
}

// Recorder listens on a channel of ChangeEvents and persists them to a history file.
type Recorder struct {
	path    string
	maxAge  time.Duration
	events  <-chan ChangeEvent
	done    chan struct{}
}

// NewRecorder creates a Recorder that writes events to path, pruning entries
// older than maxAge on every write.
func NewRecorder(path string, maxAge time.Duration, events <-chan ChangeEvent) *Recorder {
	return &Recorder{
		path:   path,
		maxAge: maxAge,
		events: events,
		done:   make(chan struct{}),
	}
}

// Run starts consuming events until the channel is closed or Stop is called.
func (r *Recorder) Run() {
	for {
		select {
		case e, ok := <-r.events:
			if !ok {
				return
			}
			r.record(e)
		case <-r.done:
			return
		}
	}
}

// Stop signals the recorder to stop processing events.
func (r *Recorder) Stop() {
	close(r.done)
}

func (r *Recorder) record(e ChangeEvent) {
	h, err := Load(r.path)
	if err != nil {
		log.Printf("history recorder: load: %v", err)
		h = &History{}
	}
	h.Append(e.Type, e.Port)
	if r.maxAge > 0 {
		h = Prune(h, r.maxAge)
	}
	if err := Save(r.path, h); err != nil {
		log.Printf("history recorder: save: %v", err)
	}
}
