package portlabel

// Label returns a human-readable service name for well-known ports.
func Label(port uint16, protocol string) string {
	key := portKey(port, protocol)
	if name, ok := wellKnown[key]; ok {
		return name
	}
	if name, ok := wellKnown[portKey(port, "")]; ok {
		return name
	}
	return ""
}

// Annotate returns the label if available, otherwise a fallback string.
func Annotate(port uint16, protocol string) string {
	if l := Label(port, protocol); l != "" {
		return l
	}
	return "unknown"
}

func portKey(port uint16, protocol string) string {
	return protocol + ":" + itoa(port)
}

func itoa(n uint16) string {
	if n == 0 {
		return "0"
	}
	buf := [8]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

var wellKnown = map[string]string{
	":22":    "ssh",
	":25":    "smtp",
	":53":    "dns",
	":80":    "http",
	":110":   "pop3",
	":143":   "imap",
	":443":   "https",
	":3306":  "mysql",
	":5432":  "postgres",
	":6379":  "redis",
	":8080":  "http-alt",
	":8443":  "https-alt",
	":27017": "mongodb",
	tcp(22):   "ssh",
	tcp(80):   "http",
	tcp(443):  "https",
	udp(53):   "dns",
}

func tcp(p uint16) string { return portKey(p, "tcp") }
func udp(p uint16) string { return portKey(p, "udp") }
