package snapshot_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func TestNewManagerCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "snaps")
	m, err := snapshot.NewManager(dir, 0)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to exist: %v", err)
	}
	_ = m
}

func TestLatestPath(t *testing.T) {
	dir := t.TempDir()
	m, _ := snapshot.NewManager(dir, 0)
	if !strings.HasSuffix(m.LatestPath(), "latest.json") {
		t.Errorf("unexpected latest path: %s", m.LatestPath())
	}
}

func TestArchivePathFormat(t *testing.T) {
	dir := t.TempDir()
	m, _ := snapshot.NewManager(dir, 0)
	p := m.ArchivePath()
	base := filepath.Base(p)
	if !strings.HasPrefix(base, "snapshot_") || !strings.HasSuffix(base, ".json") {
		t.Errorf("unexpected archive path format: %s", base)
	}
}

func TestPruneRemovesOldFiles(t *testing.T) {
	dir := t.TempDir()
	m, _ := snapshot.NewManager(dir, 7)

	// Create an old file
	oldFile := filepath.Join(dir, "snapshot_old.json")
	_ = os.WriteFile(oldFile, []byte(`{}`), 0644)
	oldTime := time.Now().AddDate(0, 0, -10)
	_ = os.Chtimes(oldFile, oldTime, oldTime)

	// Create a recent file
	newFile := filepath.Join(dir, "snapshot_new.json")
	_ = os.WriteFile(newFile, []byte(`{}`), 0644)

	if err := m.Prune(); err != nil {
		t.Fatalf("Prune: %v", err)
	}
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("expected old file to be pruned")
	}
	if _, err := os.Stat(newFile); err != nil {
		t.Error("expected new file to remain")
	}
}

func TestPruneZeroRetainKeepsAll(t *testing.T) {
	dir := t.TempDir()
	m, _ := snapshot.NewManager(dir, 0)

	oldFile := filepath.Join(dir, "snapshot_old.json")
	_ = os.WriteFile(oldFile, []byte(`{}`), 0644)
	oldTime := time.Now().AddDate(0, 0, -100)
	_ = os.Chtimes(oldFile, oldTime, oldTime)

	_ = m.Prune()
	if _, err := os.Stat(oldFile); err != nil {
		t.Error("expected file to remain when retainDays=0")
	}
}
