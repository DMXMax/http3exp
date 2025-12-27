# HTTP/3 and QUIC Examples

Educational demonstration project showcasing QUIC protocol and HTTP/3 implementation patterns. This code accompanies the QUIC overview talk and provides hands-on examples of different server and client configurations.

## Prerequisites

- **Go 1.25+** - [Download](https://go.dev/dl/)
- **OpenSSL 3.x** - For certificate generation (included with macOS via Homebrew)
- **IDE** (optional) - VSCode recommended (includes debug configurations in `.vscode/`)

## Quick Start

### 1. Generate Certificates

QUIC and HTTP/3 require TLS certificates. Generate self-signed certificates for testing:

```bash
cd certs
./createcerts.sh
cd ..
```

This creates certificates valid for 7 days. **Never use these certificates in production.**

### 2. Build the Project

```bash
go build -o http3test .
```

### 3. Run Examples

**Option A: Server only (access via browser)**
```bash
./http3test -s 1
# Visit https://localhost:8080 in your browser
```

**Option B: Server + Client**
```bash
./http3test -s 2 -c 1 -n 3
```

## Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-s` | Server type to run (0-3) | 0 |
| `-c` | Client type to run (0-3), use -1 for no client | -1 |
| `-n` | Number of times to run the client | 1 |

**Note:** When `-c -1` (default), the server runs in blocking mode until interrupted.

## Server Types

### Server 0: Basic HTTPS (TCP/TLS)
```bash
./http3test -s 0
```
- Traditional HTTPS server over TCP
- Uses standard `http.Server` with TLS configuration
- Good baseline for comparing with HTTP/3

### Server 1: HTTP/3 with TCP Fallback (Recommended)
```bash
./http3test -s 1
```
- **Production-ready pattern** - Serves both HTTP/2 and HTTP/3
- Automatically advertises HTTP/3 via Alt-Svc header at TCP layer
- Browser-compatible (use https://localhost:8080)
- Serves protocol-specific favicons:
  - `two.png` for HTTP/2 connections
  - `three.png` for HTTP/3 connections

### Server 2: Pure HTTP/3 with QUIC
```bash
./http3test -s 2 -c 1
```
- Direct QUIC/HTTP/3 implementation without TCP fallback
- Enables 0-RTT (Zero Round Trip Time) for faster reconnections
- Returns JSON responses
- Requires HTTP/3-capable clients

### Server 3: Raw QUIC Echo Server
```bash
./http3test -s 3 -c 3
```
- Low-level QUIC stream handling (no HTTP/3 wrapping)
- Echoes back received data
- Demonstrates raw QUIC protocol usage
- Must use Client 3

## Client Types

### Client 0: Standard HTTPS Client
```bash
./http3test -s 0 -c 0
```
- Traditional HTTPS client over TCP/TLS
- Uses `http.Transport` with custom TLS config
- Tests baseline HTTP/2 connections

### Client 1: HTTP/3 Client
```bash
./http3test -s 2 -c 1 -n 5
```
- HTTP/3 round-tripper using QUIC transport
- Session caching for connection reuse
- Supports 0-RTT for subsequent requests
- Compatible with Servers 1 and 2

### Client 2: HTTP/3 with 0-RTT
```bash
./http3test -s 2 -c 2
```
- Advanced HTTP/3 client with zero round-trip optimization
- Uses `http3.MethodGet0RTT` for faster initial requests
- Demonstrates performance benefits of 0-RTT

### Client 3: Raw QUIC Client
```bash
./http3test -s 3 -c 3
```
- Low-level QUIC stream client
- Opens streams and exchanges raw data
- Matches Server 3's protocol
- No HTTP/3 layer

## Server/Client Compatibility Matrix

| Server | Client 0 | Client 1 | Client 2 | Client 3 | Browser |
|--------|----------|----------|----------|----------|---------|
| Server 0 (HTTPS) | ✅ | ❌ | ❌ | ❌ | ✅ |
| Server 1 (HTTP/3+fallback) | ✅ | ✅ | ✅ | ❌ | ✅ |
| Server 2 (HTTP/3) | ❌ | ✅ | ✅ | ❌ | ⚠️* |
| Server 3 (Raw QUIC) | ❌ | ❌ | ❌ | ✅ | ❌ |

*Browser support requires Alt-Svc discovery which Server 2 doesn't provide.

## Common Usage Examples

### Compare HTTP/2 vs HTTP/3 Performance
```bash
# Terminal 1: Start HTTP/3 server with fallback
./http3test -s 1

# Terminal 2: Test with multiple clients
./http3test -s 1 -c 1 -n 10
```

### Test 0-RTT Performance
```bash
# First run establishes session
./http3test -s 2 -c 2 -n 1

# Subsequent runs use 0-RTT (faster)
./http3test -s 2 -c 2 -n 5
```

### Browser Testing
```bash
# Start server
./http3test -s 1

# Open browser to:
# https://localhost:8080
# Check DevTools → Network → Protocol column for "h3"
```

### Raw QUIC Stream Communication
```bash
./http3test -s 3 -c 3 -n 3
```

## Debugging with Wireshark

The application logs TLS keys to `~/.ssl-key.log` for protocol analysis:

1. Start Wireshark and capture on loopback interface
2. Set SSL key log: Preferences → Protocols → TLS → (Pre)-Master-Secret log filename → `~/.ssl-key.log`
3. Run the application
4. Filter: `quic` or `http3`

## Project Structure

```
.
├── main.go              # Entry point with CLI parsing
├── server/
│   ├── server.go        # 4 server implementations
│   ├── two.png          # HTTP/2 favicon
│   └── three.png        # HTTP/3 favicon
├── client/
│   └── client.go        # 4 client implementations
├── util/
│   └── util.go          # Certificate path helpers
├── certs/
│   ├── createcerts.sh   # Certificate generation script
│   └── *.pem, *.key     # Generated certificates (gitignored)
└── README.md
```

## Troubleshooting

**Certificate errors:**
```bash
cd certs && ./createcerts.sh && cd ..
```

**Port already in use:**
```bash
# Kill existing process
lsof -ti:8080 | xargs kill -9
```

**Browser doesn't use HTTP/3:**
- Ensure using Server 1 (has Alt-Svc header)
- Clear browser cache
- Chrome: Visit chrome://net-internals/#http3 to verify

**OpenSSL config error:**
- The script automatically sets `OPENSSL_CONF` for Homebrew installations
- If issues persist, ensure OpenSSL 3.x is installed: `brew install openssl@3`

## Learn More

- [QUIC Protocol](https://www.rfc-editor.org/rfc/rfc9000.html)
- [HTTP/3 Specification](https://www.rfc-editor.org/rfc/rfc9114.html)
- [quic-go Library](https://github.com/quic-go/quic-go)
- [quic-go Documentation](https://quic-go.net/)
