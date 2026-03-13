// SPDX-License-Identifier: MIT

package libnvme

import (
	"fmt"
	"testing"
)

// TestCreateGlobalCtx tests creating and freeing a global context
func TestCreateGlobalCtx(t *testing.T) {
	ctx, err := CreateGlobalCtx()
	if err != nil {
		t.Fatalf("Failed to create global context: %v", err)
	}
	if ctx == nil {
		t.Fatal("Expected non-nil context, got nil")
	}
	if ctx.root == nil {
		t.Fatal("Expected non-nil internal root, got nil")
	}

	// Test that Free() doesn't crash
	ctx.Free()

	// Test that double-free is safe
	ctx.Free()
}

// TestScanTopology tests scanning the topology
func TestScanTopology(t *testing.T) {
	ctx, err := CreateGlobalCtx()
	if err != nil {
		t.Fatalf("Failed to create global context: %v", err)
	}
	defer ctx.Free()

	// Scan topology - this should work even without NVMe devices
	err = ctx.ScanTopology()
	if err != nil {
		t.Logf("Scan topology returned error (this may be expected): %v", err)
	}
}

// TestHostIteration tests host iteration methods
func TestHostIteration(t *testing.T) {
	ctx, err := CreateGlobalCtx()
	if err != nil {
		t.Fatalf("Failed to create global context: %v", err)
	}
	defer ctx.Free()

	err = ctx.ScanTopology()
	if err != nil {
		t.Skipf("Skipping test, topology scan failed: %v", err)
	}

	// Test FirstHost (may return nil if no devices)
	h := ctx.FirstHost()
	if h == nil {
		t.Log("No hosts found (expected if no NVMe devices present)")
		return
	}

	// Test GetHostNQN
	hostnqn := h.GetHostNQN()
	t.Logf("Found host with NQN: %s", hostnqn)

	// Test NextHost
	next := ctx.NextHost(h)
	t.Logf("NextHost returned: %v", next != nil)
}

// TestSubsystemIteration tests subsystem iteration methods
func TestSubsystemIteration(t *testing.T) {
	ctx, err := CreateGlobalCtx()
	if err != nil {
		t.Fatalf("Failed to create global context: %v", err)
	}
	defer ctx.Free()

	err = ctx.ScanTopology()
	if err != nil {
		t.Skipf("Skipping test, topology scan failed: %v", err)
	}

	h := ctx.FirstHost()
	if h == nil {
		t.Skip("No hosts found, skipping subsystem iteration test")
	}

	// Test FirstSubsystem (may return nil if no subsystems)
	s := h.FirstSubsystem()
	if s == nil {
		t.Log("No subsystems found (expected if no NVMe devices present)")
		return
	}

	// Test subsystem getters
	t.Logf("Subsystem name: %s", s.GetName())
	t.Logf("Subsystem NQN: %s", s.GetNQN())
	t.Logf("Subsystem type: %s", s.GetType())
	t.Logf("Subsystem model: %s", s.GetModel())
	t.Logf("Subsystem serial: %s", s.GetSerial())
	t.Logf("Subsystem firmware: %s", s.GetFirmwareRev())
	t.Logf("Subsystem iopolicy: %s", s.GetIOPolicy())

	// Test GetHost
	host := s.GetHost()
	if host == nil {
		t.Error("Expected non-nil host from subsystem")
	}

	// Test NextSubsystem
	next := h.NextSubsystem(s)
	t.Logf("NextSubsystem returned: %v", next != nil)
}

// TestControllerIteration tests controller iteration methods
func TestControllerIteration(t *testing.T) {
	ctx, err := CreateGlobalCtx()
	if err != nil {
		t.Fatalf("Failed to create global context: %v", err)
	}
	defer ctx.Free()

	err = ctx.ScanTopology()
	if err != nil {
		t.Skipf("Skipping test, topology scan failed: %v", err)
	}

	h := ctx.FirstHost()
	if h == nil {
		t.Skip("No hosts found, skipping controller iteration test")
	}

	s := h.FirstSubsystem()
	if s == nil {
		t.Skip("No subsystems found, skipping controller iteration test")
	}

	// Test FirstCtrl (may return nil if no controllers)
	c := s.FirstCtrl()
	if c == nil {
		t.Log("No controllers found (expected if no NVMe devices present)")
		return
	}

	// Test controller getters
	t.Logf("Controller name: %s", c.GetName())
	t.Logf("Controller transport: %s", c.GetTransport())
	t.Logf("Controller address: %s", c.GetAddress())
	t.Logf("Controller state: %s", c.GetState())

	// Test NextCtrl
	next := s.NextCtrl(c)
	t.Logf("NextCtrl returned: %v", next != nil)
}

