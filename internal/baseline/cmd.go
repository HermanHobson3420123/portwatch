package baseline

import (
	"fmt"
	"io"

	"github.com/user/portwatch/internal/scanner"
)

// CaptureOptions holds parameters for a baseline capture operation.
type CaptureOptions struct {
	Path    string
	Targets []string
	Ports   string
	Out     io.Writer
}

// Capture scans the given targets and saves the result as a new baseline.
func Capture(opts CaptureOptions, sc *scanner.Scanner) error {
	var all []scanner.Port
	for _, target := range opts.Targets {
		ports, err := sc.Scan(target)
		if err != nil {
			return fmt.Errorf("scan %s: %w", target, err)
		}
		all = append(all, ports...)
	}
	if err := Save(opts.Path, all); err != nil {
		return fmt.Errorf("save baseline: %w", err)
	}
	fmt.Fprintf(opts.Out, "baseline captured: %d ports written to %s\n", len(all), opts.Path)
	return nil
}

// Check loads the baseline at path and compares it against current.
// A summary of violations is written to out.
func Check(path string, current []scanner.Port, out io.Writer) (ok bool, err error) {
	b, err := Load(path)
	if err != nil {
		return false, err
	}
	unexpected, missing := Violations(b, current)
	if len(unexpected) == 0 && len(missing) == 0 {
		fmt.Fprintln(out, "baseline check passed: no violations")
		return true, nil
	}
	for _, p := range unexpected {
		fmt.Fprintf(out, "UNEXPECTED port open: %s\n", p.String())
	}
	for _, p := range missing {
		fmt.Fprintf(out, "MISSING expected port: %s\n", p.String())
	}
	return false, nil
}
