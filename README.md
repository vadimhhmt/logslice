# logslice

A CLI tool to filter and slice structured log files by time range and field patterns.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build -o logslice .
```

## Usage

```bash
# Filter logs between two timestamps
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" app.log

# Filter by field pattern
logslice --field "level=error" --field "service=api" app.log

# Combine time range and field filters
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --field "level=error" app.log

# Read from stdin
cat app.log | logslice --from "2024-01-15T08:00:00Z" --field "status=500"
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start of time range (RFC3339) |
| `--to` | End of time range (RFC3339) |
| `--field` | Field pattern to match (key=value), repeatable |
| `--format` | Log format: `json`, `logfmt` (default: `json`) |
| `--output` | Output file path (default: stdout) |

## Supported Formats

- **JSON** — newline-delimited JSON logs
- **logfmt** — key=value structured logs

## License

MIT © [yourusername](https://github.com/yourusername)