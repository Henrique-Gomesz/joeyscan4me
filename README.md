# JoeyScan4Me

A easy-to-use recon tool kit for subdomain enumeration, HTTP probing, web crawling, and screenshot capturing.
<img width="692" height="520" alt="image" src="https://github.com/user-attachments/assets/f68e172e-d69b-4d8d-b0ef-2c4304d13eea" />

## Features

- **Subdomain Enumeration**: Uses [Subfinder](https://github.com/projectdiscovery/subfinder) to discover subdomains
- **HTTP Probing**: Uses [HTTPX](https://github.com/projectdiscovery/httpx) to identify live web services
- **Web Crawling**: Uses [Katana](https://github.com/projectdiscovery/katana) to crawl discovered websites
- **Screenshot Capture and Dashboard**: Uses [Gowitness](https://github.com/sensepost/gowitness) to capture screenshots with database storage and a web dashboard for easy viewing

## Installation
You can install JoeyScan4Me using the following methods:
Require Go 1.21 or higher.

## From source
```bash
go install github.com/Henrique-Gomesz/JoeyScan4Me/cmd/JoeyScan4Me@latest
```
## Manual build
```bash
git clone https://github.com/Henrique-Gomesz/JoeyScan4Me.git
cd JoeyScan4Me
go build -o joeyscan4me cmd/JoeyScan4Me/main.go
```

## Usage
```bash
$ joeyscan4me -h

JoeyScan4Me - Simple and helpful recon toolkit

    |\__/,|   ('\
  _.|o o  |_   ) )
-(((---(((--------
by: Henrique-Gomesz


Usage:
  joeyscan4me [flags]

Flags:
  -d string                domain to scan (e.g. example.com)
  -w string                working directory for output files, defaults to current directory (default "./")
  -server                  start gowitness server at the end of scan to view screenshots
  -profile string          scan profile: balanced, stealth, aggressive, bugbounty (default "balanced")
  -summary-json            write scan summary as JSON file (default true)
  -summary-file string     scan summary JSON output filename (default "scan_summary.json")

  # Subfinder tuning
  -subfinder-threads int   number of subfinder threads (default 10)
  -subfinder-timeout int   subfinder timeout in seconds (default 30)
  -subfinder-max-time int  subfinder max enumeration time in minutes (default 10)

  # HTTPX tuning
  -httpx-threads int       number of httpx threads (default 50)
  -httpx-rate-limit int    httpx requests per second (default 150)
  -httpx-timeout int       httpx timeout in seconds (default 10)
  -httpx-ports string      custom ports list (nmap-style), e.g. "80,443,8080" or "http:80,https:8443"

  # Katana tuning
  -katana-depth int        maximum katana crawl depth (default 3)
  -katana-timeout int      katana timeout in seconds (default 10)
  -katana-concurrency int  katana concurrent crawling goroutines (default 100)
  -katana-parallelism int  katana URL processing goroutines (default 100)
  -katana-rate-limit int   katana requests per second (default 150)
```

## Example
Running a scan on example.com and starting the gowitness server at the end:
```bash
joeyscan4me -d example.com -w /path/to/output -server
```

Aggressive recon profile with custom HTTP ports:
```bash
joeyscan4me \
  -d example.com \
  -w /path/to/output \
  -httpx-threads 500 \
  -httpx-rate-limit 500 \
  -httpx-ports "66,80,81,82,83,84,85,86,87,88,89,90,443,444,445,3000,3001,3002,8080,8081,8082,8443,8888,9000,9443,10000,20000,30000,50000,65535" \
  -katana-depth 4 \
  -katana-concurrency 150 \
  -katana-parallelism 150
```

Bug bounty profile (pre-tuned with broad web ports and high throughput):
```bash
joeyscan4me -d example.com -profile bugbounty
```

Bug bounty profile with custom override (your own target pace):
```bash
joeyscan4me \
  -d example.com \
  -profile bugbounty \
  -httpx-rate-limit 300 \
  -katana-rate-limit 180
```

Stealth profile (lower noise):
```bash
joeyscan4me -d example.com -profile stealth
```

Disable JSON summary file:
```bash
joeyscan4me -d example.com -summary-json=false
```

## Output Files
The output files will be stored in the specified working directory (or current directory by default) with the following structure:
```
/output
├── example.com/
│   ├── subdomains.txt              # List of discovered subdomains
│   ├── up_subdomains.txt           # List of live HTTP services
│   ├── up_subdomains_with_tech.txt # Live services with technology detection
│   ├── crawling_results.txt        # Web crawling results
│   ├── scan_summary.json           # Consolidated scan summary and counts
│   └── screenshots/
│       └── gowitness.sqlite3       # Screenshot database
│       └── screenshots...       # Screenshot images
```
