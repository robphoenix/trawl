package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

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
	if ipv4 == "" {
		return &Interface{
			Name:     iface.Name,
			IPv6Addr: ipv6,
		}, nil
	}

	// get IPv4 address and network
	ipv4Addr, ipv4Network, err := net.ParseCIDR(ipv4)
	if err != nil {
		log.Fatal(err)
	}

	return &Interface{
		HardwareAddr: iface.HardwareAddr.String(),
		IPv4Addr:     ipv4Addr.String(),
		IPv4Mask:     toDottedDec(ipv4Network.Mask),
		IPv4Network:  ipv4Network.String(),
		IPv6Addr:     ipv6,
		MTU:          iface.MTU,
		Name:         iface.Name,
	}, nil
}

func extractAddrs(addrs []net.Addr) (ipv4, ipv6 string) {
	for _, addr := range addrs {
		a := addr.String()
		switch {
		case strings.Contains(a, ":"):
			ipv6 = a
		case strings.Contains(a, "."):
			ipv4 = a
		}
	}
	return
}

func toDottedDec(mask net.IPMask) string {
	parts := make([]string, len(mask))
	for i, part := range mask {
		parts[i] = strconv.FormatUint(uint64(part), 10)
	}
	return strings.Join(parts, ".")
}

func (iface *Interface) String() string {
	ifaceString := "%-10s  %-15s  %-15s  %-18s  %4d  %17s  %s"
	if os == win {
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
