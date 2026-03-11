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
	"unsafe"
)

// ConnectArgs represents the arguments for NVMe-oF controller connection
type ConnectArgs struct {
	Traddr     string // Transport address (e.g., IP address)
	Trsvcid    string // Transport service ID (e.g., port number)
	Subsysnqn  string // Subsystem NQN
	Transport  string // Transport type (e.g., "tcp", "rdma")
	HostTraddr string // Host transport address (optional)
	HostIface  string // Host interface (optional)
}

// ConnectCtrl connects to an NVMe-oF controller using nvmf_connect_ctrl
func ConnectCtrl(args *ConnectArgs) error {
	if args == nil {
		return fmt.Errorf("connect arguments cannot be nil")
	}

	// Create nvme_root
	root := C.nvme_create_root(nil, C.DEFAULT_LOGLEVEL)
	if root == nil {
		return fmt.Errorf("failed to create nvme root")
	}
	defer C.nvme_free_tree(root)

	// Create nvme_host
	host := C.nvme_default_host(root)
	if host == nil {
		return fmt.Errorf("failed to get default host")
	}

	// Prepare C strings
	var cTraddr *C.char
	var cTrsvcid *C.char
	var cSubsysnqn *C.char
	var cTransport *C.char
	var cHostTraddr *C.char
	var cHostIface *C.char

	if args.Traddr != "" {
		cTraddr = C.CString(args.Traddr)
		defer C.free(unsafe.Pointer(cTraddr))
	}

	if args.Trsvcid != "" {
		cTrsvcid = C.CString(args.Trsvcid)
		defer C.free(unsafe.Pointer(cTrsvcid))
	}

	if args.Subsysnqn != "" {
		cSubsysnqn = C.CString(args.Subsysnqn)
		defer C.free(unsafe.Pointer(cSubsysnqn))
	}

	if args.Transport != "" {
		cTransport = C.CString(args.Transport)
		defer C.free(unsafe.Pointer(cTransport))
	}

	if args.HostTraddr != "" {
		cHostTraddr = C.CString(args.HostTraddr)
		defer C.free(unsafe.Pointer(cHostTraddr))
	}

	if args.HostIface != "" {
		cHostIface = C.CString(args.HostIface)
		defer C.free(unsafe.Pointer(cHostIface))
	}

	// Create controller object
	ctrl := C.nvme_create_ctrl(root, cSubsysnqn, cTransport, cTraddr,
		cHostTraddr, cHostIface, cTrsvcid)
	if ctrl == nil {
		return fmt.Errorf("failed to create controller object")
	}

	// Initialize default fabrics config
	var cfg C.struct_nvme_fabrics_config
	C.nvmf_default_config(&cfg)

	// Connect the controller with the host
	ret := C.nvmf_add_ctrl(host, ctrl, &cfg)
	if ret != 0 {
		return fmt.Errorf("failed to connect to NVMe-oF controller: %d", ret)
	}

	return nil
}

// DisconnectCtrl disconnects an NVMe-oF controller using nvme_disconnect_ctrl
func DisconnectCtrl(args *ConnectArgs) error {
	if args == nil {
		return fmt.Errorf("disconnect arguments cannot be nil")
	}

	// Create nvme_root
	root := C.nvme_create_root(nil, C.DEFAULT_LOGLEVEL)
	if root == nil {
		return fmt.Errorf("failed to create nvme root")
	}
	defer C.nvme_free_tree(root)

	// Scan the topology to find existing controllers
	ret := C.nvme_scan_topology(root, nil, nil)
	if ret != 0 {
		return fmt.Errorf("failed to scan NVMe topology: %d", ret)
	}

	// Get default host
	host := C.nvme_default_host(root)
	if host == nil {
		return fmt.Errorf("failed to get default host")
	}

	// Prepare C strings
	var cTraddr *C.char
	var cTrsvcid *C.char
	var cSubsysnqn *C.char
	var cTransport *C.char
	var cHostTraddr *C.char
	var cHostIface *C.char

	if args.Traddr != "" {
		cTraddr = C.CString(args.Traddr)
		defer C.free(unsafe.Pointer(cTraddr))
	}

	if args.Trsvcid != "" {
		cTrsvcid = C.CString(args.Trsvcid)
		defer C.free(unsafe.Pointer(cTrsvcid))
	}

	if args.Subsysnqn != "" {
		cSubsysnqn = C.CString(args.Subsysnqn)
		defer C.free(unsafe.Pointer(cSubsysnqn))
	}

	if args.Transport != "" {
		cTransport = C.CString(args.Transport)
		defer C.free(unsafe.Pointer(cTransport))
	}

	if args.HostTraddr != "" {
		cHostTraddr = C.CString(args.HostTraddr)
		defer C.free(unsafe.Pointer(cHostTraddr))
	}

	if args.HostIface != "" {
		cHostIface = C.CString(args.HostIface)
		defer C.free(unsafe.Pointer(cHostIface))
	}

	// Lookup the subsystem
	subsys := C.nvme_lookup_subsystem(host, nil, cSubsysnqn)
	if subsys == nil {
		return fmt.Errorf("subsystem not found: %s", args.Subsysnqn)
	}

	// Lookup the controller within the subsystem
	ctrl := C.nvme_lookup_ctrl(subsys, cTransport, cTraddr, cHostTraddr,
		cHostIface, cTrsvcid, nil)
	if ctrl == nil {
		return fmt.Errorf("controller not found")
	}

	// Disconnect the controller
	ret = C.nvme_disconnect_ctrl(ctrl)
	if ret != 0 {
		return fmt.Errorf("failed to disconnect controller: %d", ret)
	}

	return nil
}
