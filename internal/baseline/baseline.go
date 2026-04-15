package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// ErrNoBaseline is returned when no baseline file exists.
var ErrNoBaseline = errors.New("baseline: no baseline file found")

// Baseline represents a saved set of expected open ports.
type Baseline struct {
	CreatedAt time.Time      `json:"created_at"`
	Ports     []scanner.Port `json:"ports"`
}

// Save writes the given ports as the current baseline to path.
func Save(path string, ports []scanner.Port) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Ports:     ports,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads the baseline from path.
// Returns ErrNoBaseline if the file does not exist.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoBaseline
		}
		return nil, err
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Violations returns ports that are open but not present in the baseline,
// and ports that are in the baseline but no longer open.
func Violations(base *Baseline, current []scanner.Port) (unexpected, missing []scanner.Port) {
	baseSet := make(map[string]struct{}, len(base.Ports))
	for _, p := range base.Ports {
		baseSet[p.String()] = struct{}{}
	}
	currentSet := make(map[string]struct{}, len(current))
	for _, p := range current {
		currentSet[p.String()] = struct{}{}
	}
	for _, p := range current {
		if _, ok := baseSet[p.String()]; !ok {
			unexpected = append(unexpected, p)
		}
	}
	for _, p := range base.Ports {
		if _, ok := currentSet[p.String()]; !ok {
			missing = append(missing, p)
		}
	}
	return unexpected, missing
}
