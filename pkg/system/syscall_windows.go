package system

import (
	"unsafe"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

const (
	SeTakeOwnershipPrivilege = "SeTakeOwnershipPrivilege"
)

const (
	LabniAdministratorSidString = "S-1-8-73-2-1"
	LabniUserSidString          = "S-1-8-73-2-2"
)

var (
	ntuserApiset      = windows.NewLazyDLL("ext-ms-win-ntuser-window-l1-1-0")
	procGetVersionExW = modkernel32.NewProc("GetVersionExW")
)

// TODO: use golang.org/x/sys/windows.OsVersionInfoEx (needs OSVersionInfoSize to be exported)
type osVersionInfoEx struct {
	OSVersionInfoSize uint32
	MajorVersion      uint32
	MinorVersion      uint32
	BuildNumber       uint32
	PlatformID        uint32
	CSDVersion        [128]uint16
	ServicePackMajor  uint16
	ServicePackMinor  uint16
	SuiteMask         uint16
	ProductType       byte
	Reserve           byte
}

// IsWindowsClient returns true if the SKU is client. It returns false on
// Windows server, or if an error occurred when making the GetVersionExW
// syscall.
func IsWindowsClient() bool {
	osviex := &osVersionInfoEx{OSVersionInfoSize: 284}
	r1, _, err := procGetVersionExW.Call(uintptr(unsafe.Pointer(osviex)))
	if r1 == 0 {
		logrus.WithError(err).Warn("GetVersionExW failed - assuming server SKU")
		return false
	}
	// VER_NT_WORKSTATION
	const verNTWorkstation = 0x00000001 // VER_NT_WORKSTATION
	return osviex.ProductType == verNTWorkstation
}

// HasWin32KSupport determines whether Labni(s) that depend on win32k can
// run on this machine. Win32k is the driver used to implement windowing.
func HasWin32KSupport() bool {
	// For now, check for ntuser API support on the host. In the future, a host
	// may support win32k in Labni(s) even if the host does not support ntuser
	// APIs.
	return ntuserApiset.Load() == nil
}
