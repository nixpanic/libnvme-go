// SPDX-License-Identifier: MIT

package libnvme

import (
	"testing"
)

// TestCtrlNamespaceIteration tests namespace iteration from controller
func TestCtrlNamespaceIteration(t *testing.T) {
	root, err := CreateRoot()
	if err != nil {
		t.Fatalf("Failed to create root: %v", err)
	}
	defer root.Free()

	err = root.ScanTopology()
	if err != nil {
		t.Skipf("Skipping test, topology scan failed: %v", err)
	}

	h := root.FirstHost()
	if h == nil {
		t.Skip("No hosts found, skipping controller namespace iteration test")
	}

	s := h.FirstSubsystem()
	if s == nil {
		t.Skip("No subsystems found, skipping controller namespace iteration test")
	}

	c := s.FirstCtrl()
	if c == nil {
		t.Skip("No controllers found, skipping controller namespace iteration test")
	}

	// Test FirstNs (may return nil if no namespaces)
	n := c.FirstNs()
	if n == nil {
		t.Log("No namespaces found (expected if no NVMe namespaces present)")
		return
	}

	// Test namespace getters
	t.Logf("Namespace name: %s", n.GetName())
	t.Logf("Namespace generic name: %s", n.GetGenericName())
	t.Logf("Namespace NSID: %d", n.GetNsid())
	t.Logf("Namespace LBA size: %d", n.GetLbaSize())
	t.Logf("Namespace meta size: %d", n.GetMetaSize())
	t.Logf("Namespace LBA count: %d", n.GetLbaCount())
	t.Logf("Namespace LBA util: %d", n.GetLbaUtil())
	t.Logf("Namespace firmware: %s", n.GetFirmware())
	t.Logf("Namespace serial: %s", n.GetSerial())
	t.Logf("Namespace model: %s", n.GetModel())

	// Test GetSubsystem
	subsys := n.GetSubsystem()
	if subsys == nil {
		t.Error("Expected non-nil subsystem from namespace")
	}

	// Test GetCtrl
	ctrl := n.GetCtrl()
	if ctrl == nil {
		t.Error("Expected non-nil controller from namespace")
	}

	// Test NextNs
	next := c.NextNs(n)
	t.Logf("NextNs returned: %v", next != nil)
}

// TestSubsystemNamespaceIteration tests namespace iteration from subsystem
func TestSubsystemNamespaceIteration(t *testing.T) {
	root, err := CreateRoot()
	if err != nil {
		t.Fatalf("Failed to create root: %v", err)
	}
	defer root.Free()

	err = root.ScanTopology()
	if err != nil {
		t.Skipf("Skipping test, topology scan failed: %v", err)
	}

	h := root.FirstHost()
	if h == nil {
		t.Skip("No hosts found, skipping subsystem namespace iteration test")
	}

	s := h.FirstSubsystem()
	if s == nil {
		t.Skip("No subsystems found, skipping subsystem namespace iteration test")
	}

	// Test FirstNs (may return nil if no namespaces)
	n := s.FirstNs()
	if n == nil {
		t.Log("No namespaces found (expected if no NVMe namespaces present)")
		return
	}

	t.Logf("Found namespace: %s (NSID: %d)", n.GetName(), n.GetNsid())

	// Test NextNs
	next := s.NextNs(n)
	t.Logf("NextNs returned: %v", next != nil)
}

