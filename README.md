# libnvme-go

Go bindings for libnvme - NVMe management library for Linux.

## Overview

libnvme-go provides Go language bindings to the [libnvme](https://github.com/linux-nvme/libnvme) C library, enabling Go applications to manage NVMe and NVMe-oF (NVMe over Fabrics) devices.

## Features

- Connect to NVMe-oF controllers
- Disconnect from NVMe-oF controllers
- Support for multiple transport types (TCP, RDMA)
- Comprehensive error handling
- Type-safe Go API wrapping libnvme C functions

## Requirements

- Go 1.25 or later
- libnvme library and development headers
- GCC (for CGO compilation)

### Installing Dependencies

**Fedora/RHEL:**
```bash
sudo dnf install libnvme-devel gcc golang
```

**Debian/Ubuntu:**
```bash
sudo apt install libnvme-dev gcc golang
```

## Installation

```bash
go get github.com/nixpanic/libnvme-go
```

## Usage

### Connecting to an NVMe-oF Controller

```go
package main

import (
    "fmt"
    "log"

    "github.com/nixpanic/libnvme-go"
)

func main() {
    args := &libnvme.ConnectArgs{
        Traddr:    "192.168.1.100",         // Target IP address
        Trsvcid:   "4420",                  // Port number
        Subsysnqn: "nqn.2014-08.org.nvmexpress:uuid:my-subsystem",
        Transport: "tcp",                   // or "rdma"
    }

    err := libnvme.ConnectCtrl(args)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }

    fmt.Println("Successfully connected to NVMe-oF controller")
}
```

### Disconnecting from an NVMe-oF Controller

```go
err := libnvme.DisconnectCtrl(args)
if err != nil {
    log.Fatalf("Failed to disconnect: %v", err)
}

fmt.Println("Successfully disconnected from NVMe-oF controller")
```

## Testing

### Unit Tests

Run the unit tests:
```bash
go test -v
```

Run tests in short mode (skips integration tests):
```bash
go test -v -short
```

### Container-based Testing

Build and test in a containerized environment:
```bash
# Build the container
podman build -t libnvme-go -f tests/Containerfile .

# Run tests in the container
podman run --rm libnvme-go go test -v
```

## Development

### Building

```bash
go build
```

### Code Formatting

```bash
go fmt ./...
```

## CI/CD

This project uses GitHub Actions for continuous integration. On every pull request and push to main:
- Container is built with all dependencies
- Unit tests are executed
- Code is compiled to verify build success

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Acknowledgments

- This library wraps the [libnvme](https://github.com/linux-nvme/libnvme) C library
- Code generated with AI assistance from Claude Code (Anthropic)

## Resources

- [NVMe Specifications](https://nvmexpress.org/specifications/)
- [libnvme Documentation](https://github.com/linux-nvme/libnvme)
- [NVMe-oF (NVMe over Fabrics)](https://nvmexpress.org/nvme-over-fabrics/)
