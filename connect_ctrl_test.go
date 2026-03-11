// SPDX-License-Identifier: MIT

package libnvme

import (
	"testing"
)

func TestConnectCtrl_NilArgs(t *testing.T) {
	err := ConnectCtrl(nil)
	if err == nil {
		t.Error("Expected error when args is nil, got nil")
	}
	expectedMsg := "connect arguments cannot be nil"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestDisconnectCtrl_NilArgs(t *testing.T) {
	err := DisconnectCtrl(nil)
	if err == nil {
		t.Error("Expected error when args is nil, got nil")
	}
	expectedMsg := "disconnect arguments cannot be nil"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestConnectCtrl_EmptyArgs(t *testing.T) {
	args := &ConnectArgs{}
	err := ConnectCtrl(args)
	if err == nil {
		t.Error("Expected error when args are empty, got nil")
	}
}

func TestDisconnectCtrl_EmptyArgs(t *testing.T) {
	args := &ConnectArgs{}
	err := DisconnectCtrl(args)
	if err == nil {
		t.Error("Expected error when args are empty, got nil")
	}
}

// Integration test that requires actual NVMe-oF target
// Set NVME_TEST_TARGET environment variable to enable this test
// Example: NVME_TEST_TARGET=1 go test -v
func TestConnectDisconnect_Integration(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires:
	// 1. An actual NVMe-oF target available
	// 2. Root privileges or CAP_SYS_ADMIN
	// 3. Proper network connectivity to the target
	t.Skip("Integration test requires actual NVMe-oF target setup")

	// Example test that would run with a real target:
	/*
		args := &ConnectArgs{
			Traddr:    "192.168.1.100",
			Trsvcid:   "4420",
			Subsysnqn: "nqn.2014-08.org.nvmexpress:uuid:test-subsystem",
			Transport: "tcp",
		}

		// Test Connect
		err := ConnectCtrl(args)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}

		// Test Disconnect
		err = DisconnectCtrl(args)
		if err != nil {
			t.Fatalf("Failed to disconnect: %v", err)
		}
	*/
}

func TestConnectArgs_Fields(t *testing.T) {
	args := &ConnectArgs{
		Traddr:     "192.168.1.100",
		Trsvcid:    "4420",
		Subsysnqn:  "nqn.2014-08.org.nvmexpress:uuid:test",
		Transport:  "tcp",
		HostTraddr: "192.168.1.1",
		HostIface:  "eth0",
	}

	if args.Traddr != "192.168.1.100" {
		t.Errorf("Expected Traddr '192.168.1.100', got '%s'", args.Traddr)
	}
	if args.Trsvcid != "4420" {
		t.Errorf("Expected Trsvcid '4420', got '%s'", args.Trsvcid)
	}
	if args.Subsysnqn != "nqn.2014-08.org.nvmexpress:uuid:test" {
		t.Errorf("Expected Subsysnqn 'nqn.2014-08.org.nvmexpress:uuid:test', got '%s'", args.Subsysnqn)
	}
	if args.Transport != "tcp" {
		t.Errorf("Expected Transport 'tcp', got '%s'", args.Transport)
	}
	if args.HostTraddr != "192.168.1.1" {
		t.Errorf("Expected HostTraddr '192.168.1.1', got '%s'", args.HostTraddr)
	}
	if args.HostIface != "eth0" {
		t.Errorf("Expected HostIface 'eth0', got '%s'", args.HostIface)
	}
}
