package tui

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Row represents a single port row in the TUI table.
type Row struct {
	Port     scanner.Port
	Status   string
	SeenAt   time.Time
}

// Table holds the current state of monitored ports for display.
type Table struct {
	mu   sync.RWMutex
	rows map[string]Row
}

// New returns an initialised Table.
func New() *Table {
	return &Table{rows: make(map[string]Row)}
}

func rowKey(p scanner.Port) string {
	return fmt.Sprintf("%s/%d", p.Proto, p.Number)
}

// Upsert adds or updates a port row.
func (t *Table) Upsert(p scanner.Port, status string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.rows[rowKey(p)] = Row{Port: p, Status: status, SeenAt: time.Now()}
}

// Remove deletes a port row.
func (t *Table) Remove(p scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.rows, rowKey(p))
}

// Render writes a formatted table to w.
func (t *Table) Render(w io.Writer) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	fmt.Fprintln(w, strings.Repeat("-", 40))
	fmt.Fprintf(w, "%-8s %-6s %-10s %s\n", "PROTO", "PORT", "STATUS", "SEEN")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, r := range t.rows {
		fmt.Fprintf(w, "%-8s %-6d %-10s %s\n",
			r.Port.Proto,
			r.Port.Number,
			r.Status,
			r.SeenAt.Format("15:04:05"),
		)
	}
	fmt.Fprintln(w, strings.Repeat("-", 40))
}

// Count returns the number of rows currently tracked.
func (t *Table) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.rows)
}
