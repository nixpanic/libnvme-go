# AI Agent Instructions

This document provides guidance for AI agents working on the libnvme-go project.

## Project Overview

libnvme-go provides Go language bindings to the libnvme C library for managing NVMe and NVMe-oF devices on Linux.

**Key Files:**
- `connect_ctrl.go` - NVMe-oF controller connection/disconnection
- `list_subsys.go` - Subsystem listing and topology scanning
- `*_test.go` - Unit and integration tests
- `tests/Containerfile` - Container-based testing environment

## Development Workflow

### 1. Adding New Functions

When adding new libnvme bindings:

1. Check the libnvme headers to find the correct function signatures:
   - Main header: `/usr/include/libnvme.h`
   - Tree functions: `/usr/include/nvme/tree.h`
   - Fabrics functions: `/usr/include/nvme/fabrics.h`

2. Create Go wrapper types and functions in a new `.go` file

3. **Add SPDX license identifier** at the top of every source file:
   ```go
   // SPDX-License-Identifier: MIT

   package libnvme
   ```

4. Follow existing patterns:
   - Wrap C types in Go structs
   - Use `C.GoString()` for string conversion
   - Provide proper `Free()` methods for cleanup
   - Handle nil pointers gracefully

5. Create corresponding `_test.go` file with unit tests (also with SPDX header)

### 2. Testing

**IMPORTANT:** Always test changes using the containerized environment to ensure consistency.

#### Build and Run Tests in Container

```bash
# Build the container with all dependencies
podman build -t libnvme-go -f tests/Containerfile .

# Run all tests (default command)
podman run --rm libnvme-go

# Run tests in short mode (skips integration tests)
podman run --rm libnvme-go go test -v -short ./...

# Run specific test
podman run --rm libnvme-go go test -v -run TestFunctionName

# Get a shell in the container for debugging
podman run --rm -it libnvme-go /bin/bash
```

#### Test Requirements

- All new functions must have unit tests
- Tests should handle systems without NVMe devices gracefully
- Use `t.Skip()` for tests requiring special hardware/setup
- Integration tests should check for `testing.Short()` flag
- Tests must pass in the containerized environment

#### Example Test Pattern

```go
func TestNewFunction(t *testing.T) {
    // Create context
    root, err := CreateRoot()
    if err != nil {
        t.Fatalf("Failed to create root: %v", err)
    }
    defer root.Free()

    // Scan topology - may fail if no devices
    err = root.ScanTopology()
    if err != nil {
        t.Skipf("Skipping test, topology scan failed: %v", err)
    }

    // Test the function
    result := NewFunction(root)
    if result == nil {
        t.Log("No results (expected if no NVMe devices present)")
        return
    }

    // Verify result
    t.Logf("Result: %v", result)
}
```

### 3. Code Style

- **REQUIRED:** All `.go` source files must start with SPDX license identifier:
  ```go
  // SPDX-License-Identifier: MIT
  ```
- Follow Go conventions: `gofmt`, clear names, exported types start with capital letters
- Use CGO only in `*.go` files, not in tests
- Keep C interop minimal and isolated
- Add comments for all exported functions and types
- Handle errors explicitly, don't ignore return values

### 4. File Structure

Every Go source file must follow this structure:

```go
// SPDX-License-Identifier: MIT

package libnvme

/*
#cgo LDFLAGS: -lnvme
#include <stdlib.h>
#include <libnvme.h>
*/
import "C"
import (
    "fmt"
)

// Type definitions and functions...
```

### 5. Common Patterns

#### Creating Context and Scanning

```go
root, err := CreateRoot()
if err != nil {
    return err
}
defer root.Free()

err = root.ScanTopology()
if err != nil {
    return err
}
```

#### Iterating Through Topology

```go
for h := root.FirstHost(); h != nil; h = root.NextHost(h) {
    for s := h.FirstSubsystem(); s != nil; s = h.NextSubsystem(s) {
        for c := s.FirstCtrl(); c != nil; c = s.NextCtrl(c) {
            // Process controller
        }
    }
}
```

#### String Getters

```go
func (s *Subsystem) GetName() string {
    cstr := C.nvme_subsystem_get_name(s.subsys)
    if cstr == nil {
        return ""
    }
    return C.GoString(cstr)
}
```

## Checking libnvme Functions

To verify available libnvme functions:

```bash
# In container
podman run --rm libnvme-go bash -c "grep -E 'function_pattern' /usr/include/nvme/tree.h"

# Check libnvme version
podman run --rm libnvme-go bash -c "pkg-config --modversion libnvme"

# Find all header files
podman run --rm libnvme-go bash -c "find /usr/include -name '*nvme*.h'"
```

## Committing Changes

Follow the existing commit message style:

```
Brief description of changes

Detailed explanation of what was added/changed:
- Point 1
- Point 2
- Point 3

Additional context or notes.

Assisted-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

## Common Tasks

### Adding a New Subsystem Function

1. Find the function in `/usr/include/nvme/tree.h`
2. Add wrapper to `list_subsys.go`
3. Add test to `list_subsys_test.go`
4. Build and test in container
5. Commit changes

### Adding a New Connect Function

1. Find the function in `/usr/include/nvme/fabrics.h`
2. Add wrapper to `connect_ctrl.go`
3. Add test to `connect_ctrl_test.go`
4. Build and test in container
5. Commit changes

## Troubleshooting

### CGO Errors

If you see "could not determine what C.function refers to":
- Function may not exist in the installed libnvme version
- Check function name spelling and signature
- Verify the function is in libnvme headers in the container

### Test Failures

- Always verify tests pass in the container, not just locally
- Container uses Fedora with libnvme 1.16.1
- Local environment may have different versions or configurations

### Build Issues

- Ensure all C code is in `import "C"` blocks
- Don't use `unsafe` package unless absolutely necessary
- Check that all exported Go functions have proper documentation

## Resources

- [libnvme GitHub](https://github.com/linux-nvme/libnvme)
- [nvme-cli GitHub](https://github.com/linux-nvme/nvme-cli)
- [Go CGO Documentation](https://pkg.go.dev/cmd/cgo)
