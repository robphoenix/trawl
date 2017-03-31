package main

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
)

const decBase = 10

// Interface provides the information for a device interface
type Interface struct {
	HardwareAddr string
	IPv4Addr     string
	IPv4Mask     string
	IPv4Network  string
	IPv6Addr     string
	MTU          int
	Name         string
}

// New instantiates an Interface object for the passed in net.Interface type
// representing a device interface
func New(iface net.Interface) (i *Interface, err error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return i, err
	}

	// we can't rely on the order of the addresses in the addrs array
	ipv4, ipv6 := extractAddrs(addrs)

	// if we have an IPv6 only interface
	if ipv4 == nil {
		return &Interface{
			Name:     iface.Name,
			IPv6Addr: safeIPNetToString(ipv6),
		}, nil
	}

	return &Interface{
		HardwareAddr: iface.HardwareAddr.String(),
		IPv4Addr:     ipv4.IP.String(),
		IPv4Mask:     toDottedDec(ipv4.Mask),
		IPv4Network:  maskedIPString(ipv4),
		IPv6Addr:     safeIPNetToString(ipv6),
		MTU:          iface.MTU,
		Name:         iface.Name,
	}, nil
}

func extractAddrs(addrs []net.Addr) (ipv4, ipv6 *net.IPNet) {
	for _, addr := range addrs {
		switch ipnet := addr.(type) {
		case *net.IPNet:
			if ip := ipnet.IP.To4(); ip != nil {
				ipv4 = ipnet
			} else {
				ipv6 = ipnet
			}
		}
	}
	return
}

func toDottedDec(mask net.IPMask) string {
	parts := make([]string, len(mask))
	for i, part := range mask {
		parts[i] = strconv.FormatUint(uint64(part), decBase)
	}
	return strings.Join(parts, ".")
}

func maskedIPString(ipnet *net.IPNet) string {
	ip := ipnet.IP
	mask := ipnet.Mask
	maskOnes, _ := mask.Size()
	suffix := "/" + strconv.FormatInt(int64(maskOnes), decBase)
	return ip.Mask(mask).String() + suffix
}

func (iface *Interface) String() string {
	ifaceString := "%-10s  %-15s  %-15s  %-18s  %4d  %17s  %s"
	if runtime.GOOS == "windows" {
		ifaceString = "%-35s  %-15s  %-15s  %-18s  %4d  %17s  %s"
	}
	return fmt.Sprintf(
		ifaceString,
		iface.Name,
		iface.IPv4Addr,
		iface.IPv4Mask,
		iface.IPv4Network,
		iface.MTU,
		iface.HardwareAddr,
		iface.IPv6Addr,
	)
}

func safeIPNetToString(ipnet *net.IPNet) string {
	if ipnet == nil {
		return ""
	}
	return ipnet.String()
}
