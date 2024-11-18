# DNSLY

![Go](https://img.shields.io/badge/Go-1.18%2B-blue) ![License](https://img.shields.io/badge/License-MIT-blue) ![Release](https://img.shields.io/github/v/release/vehbiu/dnsly)  
A high-performance DNS blocker written in Go that blocks DNS requests to selected websites. Perfect for district-wide or network-wide filtering.

## Features

- Blocks DNS requests to specified domains (e.g., `example.com`, `blocked.com`, etc.)
- Supports DNS over HTTPS (DoH) for DNS resolution
- Logging of blocked requests with source IPs
- Customizable configuration with a simple JSON file

## Installation

### Prerequisites

- Go 1.18+
- DNS server running on port 53

### Clone the repository

```bash
git clone https://github.com/vehbiu/dnsly.git
cd dnsly
```

### Build the application

```bash
go build -o dnsly main.go
```

## Configuration

Edit the `config.json` file to customize blocked domains and DNS resolver endpoint. **CHANGE THE PASSWORD** in the configuration file for security.

```json
{
  "blocked_domains": ["example.com", "blocked.com", "adult-website.com"],
  "search_dns": "1.1.1.1",
  "http_password": "Change-This-Password"
}
```

### Configuration Fields

- `blocked_domains`: List of domains to block.
- `search_dns`: The DNS server used for resolving non-blocked domains.
- `http_password`: Password to see logs.

## Usage

### Start the DNS Server

```bash
./dnsly
```

The DNS server will start listening on port 53 and block requests to any domains listed in the configuration.

### Access Blocked Requests

The server logs all blocked requests, including the domain, blocked IP, and source IP. You can access these logs for monitoring.

## How It Works

1. **DNS Request Handling**: When a DNS request is made, the server checks if the domain is in the `blocked_domains` list.
2. **Blocking Requests**: If the domain is blocked, the server responds with an invalid IP (`0.0.0.0`).
3. **DNS Resolution**: Non-blocked requests are forwarded to a DNS resolver via DoH (DNS over HTTPS) for resolution.


## License

MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contributors

- **Vehbi Unal** - [@vehbiu](https://github.com/vehbiu)

