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

// Root represents an nvme_root_t
type Root struct {
	root C.nvme_root_t
}

// Host represents an nvme_host_t
type Host struct {
	host C.nvme_host_t
}

// Subsystem represents an nvme_subsystem_t
type Subsystem struct {
	subsys C.nvme_subsystem_t
}

// Ctrl represents an nvme_ctrl_t
type Ctrl struct {
	ctrl C.nvme_ctrl_t
}

// CreateRoot creates a new NVMe root context
func CreateRoot() (*Root, error) {
	root := C.nvme_create_root(nil, C.DEFAULT_LOGLEVEL)
	if root == nil {
		return nil, fmt.Errorf("failed to create root context")
	}
	return &Root{root: root}, nil
}

// Free frees the root context and associated tree
func (r *Root) Free() {
	if r.root != nil {
		C.nvme_free_tree(r.root)
		r.root = nil
	}
}

// ScanTopology scans the NVMe topology
func (r *Root) ScanTopology() error {
	ret := C.nvme_scan_topology(r.root, nil, nil)
	if ret != 0 {
		return fmt.Errorf("failed to scan topology: error code %d", ret)
	}
	return nil
}

// FirstHost returns the first host in the root context
func (r *Root) FirstHost() *Host {
	h := C.nvme_first_host(r.root)
	if h == nil {
		return nil
	}
	return &Host{host: h}
}

// NextHost returns the next host after the given host
func (r *Root) NextHost(h *Host) *Host {
	next := C.nvme_next_host(r.root, h.host)
	if next == nil {
		return nil
	}
	return &Host{host: next}
}

// FirstSubsystem returns the first subsystem for the host
func (h *Host) FirstSubsystem() *Subsystem {
	s := C.nvme_first_subsystem(h.host)
	if s == nil {
		return nil
	}
	return &Subsystem{subsys: s}
}

// NextSubsystem returns the next subsystem after the given subsystem
func (h *Host) NextSubsystem(s *Subsystem) *Subsystem {
	next := C.nvme_next_subsystem(h.host, s.subsys)
	if next == nil {
		return nil
	}
	return &Subsystem{subsys: next}
}

// GetHostNQN returns the host NQN
func (h *Host) GetHostNQN() string {
	cstr := C.nvme_host_get_hostnqn(h.host)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// FirstCtrl returns the first controller in the subsystem
func (s *Subsystem) FirstCtrl() *Ctrl {
	c := C.nvme_subsystem_first_ctrl(s.subsys)
	if c == nil {
		return nil
	}
	return &Ctrl{ctrl: c}
}

// NextCtrl returns the next controller after the given controller
func (s *Subsystem) NextCtrl(c *Ctrl) *Ctrl {
	next := C.nvme_subsystem_next_ctrl(s.subsys, c.ctrl)
	if next == nil {
		return nil
	}
	return &Ctrl{ctrl: next}
}

// GetName returns the subsystem name
func (s *Subsystem) GetName() string {
	cstr := C.nvme_subsystem_get_name(s.subsys)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetNQN returns the subsystem NQN
func (s *Subsystem) GetNQN() string {
	cstr := C.nvme_subsystem_get_nqn(s.subsys)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetHost returns the host associated with the subsystem
func (s *Subsystem) GetHost() *Host {
	h := C.nvme_subsystem_get_host(s.subsys)
	if h == nil {
		return nil
	}
	return &Host{host: h}
}

// GetIOPolicy returns the subsystem I/O policy
func (s *Subsystem) GetIOPolicy() string {
	cstr := C.nvme_subsystem_get_iopolicy(s.subsys)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetModel returns the subsystem model
func (s *Subsystem) GetModel() string {
	cstr := C.nvme_subsystem_get_model(s.subsys)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetSerial returns the subsystem serial number
func (s *Subsystem) GetSerial() string {
	cstr := C.nvme_subsystem_get_serial(s.subsys)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetFirmwareRev returns the subsystem firmware revision
func (s *Subsystem) GetFirmwareRev() string {
	cstr := C.nvme_subsystem_get_fw_rev(s.subsys)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetType returns the subsystem type
func (s *Subsystem) GetType() string {
	cstr := C.nvme_subsystem_get_type(s.subsys)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetName returns the controller name
func (c *Ctrl) GetName() string {
	cstr := C.nvme_ctrl_get_name(c.ctrl)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetTransport returns the controller transport type
func (c *Ctrl) GetTransport() string {
	cstr := C.nvme_ctrl_get_transport(c.ctrl)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetAddress returns the controller address
func (c *Ctrl) GetAddress() string {
	cstr := C.nvme_ctrl_get_address(c.ctrl)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetState returns the controller state
func (c *Ctrl) GetState() string {
	cstr := C.nvme_ctrl_get_state(c.ctrl)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// Backward compatibility aliases
type GlobalCtx = Root

// CreateGlobalCtx is an alias for CreateRoot for backward compatibility
var CreateGlobalCtx = CreateRoot