// TestNamespaceProperties tests all namespace property getters
func TestNamespaceProperties(t *testing.T) {
	root, err := CreateRoot()
	if err != nil {
		t.Fatalf("Failed to create root: %v", err)
	}
	defer root.Free()

	err = root.ScanTopology()
	if err != nil {
		t.Skipf("Skipping test, topology scan failed: %v", err)
	}

	h := root.FirstHost()
	if h == nil {
		t.Skip("No hosts found, skipping namespace properties test")
	}

	s := h.FirstSubsystem()
	if s == nil {
		t.Skip("No subsystems found, skipping namespace properties test")
	}

	c := s.FirstCtrl()
	if c == nil {
		t.Skip("No controllers found, skipping namespace properties test")
	}

	n := c.FirstNs()
	if n == nil {
		t.Skip("No namespaces found, skipping namespace properties test")
	}

	// Test all string getters return valid strings (may be empty)
	name := n.GetName()
	if name == "" {
		t.Error("Expected non-empty namespace name")
	}

	// Test numeric getters
	nsid := n.GetNsid()
	if nsid <= 0 {
		t.Errorf("Expected positive NSID, got %d", nsid)
	}

	lbaSize := n.GetLbaSize()
	if lbaSize <= 0 {
		t.Errorf("Expected positive LBA size, got %d", lbaSize)
	}

	lbaCount := n.GetLbaCount()
	if lbaCount == 0 {
		t.Error("Expected non-zero LBA count")
	}

	// Meta size can be 0, so just check it doesn't panic
	_ = n.GetMetaSize()

	// LBA util can be 0 for unused namespaces
	_ = n.GetLbaUtil()

	// Test relationship getters
	subsys := n.GetSubsystem()
	if subsys == nil {
		t.Error("Expected non-nil subsystem")
	} else {
		t.Logf("Namespace belongs to subsystem: %s", subsys.GetName())
	}

	ctrl := n.GetCtrl()
	if ctrl == nil {
		t.Error("Expected non-nil controller")
	} else {
		t.Logf("Namespace belongs to controller: %s", ctrl.GetName())
	}
}

// TestListNamespaces is an integration test that demonstrates listing all namespaces
func TestListNamespaces(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	root, err := CreateRoot()
	if err != nil {
		t.Fatalf("Failed to create root: %v", err)
	}
	defer root.Free()

	err = root.ScanTopology()
	if err != nil {
		t.Fatalf("Failed to scan topology: %v", err)
	}

	foundAny := false
	for h := root.FirstHost(); h != nil; h = root.NextHost(h) {
		for s := h.FirstSubsystem(); s != nil; s = h.NextSubsystem(s) {
			for c := s.FirstCtrl(); c != nil; c = s.NextCtrl(c) {
				for n := c.FirstNs(); n != nil; n = c.NextNs(n) {
					foundAny = true

					// Print namespace information like "nvme list-ns"
					t.Logf("%-12s %-12s %-20s %-40s NSID:%-4d Size:%-12d LBA:%-8d",
						n.GetName(),
						n.GetGenericName(),
						n.GetSerial(),
						n.GetModel(),
						n.GetNsid(),
						n.GetLbaCount()*uint64(n.GetLbaSize()),
						n.GetLbaSize())
				}
			}
		}
	}

	if !foundAny {
		t.Log("No NVMe namespaces found (this may be expected if no NVMe devices are present)")
	}
}

// Example_listNamespaces shows how to list all NVMe namespaces
func Example_listNamespaces() {
	// Create root context
	root, err := CreateRoot()
	if err != nil {
		return
	}
	defer root.Free()

	// Scan topology
	err = root.ScanTopology()
	if err != nil {
		return
	}

	// Iterate through all namespaces
	for h := root.FirstHost(); h != nil; h = root.NextHost(h) {
		for s := h.FirstSubsystem(); s != nil; s = h.NextSubsystem(s) {
			for c := s.FirstCtrl(); c != nil; c = s.NextCtrl(c) {
				for n := c.FirstNs(); n != nil; n = c.NextNs(n) {
					// Calculate total size in bytes
					totalSize := n.GetLbaCount() * uint64(n.GetLbaSize())

					// Print namespace information
					println("Namespace:", n.GetName())
					println("  NSID:", n.GetNsid())
					println("  Size:", totalSize, "bytes")
					println("  LBA Size:", n.GetLbaSize(), "bytes")
					println("  Model:", n.GetModel())
					println("  Serial:", n.GetSerial())
					println("  Firmware:", n.GetFirmware())
				}
			}
		}
	}
}
