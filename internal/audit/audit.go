package audit

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	KindOpened EventKind = "opened"
	KindClosed EventKind = "closed"
	KindAlert  EventKind = "alert"
)

// Event is a single audit log entry.
type Event struct {
	Timestamp time.Time    `json:"timestamp"`
	Kind      EventKind    `json:"kind"`
	Port      scanner.Port `json:"port"`
	Message   string       `json:"message,omitempty"`
}

// Log holds an ordered list of audit events.
type Log struct {
	mu     sync.Mutex
	path   string
	events []Event
}

// New opens (or creates) an audit log at path.
func New(path string) (*Log, error) {
	l := &Log{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return l, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &l.events); err != nil {
		return nil, err
	}
	return l, nil
}

// Append adds an event and persists the log.
func (l *Log) Append(kind EventKind, p scanner.Port, msg string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events, Event{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Port:      p,
		Message:   msg,
	})
	return l.flush()
}

// Events returns a copy of all recorded events.
func (l *Log) Events() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

func (l *Log) flush() error {
	data, err := json.MarshalIndent(l.events, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}
