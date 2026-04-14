package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Level controls the verbosity of alert output.
type Level int

const (
	LevelInfo  Level = iota
	LevelWarn
)

// Alerter writes human-readable change notifications to a writer.
type Alerter struct {
	out   io.Writer
	level Level
}

// New returns an Alerter that writes to out.
func New(out io.Writer, level Level) *Alerter {
	if out == nil {
		out = os.Stdout
	}
	return &Alerter{out: out, level: level}
}

// Notify formats and writes a single Change event.
func (a *Alerter) Notify(c monitor.Change) {
	timestamp := time.Now().Format(time.RFC3339)
	var action, symbol string
	if c.Opened {
		action = "OPENED"
		symbol = "+"
	} else {
		action = "CLOSED"
		symbol = "-"
	}
	fmt.Fprintf(a.out, "[%s] [%s] port %s %s\n", timestamp, symbol, c.Port.Address, action)
}

// Watch consumes changes from the channel until it is closed.
func (a *Alerter) Watch(changes <-chan monitor.Change) {
	for c := range changes {
		a.Notify(c)
	}
}
