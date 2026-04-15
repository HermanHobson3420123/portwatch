package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format controls how reports are rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Event represents a port state change event for reporting.
type Event struct {
	Timestamp time.Time    `json:"timestamp"`
	Kind      string       `json:"kind"` // "opened" or "closed"
	Port      scanner.Port `json:"port"`
}

// Reporter writes port change events to an output stream.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter that writes to out using the given format.
func New(out io.Writer, format Format) *Reporter {
	return &Reporter{out: out, format: format}
}

// ReportOpened emits an event for a newly opened port.
func (r *Reporter) ReportOpened(p scanner.Port) error {
	return r.emit(Event{
		Timestamp: time.Now().UTC(),
		Kind:      "opened",
		Port:      p,
	})
}

// ReportClosed emits an event for a recently closed port.
func (r *Reporter) ReportClosed(p scanner.Port) error {
	return r.emit(Event{
		Timestamp: time.Now().UTC(),
		Kind:      "closed",
		Port:      p,
	})
}

func (r *Reporter) emit(e Event) error {
	switch r.format {
	case FormatJSON:
		data, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("reporter: marshal: %w", err)
		}
		_, err = fmt.Fprintf(r.out, "%s\n", data)
		return err
	default:
		symbol := "+"
		if e.Kind == "closed" {
			symbol = "-"
		}
		_, err := fmt.Fprintf(r.out, "[%s] %s [%s] %s\n",
			e.Timestamp.Format(time.RFC3339),
			symbol,
			e.Kind,
			e.Port.String(),
		)
		return err
	}
}
