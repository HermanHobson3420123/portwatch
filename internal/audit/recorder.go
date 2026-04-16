package audit

import (
	"context"

	"portwatch/internal/scanner"
)

// PortEvent carries a port change and its direction.
type PortEvent struct {
	Kind EventKind
	Port scanner.Port
}

// Recorder listens on a channel and writes events to a Log.
type Recorder struct {
	log *Log
}

// NewRecorder wraps a Log for channel-based recording.
func NewRecorder(l *Log) *Recorder {
	return &Recorder{log: l}
}

// Watch consumes events from ch until ctx is cancelled or ch is closed.
func (r *Recorder) Watch(ctx context.Context, ch <-chan PortEvent) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-ch:
			if !ok {
				return
			}
			_ = r.log.Append(ev.Kind, ev.Port, "")
		}
	}
}
