<p align="center">
  <img src="assets/logo.png" alt="GoScan Logo" width="200" style="border-radius: 50%;" />
</p>

# GoScan · 📡🔍

[![Go](https://img.shields.io/badge/Go-1.24.2-00ADD8?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-v0.9.0-yellowgreen)]()

A fast and minimal TCP port scanner built in Go. Supports CIDR/IP ranges, port ranges, concurrency, banner grabbing, and structured output.

---

## ✨ Features

- Parse IPs (CIDR, ranges, comma-separated lists).
- Parse ports (single, range, list).
- Concurrent TCP scanning using goroutines.
- Customizable timeout and worker count.
- Banner grabbing support (basic).
- Output in plaintext or JSON.
- Easy to extend and clean project structure.

---

## 📦 Installation

### 🔧 Option 1: Build from source

```bash
git clone https://github.com/0xM3H51N/GoScan.git
cd GoScan
go build -o goscan
```

### 📥 Option 2: Install via Go

```bash
go install github.com/0xM3H51N/GoScan@latest
```

Make sure `GOBIN` is in your `$PATH`.

---

## 🚀 Usage

```bash
goscan --ip 192.168.1.1,192.168.1.10-20 --port 80,443,8000-8080 --workers 50 --timeout 3 --output result.json --format json
```

### Flags

| Flag         | Description                                                         | Default |
|--------------|---------------------------------------------------------------------|---------|
| `-i, --ip`   | Target IP(s). Supports single, list, CIDR, or range                 |         |
| `-f, --file` | Load IPs from file (one per line)                                   |         |
| `-p, --port` | Target ports. Accepts single, list, or range                        |         |
| `-w`         | Number of concurrent workers                                        | 5       |
| `--timeout`  | Connection timeout in seconds                                       | 5       |
| `-x`         | Output format: `PLAINTEXT` or `JSON`                                |         |
| `-o`         | Write results to file instead of stdout                             |         |
| `--version`  | Show version and exit                                               |         |

---

## 🧪 Output Examples

### Plaintext
```                           
=================
10.0.2.15
=================
port	status	version
8000	close	
22	open	SSH-2.0-OpenSSH_9.9p1 Ubuntu-3ubuntu3.1

=================
192.168.1.113
=================
port	status	version
8000	open	
22	close	
```

### JSON
```json
{
  "10.0.2.15": [
    {
      "port": 8000,
      "status": "close",
      "banner": ""
    },
    {
      "port": 22,
      "status": "open",
      "banner": "SSH-2.0-OpenSSH_9.9p1 Ubuntu-3ubuntu3.1\r\n"
    }
  ],
  "192.168.1.113": [
    {
      "port": 8000,
      "status": "open",
      "banner": ""
    },
    {
      "port": 22,
      "status": "close",
      "banner": ""
    }
  ]
}
```

---

## 📁 Project Structure

<pre>
GoScan/
├── assets/
├── cmd/         # Cobra CLI and command entry
├── internal/    # Helpers (IP/port parsing, banner grabbing)
├── core/        # Data models and types
├── main.go      # App entry point
├── go.mod
└── README.md
</pre>

---

## ✅ Completed Features

- IP & port parsing (CIDR, ranges, etc.)
- Concurrent port scanning
- Basic banner grabbing
- JSON/plaintext output
- Cobra CLI integration

---

## 🧪 Coming Soon (maybe)

- Advanced banner grabbing (protocol detection)
- Integration/unit tests
---

## 🔖 Versioning

This project uses [Semantic Versioning](https://semver.org/).  
You are currently viewing **v0.9.0**, an early development preview.

---
## 🛠️ Craftsmanship

GoScan was **entirely built from scratch** without AI code generation tools.  
Every line of code, design decision, and feature was written manually to maximize learning, reinforce Go fundamentals, and demonstrate full ownership of the development process.

---

## 🧑‍💻 Author

**[@0xM3H51N](https://github.com/0xM3H51N)**
**[email](m3h51n@protonmail.com)**

---

Licensed under the MIT License.