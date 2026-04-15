package notify

import (
	"portwatch/internal/scanner"
)

// OpenedEvent constructs an ALERT-level Event for a newly opened port.
func OpenedEvent(p scanner.Port) Event {
	return Event{
		Level:    LevelAlert,
		Message:  "port opened",
		Port:     p.Number,
		Protocol: p.Protocol,
	}
}

// ClosedEvent constructs a WARN-level Event for a port that has closed.
func ClosedEvent(p scanner.Port) Event {
	return Event{
		Level:    LevelWarn,
		Message:  "port closed",
		Port:     p.Number,
		Protocol: p.Protocol,
	}
}

// InfoEvent constructs an INFO-level Event with a custom message.
func InfoEvent(p scanner.Port, msg string) Event {
	return Event{
		Level:    LevelInfo,
		Message:  msg,
		Port:     p.Number,
		Protocol: p.Protocol,
	}
}
