package ciutil

import (
	"fmt"
	"runtime"
)

// HostInfoInput allows different tools to specify the different possible values for their host
// information, such as when downloading a release artifact. This is necessary to provide because
// different tools use different values for their host info, e.g. some use "x86_64" and some use
// "amd64", some additionally specify an OS value and some do not, etc.
type HostInfoInput struct {
	// The name a tool uses for Linux OSes (e.g. "linux", "unknown", etc.)
	OSLinux string
	// The name a tool uses for the Linux kernel (e.g. "linux", "linux-gnu", etc.)
	KernelLinux string
	// The name a tool uses for macOS (e.g. "apple")
	OSMacOS string
	// The name a tool uses for the macOS kernel (e.g. "darwin")
	KernelMacOS string
	// The name a tool uses for AMD64 architectures (e.g. "x86_64", "amd64", etc.)
	ArchAMD64 string
	// The name a tool uses for ARM64 architectures (e.g. "aarch64", "arm64", etc.)
	ArchARM64 string
}

// HostInfo holds the final host information values to be used in e.g. release downloads of a tool.
type HostInfo struct {
	OS     string
	Arch   string
	Kernel string
}

// GetHostInfo returns a populated [HostInfo], based on the provided [HostInfoInput] mappings.
func GetHostInfo(i HostInfoInput) (HostInfo, error) {
	var out HostInfo

	switch runtime.GOOS {
	case "linux":
		out.OS = i.OSLinux
		out.Kernel = i.KernelLinux
	case "darwin":
		out.OS = i.OSMacOS
		out.Kernel = i.KernelMacOS
	default:
		return HostInfo{}, fmt.Errorf("unsupported operating system/kernel '%s'", runtime.GOOS)
	}

	switch runtime.GOARCH {
	case "amd64":
		out.Arch = i.ArchAMD64
	case "arm64":
		out.Arch = i.ArchARM64
	default:
		return HostInfo{}, fmt.Errorf("unsupported CPU architecture '%s'", runtime.GOARCH)
	}

	return out, nil
}
