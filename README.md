# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes in real time.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git
cd portwatch && go build -o portwatch .
```

---

## Usage

Start monitoring with default settings (checks every 30 seconds):

```bash
portwatch start
```

Specify a custom scan interval and alert on any new or closed ports:

```bash
portwatch start --interval 10s --notify
```

Take a snapshot of the current open ports to use as a baseline:

```bash
portwatch snapshot
```

Run in the foreground with verbose output:

```bash
portwatch start --verbose --foreground
```

### Example Output

```
[INFO]  Baseline captured: 4 ports open
[ALERT] New port detected: 0.0.0.0:8080 (tcp)
[ALERT] Port closed:       127.0.0.1:5432 (tcp)
```

---

## Configuration

`portwatch` can be configured via a YAML file at `~/.portwatch.yaml`:

```yaml
interval: 15s
notify: true
ignore:
  - 127.0.0.1:631
```

---

## License

MIT © 2024 Your Name