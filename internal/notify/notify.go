package notify

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event holds the data for a single notification event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      uint16
	Protocol  string
}

const defaultTemplate = "[{{.Timestamp.Format \"2006-01-02 15:04:05\"}}] {{.Level}} port={{.Port}}/{{.Protocol}} {{.Message}}\n"

// Notifier writes formatted notification events to a writer.
type Notifier struct {
	out  io.Writer
	tmpl *template.Template
}

// New creates a Notifier writing to out using the default template.
// If out is nil, os.Stderr is used.
func New(out io.Writer) (*Notifier, error) {
	if out == nil {
		out = os.Stderr
	}
	tmpl, err := template.New("event").Parse(defaultTemplate)
	if err != nil {
		return nil, fmt.Errorf("notify: parse template: %w", err)
	}
	return &Notifier{out: out, tmpl: tmpl}, nil
}

// NewWithTemplate creates a Notifier using a custom Go text/template string.
func NewWithTemplate(out io.Writer, tmplStr string) (*Notifier, error) {
	if out == nil {
		out = os.Stderr
	}
	tmpl, err := template.New("event").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("notify: parse template: %w", err)
	}
	return &Notifier{out: out, tmpl: tmpl}, nil
}

// Send formats and writes a single Event.
func (n *Notifier) Send(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	if err := n.tmpl.Execute(n.out, e); err != nil {
		return fmt.Errorf("notify: render event: %w", err)
	}
	return nil
}

// Watch consumes events from ch and calls Send for each one.
// It blocks until ch is closed.
func (n *Notifier) Watch(ch <-chan Event) {
	for e := range ch {
		_ = n.Send(e)
	}
}
