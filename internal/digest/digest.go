// Package digest produces a stable fingerprint for a port scan result set,
// allowing quick equality checks without deep comparison.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"portwatch/internal/scanner"
)

// Sum returns a deterministic SHA-256 hex string for the given port list.
// Ports are sorted by protocol+number before hashing so that order does not
// affect the result.
func Sum(ports []scanner.Port) string {
	keys := make([]string, 0, len(ports))
	for _, p := range ports {
		keys = append(keys, portKey(p))
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		_, _ = fmt.Fprintln(h, k)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Equal returns true when two port lists produce the same digest.
func Equal(a, b []scanner.Port) bool {
	return Sum(a) == Sum(b)
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Number)
}
