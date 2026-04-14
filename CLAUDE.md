# CLAUDE.md — JoeyScan4Me

> This file provides context for AI coding assistants working on this codebase.

## Project Overview

JoeyScan4Me is a CLI recon toolkit written in Go that orchestrates a multi-stage pipeline of best-in-class open-source security tools:

1. **Subfinder** — subdomain enumeration
2. **HTTPX** — HTTP probing to identify live web services + technology detection
3. **Katana** — web crawling (runs concurrently with Gowitness)
4. **Gowitness** — screenshot capture + SQLite database + optional web dashboard (runs concurrently with Katana)

All stages are wired together in `pkg/runner/runner.go`. Katana and Gowitness run in parallel goroutines after Subfinder and HTTPX finish in sequence.

## Module & Version

```
module github.com/henrique-gomesz/joeyscan4me
go 1.25.1
current version: 1.1.2  (set in cmd/joeyscan4me/main.go as var Version string)
```

## Directory Structure

```
joeyscan4me/
├── cmd/
│   └── joeyscan4me/
│       └── main.go          # Entry point: prints banner, parses opts, calls StartScan
├── pkg/
│   ├── logging/
│   │   └── logger.go        # LogInfo (cyan), LogSuccess (green), LogError (red), PrintBanner
│   └── runner/
│       ├── runner.go        # StartScan — pipeline orchestration
│       ├── options.go       # CLI flags, profile presets, domain/tunable validation
│       ├── subfinder.go     # Subfinder integration
│       ├── httpx.go         # HTTPX integration (probing + tech detection)
│       ├── katana.go        # Katana integration (web crawling)
│       ├── gowitness.go     # Gowitness integration (screenshots + server)
│       ├── output.go        # File I/O helpers and output path constants
│       └── summary.go       # JSON scan summary writer
├── go.mod
├── go.sum
└── README.md
```

## Key Design Patterns

### Adding a New Tool Stage

Every tool stage follows the same pattern — see `subfinder.go` as the reference:

1. Create `pkg/runner/<toolname>.go`
2. Implement `Run<Toolname>(opt *Options) error`
3. Add resume guard at the top:
   ```go
   if opt.Resume && FileNonEmpty(filePath) {
       logging.LogInfo("Skipping <toolname> — output already exists: " + filePath)
       return nil
   }
   ```
4. Register the output filename constant in `output.go`:
   ```go
   var <Toolname>OutputFile = "output_name.txt"
   ```
5. Wire it into `runner.go` (sequential → add after HTTPX; concurrent → add a goroutine like Katana/Gowitness)
6. If it produces countable output, add a count field to `summaryCounts` in `summary.go`

### Adding a New CLI Flag

All flags live in `pkg/runner/options.go`:
1. Add a field to the `Options` struct
2. Register it in `ParseOptions()` using the `goflags` flagset
3. If it's a per-tool tunable, add it to `profilePreset` and `validateTunables()`
4. Update `applyProfile()` to apply the preset value when the flag wasn't explicitly set
5. Add the new flag to all four profile presets in `getProfilePreset()`

### Adding a New Profile Preset

Add a new `case` in `getProfilePreset()` in `options.go`. The four existing profiles are: `balanced`, `stealth`, `aggressive`, `bugbounty`.

### Resume Support

Resume is checked per-stage using `FileNonEmpty(filePath)`. Each stage must check `opt.Resume` itself. The convention is to skip when the output file already exists and is non-empty. To force a stage to re-run, the user deletes its output file.

## Output Files

All output is generated under `<workdir>/output/<domain>/`:

| File | Constant | Produced by |
|---|---|---|
| `subdomains.txt` | `SubfinderOutputFile` | Subfinder |
| `up_subdomains.txt` | `HttpxOutputFile` | HTTPX |
| `up_subdomains_with_tech.txt` | `HttpxTechOutputFile` | HTTPX |
| `crawling_results.txt` | `KatanaOutputFile` | Katana |
| `scan_summary.json` | `opt.SummaryFile` | summary.go |
| `screenshots/gowitness.sqlite3` | — | Gowitness |
| `screenshots/*.png` | — | Gowitness |

Helper: `GetOutputFilePath(workdir, domain)` returns the base output dir.

## Logging

Use `pkg/logging` for all terminal output. Never use `fmt.Println` directly in runner code.

```go
logging.LogInfo("...")    // [i] cyan — progress/status
logging.LogSuccess("...") // [✓] green — stage completed
logging.LogError("...", err) // [x] red — errors (err may be nil)
logging.LogText("...")    // white — plain text
```

## Core Dependencies

| Package | Role |
|---|---|
| `github.com/projectdiscovery/subfinder/v2/pkg/runner` | Subdomain enumeration |
| `github.com/projectdiscovery/httpx` | HTTP probing + tech detection |
| `github.com/projectdiscovery/katana` | Web crawling |
| `github.com/sensepost/gowitness` | Screenshots + web dashboard |
| `github.com/projectdiscovery/goflags` | CLI flag parsing |
| `github.com/fatih/color` | Colored terminal output |

## Common Commands

```bash
# Build
go build -o joeyscan4me cmd/joeyscan4me/main.go

# Install from source
go install github.com/henrique-gomesz/joeyscan4me/cmd/joeyscan4me@main

# Run (basic)
go run cmd/joeyscan4me/main.go -d example.com

# Run with profile
go run cmd/joeyscan4me/main.go -d example.com -profile bugbounty

# Run tests
go test ./...

# Tidy dependencies
go mod tidy
```

## Conventions

- **Error wrapping**: Always wrap errors with `fmt.Errorf("context: %w", err)` before returning
- **Stage files**: One file per tool (`subfinder.go`, `httpx.go`, etc.) — do not mix tool logic
- **No global state**: Everything is passed through `*Options`
- **Exit on bad input**: Call `os.Exit(1)` only in `options.go` (validation), never in runner stages
- **Concurrent stages**: Only stages that don't depend on each other's output run in parallel (currently Katana + Gowitness)
- **Version bump**: Update `var Version string` in `cmd/joeyscan4me/main.go` for each release