// TestListSubsystems is an integration test that demonstrates how to list NVMe subsystems
func TestListSubsystems(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create global context
	ctx, err := CreateGlobalCtx()
	if err != nil {
		t.Fatalf("Failed to create global context: %v", err)
	}
	defer ctx.Free()

	// Scan the topology
	err = ctx.ScanTopology()
	if err != nil {
		t.Fatalf("Failed to scan topology: %v", err)
	}

	// Iterate through hosts
	foundAny := false
	for h := ctx.FirstHost(); h != nil; h = ctx.NextHost(h) {
		// Iterate through subsystems for each host
		for s := h.FirstSubsystem(); s != nil; s = h.NextSubsystem(s) {
			// Check if subsystem has any controllers
			hasCtrl := false
			for c := s.FirstCtrl(); c != nil; c = s.NextCtrl(c) {
				hasCtrl = true
				break
			}

			// Skip subsystems without controllers (like nvme list-subsys does)
			if !hasCtrl {
				continue
			}

			foundAny = true

			// Print subsystem information
			t.Logf("%s - NQN=%s", s.GetName(), s.GetNQN())
			t.Logf("  hostnqn=%s", h.GetHostNQN())

			// Print verbose information if available
			if model := s.GetModel(); model != "" {
				t.Logf("  model=%s", model)
			}
			if serial := s.GetSerial(); serial != "" {
				t.Logf("  serial=%s", serial)
			}
			if firmware := s.GetFirmwareRev(); firmware != "" {
				t.Logf("  firmware=%s", firmware)
			}
			if subsysType := s.GetType(); subsysType != "" {
				t.Logf("  type=%s", subsysType)
			}

			// Print controllers
			for c := s.FirstCtrl(); c != nil; c = s.NextCtrl(c) {
				t.Logf(" +- %s %s %s %s",
					c.GetName(),
					c.GetTransport(),
					c.GetAddress(),
					c.GetState())
			}
		}
	}

	if !foundAny {
		t.Log("No NVMe subsystems found (this may be expected if no NVMe devices are present)")
	}
}

// Example_listSubsystems shows how to list all NVMe subsystems
func Example_listSubsystems() {
	// Create global context
	ctx, err := CreateGlobalCtx()
	if err != nil {
		fmt.Printf("Failed to create global context: %v\n", err)
		return
	}
	defer ctx.Free()

	// Scan the topology
	err = ctx.ScanTopology()
	if err != nil {
		fmt.Printf("Failed to scan topology: %v\n", err)
		return
	}

	// Iterate through hosts
	for h := ctx.FirstHost(); h != nil; h = ctx.NextHost(h) {
		// Iterate through subsystems for each host
		for s := h.FirstSubsystem(); s != nil; s = h.NextSubsystem(s) {
			// Check if subsystem has any controllers
			hasCtrl := false
			for c := s.FirstCtrl(); c != nil; c = s.NextCtrl(c) {
				hasCtrl = true
				break
			}

			// Skip subsystems without controllers
			if !hasCtrl {
				continue
			}

			// Print subsystem information
			fmt.Printf("%s - NQN=%s\n", s.GetName(), s.GetNQN())
			fmt.Printf("  hostnqn=%s\n", h.GetHostNQN())

			// Print controllers
			for c := s.FirstCtrl(); c != nil; c = s.NextCtrl(c) {
				fmt.Printf(" +- %s %s %s %s\n",
					c.GetName(),
					c.GetTransport(),
					c.GetAddress(),
					c.GetState())
			}
		}
	}
}
