# Go Port Scanner

A simple and fast TCP port scanner written in Go that scans open ports on a given host using concurrent connections.

---

## Features

- Scans a specified range of TCP ports on a host (e.g., `localhost`, IP address, or domain)
- Detects open ports (ports that accept TCP connections)
- Runs scans in parallel using goroutines for increased speed
- Default connection timeout set to 1 second

---

## Requirements

- Go 1.20 or newer

---

## Usage

Clone the repository:

```bash
git clone https://github.com/ManceriuMadalin/go-port-scanner.git
cd go-port-scanner
go run main.go <host> <startPort> <endPort>
```

Example:
go run main.go localhost 20 1024
This will scan ports from 20 to 1024 on the local machine.

## Build

You can build an executable binary with:

```bash
go build -o portscanner
```

Then run it like this:

``` bash
./portscanner localhost 20 1024
```

(On Windows: portscanner.exe)


## Author

- Created with ❤️ by Manceriu Mădălin
