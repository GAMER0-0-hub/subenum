# SUBENUM v3

<p align="center">
  <img src="assets/demo.png" alt="Subenum Demo" width="800">
</p>

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![Recon](https://img.shields.io/badge/Recon-Subdomain%20Enumeration-blue)
![Status](https://img.shields.io/badge/Status-Stable-green)

SUBENUM is a professional subdomain enumeration framework written in Go, designed for bug bounty hunters, penetration testers, and security researchers.

It focuses on:
- Accuracy
- Clean output
- Strong UX
- Full deduplication
- Real-world recon workflows

---

## Features

- Supports:
  - example.com
  - *.example.com
  - *.example.*
- Multiple targets via file input
- 100% deduplication (tools + manual input)
- Integrated tools:
  - subfinder
  - assetfinder
  - findomain
  - chaos (API supported)
  - crt.sh
  - httpx
- Interactive UX with spinner & colored output
- Chaos API key auto-save (~/.config/subenum/chaos.key)
- Live host detection using httpx
- Automatic filtering for HTTP 200 responses
- Clean output structure per target
- Cross-platform (Linux / macOS)

---

## Workflow

1. Load targets (single domain or file)
2. Run enumeration tools concurrently
3. Merge and deduplicate all results
4. Optional manual subdomain input
5. Check live hosts using httpx
6. Filter HTTP 200 responses
7. Save clean recon results

---

## Installation (Go – Recommended)

### Requirements

- Go 1.20+
- The following tools installed and available in PATH:
  - subfinder
  - assetfinder
  - findomain
  - chaos (optional – requires API key)
  - httpx

Missing tools are skipped automatically.

### Install

go install github.com/GAMER0_0/subenum@latest

Ensure Go binaries are in your PATH:

export PATH=$PATH:$(go env GOPATH)/bin

---

## Usage

subenum [options]

Target Options:
- -d string    Target domain (example.com | *.example.com | *.example.*)
- -l file      File with list of domains

Output Options:
- -o dir       Output directory (default: subdomain_enu)

General:
- -h           Show help message

---

## Examples

Single domain:
subenum -d example.com

Wildcard scope:
subenum -d "*.example.*"

Scope file with custom output:
subenum -l scope.txt -o /path/recon

---

## Output Structure

subdomain_enu/
└── example.com/
    ├── subdomains.txt
    ├── all_subdomains.txt   (only if manual input added)
    ├── httpx.txt
    └── httpx_200.txt

---

## Disclaimer

This tool is intended only for assets you own or have explicit permission to test.
The author is not responsible for misuse.

---

## Author

GAMER_0_0  
Crafted with precision for professional recon workflows.
