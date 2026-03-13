// SPDX-License-Identifier: MIT

package libnvme

/*
#cgo LDFLAGS: -lnvme
#include <stdlib.h>
#include <libnvme.h>
*/
import "C"

// Namespace represents an nvme_ns_t
type Namespace struct {
	ns C.nvme_ns_t
}

// FirstNs returns the first namespace for the controller
func (c *Ctrl) FirstNs() *Namespace {
	n := C.nvme_ctrl_first_ns(c.ctrl)
	if n == nil {
		return nil
	}
	return &Namespace{ns: n}
}

// NextNs returns the next namespace after the given namespace
func (c *Ctrl) NextNs(n *Namespace) *Namespace {
	next := C.nvme_ctrl_next_ns(c.ctrl, n.ns)
	if next == nil {
		return nil
	}
	return &Namespace{ns: next}
}

// FirstNs returns the first namespace for the subsystem
func (s *Subsystem) FirstNs() *Namespace {
	n := C.nvme_subsystem_first_ns(s.subsys)
	if n == nil {
		return nil
	}
	return &Namespace{ns: n}
}

// NextNs returns the next namespace after the given namespace
func (s *Subsystem) NextNs(n *Namespace) *Namespace {
	next := C.nvme_subsystem_next_ns(s.subsys, n.ns)
	if next == nil {
		return nil
	}
	return &Namespace{ns: next}
}

// GetName returns the namespace name (e.g., "nvme0n1")
func (n *Namespace) GetName() string {
	cstr := C.nvme_ns_get_name(n.ns)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetGenericName returns the generic namespace chardev name
func (n *Namespace) GetGenericName() string {
	cstr := C.nvme_ns_get_generic_name(n.ns)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetNsid returns the namespace ID
func (n *Namespace) GetNsid() int {
	return int(C.nvme_ns_get_nsid(n.ns))
}

// GetLbaSize returns the LBA size in bytes
func (n *Namespace) GetLbaSize() int {
	return int(C.nvme_ns_get_lba_size(n.ns))
}

// GetMetaSize returns the metadata size in bytes
func (n *Namespace) GetMetaSize() int {
	return int(C.nvme_ns_get_meta_size(n.ns))
}

// GetLbaCount returns the total number of LBAs
func (n *Namespace) GetLbaCount() uint64 {
	return uint64(C.nvme_ns_get_lba_count(n.ns))
}

// GetLbaUtil returns the number of utilized LBAs
func (n *Namespace) GetLbaUtil() uint64 {
	return uint64(C.nvme_ns_get_lba_util(n.ns))
}

// GetFirmware returns the namespace firmware revision
func (n *Namespace) GetFirmware() string {
	cstr := C.nvme_ns_get_firmware(n.ns)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetSerial returns the namespace serial number
func (n *Namespace) GetSerial() string {
	cstr := C.nvme_ns_get_serial(n.ns)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetModel returns the namespace model
func (n *Namespace) GetModel() string {
	cstr := C.nvme_ns_get_model(n.ns)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// GetSubsystem returns the subsystem associated with the namespace
func (n *Namespace) GetSubsystem() *Subsystem {
	s := C.nvme_ns_get_subsystem(n.ns)
	if s == nil {
		return nil
	}
	return &Subsystem{subsys: s}
}

// GetCtrl returns the controller associated with the namespace
func (n *Namespace) GetCtrl() *Ctrl {
	c := C.nvme_ns_get_ctrl(n.ns)
	if c == nil {
		return nil
	}
	return &Ctrl{ctrl: c}
}
