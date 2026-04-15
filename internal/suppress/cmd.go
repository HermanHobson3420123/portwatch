// Package suppress — cmd.go provides CLI-facing helpers for managing
// suppression windows (add, list, clear).
package suppress

import (
	"fmt"
	"io"
	"time"
)

// AddWindow adds a suppression window for port/protocol lasting duration d
// and persists the updated list to path.
func AddWindow(path string, port uint16, protocol string, d time.Duration, w io.Writer) error {
	l := New()
	if err := l.LoadFromFile(path); err != nil {
		return fmt.Errorf("load suppressions: %w", err)
	}
	until := time.Now().Add(d)
	l.Add(port, protocol, until)
	if err := l.SaveToFile(path); err != nil {
		return fmt.Errorf("save suppressions: %w", err)
	}
	fmt.Fprintf(w, "suppressed %d/%s until %s\n", port, protocol, until.Format(time.RFC3339))
	return nil
}

// ListWindows prints all active suppression windows to w.
func ListWindows(path string, w io.Writer) error {
	l := New()
	if err := l.LoadFromFile(path); err != nil {
		return fmt.Errorf("load suppressions: %w", err)
	}
	active := l.Active()
	if len(active) == 0 {
		fmt.Fprintln(w, "no active suppression windows")
		return nil
	}
	for _, e := range active {
		fmt.Fprintf(w, "%d/%s\tuntil %s\n", e.Port, e.Protocol, e.Until.Format(time.RFC3339))
	}
	return nil
}

// ClearExpired removes expired entries from the persisted file.
func ClearExpired(path string, w io.Writer) error {
	l := New()
	if err := l.LoadFromFile(path); err != nil {
		return fmt.Errorf("load suppressions: %w", err)
	}
	l.Purge()
	if err := l.SaveToFile(path); err != nil {
		return fmt.Errorf("save suppressions: %w", err)
	}
	fmt.Fprintln(w, "expired suppression windows cleared")
	return nil
}
