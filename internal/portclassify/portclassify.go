// Package portclassify categorises ports into well-known service tiers.
package portclassify

import "github.com/user/portwatch/internal/scanner"

// Tier represents a broad service category for a port.
type Tier string

const (
	TierSystem    Tier = "system"    // 1–1023
	TierRegistered Tier = "registered" // 1024–49151
	TierDynamic   Tier = "dynamic"   // 49152–65535
	TierUnknown   Tier = "unknown"
)

// Class holds classification details for a single port.
type Class struct {
	Port    scanner.Port
	Tier    Tier
	Service string // human-readable service name, empty if unknown
}

// wellKnown maps port numbers to common service names.
var wellKnown = map[uint16]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Classify returns a Class for the given port.
func Classify(p scanner.Port) Class {
	tier := tierOf(p.Number)
	service := wellKnown[p.Number]
	return Class{Port: p, Tier: tier, Service: service}
}

// ClassifyAll classifies a slice of ports.
func ClassifyAll(ports []scanner.Port) []Class {
	out := make([]Class, len(ports))
	for i, p := range ports {
		out[i] = Classify(p)
	}
	return out
}

// tierOf returns the Tier for a port number.
func tierOf(n uint16) Tier {
	switch {
	case n >= 1 && n <= 1023:
		return TierSystem
	case n >= 1024 && n <= 49151:
		return TierRegistered
	case n >= 49152:
		return TierDynamic
	default:
		return TierUnknown
	}
}
