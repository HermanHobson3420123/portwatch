package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

func makePorts(specs ...string) []scanner.Port {
	ports := make([]scanner.Port, 0, len(specs))
	for _, s := range specs {
		var p scanner.Port
		if _, err := fmt.Sscanf(s, "%d/%s", &p.Number, &p.Protocol); err == nil {
			ports = append(ports, p)
		}
	}
	return ports
}

func port(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	ports := []scanner.Port{port(80, "tcp"), port(443, "tcp")}
	if err := baseline.Save(path, ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(b.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(b.Ports))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err != baseline.ErrNoBaseline {
		t.Errorf("expected ErrNoBaseline, got %v", err)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	_ = os.WriteFile(path, []byte("not-json{"), 0o644)

	_, err := baseline.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestViolations_Unexpected(t *testing.T) {
	b := &baseline.Baseline{Ports: []scanner.Port{port(80, "tcp")}}
	current := []scanner.Port{port(80, "tcp"), port(9000, "tcp")}

	unexpected, missing := baseline.Violations(b, current)
	if len(unexpected) != 1 || unexpected[0].Number != 9000 {
		t.Errorf("expected unexpected port 9000, got %v", unexpected)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing ports, got %v", missing)
	}
}

func TestViolations_Missing(t *testing.T) {
	b := &baseline.Baseline{Ports: []scanner.Port{port(80, "tcp"), port(443, "tcp")}}
	current := []scanner.Port{port(80, "tcp")}

	unexpected, missing := baseline.Violations(b, current)
	if len(unexpected) != 0 {
		t.Errorf("expected no unexpected ports, got %v", unexpected)
	}
	if len(missing) != 1 || missing[0].Number != 443 {
		t.Errorf("expected missing port 443, got %v", missing)
	}
}

func TestViolations_Clean(t *testing.T) {
	b := &baseline.Baseline{Ports: []scanner.Port{port(80, "tcp")}}
	current := []scanner.Port{port(80, "tcp")}

	unexpected, missing := baseline.Violations(b, current)
	if len(unexpected) != 0 || len(missing) != 0 {
		t.Errorf("expected no violations, got unexpected=%v missing=%v", unexpected, missing)
	}
}
