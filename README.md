# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

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

Start the daemon with default settings:

```bash
portwatch start
```

Specify a custom polling interval and alert on any new or closed ports:

```bash
portwatch start --interval 30s --notify
```

Take a snapshot of the current port state to use as a baseline:

```bash
portwatch snapshot --output baseline.json
```

Watch against an existing baseline:

```bash
portwatch start --baseline baseline.json
```

### Example Output

```
[INFO]  Watching ports... (interval: 30s)
[ALERT] New port detected: 0.0.0.0:8080 (tcp)
[ALERT] Port closed:       127.0.0.1:5432 (tcp)
```

---

## Configuration

`portwatch` can be configured via a `portwatch.yaml` file in the working directory or via CLI flags. Run `portwatch --help` for a full list of options.

---

## License

MIT © 2024 yourusername