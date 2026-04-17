package tui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func p(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestUpsertAndCount(t *testing.T) {
	tbl := New()
	tbl.Upsert(p("tcp", 80), "open")
	tbl.Upsert(p("tcp", 443), "open")
	if got := tbl.Count(); got != 2 {
		t.Fatalf("expected 2 rows, got %d", got)
	}
}

func TestUpsertOverwrites(t *testing.T) {
	tbl := New()
	tbl.Upsert(p("tcp", 80), "open")
	tbl.Upsert(p("tcp", 80), "closed")
	if got := tbl.Count(); got != 1 {
		t.Fatalf("expected 1 row, got %d", got)
	}
}

func TestRemove(t *testing.T) {
	tbl := New()
	tbl.Upsert(p("tcp", 22), "open")
	tbl.Remove(p("tcp", 22))
	if got := tbl.Count(); got != 0 {
		t.Fatalf("expected 0 rows after remove, got %d", got)
	}
}

func TestRenderContainsHeaders(t *testing.T) {
	tbl := New()
	var buf bytes.Buffer
	tbl.Render(&buf)
	out := buf.String()
	for _, hdr := range []string{"PROTO", "PORT", "STATUS", "SEEN"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestRenderContainsRowData(t *testing.T) {
	tbl := New()
	tbl.Upsert(p("udp", 53), "open")
	var buf bytes.Buffer
	tbl.Render(&buf)
	out := buf.String()
	if !strings.Contains(out, "udp") {
		t.Error("expected proto 'udp' in rendered output")
	}
	if !strings.Contains(out, "53") {
		t.Error("expected port '53' in rendered output")
	}
}
