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
	MTU          string
	Name         string
}

// New instantiates an Interface object representing a device interface
func New(iface net.Interface) (*Interface, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return &Interface{}, err
	}

	// we can't rely on the order of the addresses in the addrs array
	ipv4, ipv6 := extractAddrs(addrs)

	// get IPv4 address and network
	var ipv4Addr net.IP
	var ipv4Network *net.IPNet
	if len(ipv4) > 0 {
		ipv4Addr, ipv4Network, err = net.ParseCIDR(ipv4)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &Interface{
		HardwareAddr: iface.HardwareAddr.String(),
		IPv4Addr:     ipv4Addr.String(),
		IPv4Mask:     toDottedDec(ipv4Network.Mask),
		IPv4Network:  ipv4Network.String(),
		IPv6Addr:     ipv6,
		MTU:          strconv.Itoa(iface.MTU),
		Name:         iface.Name,
	}, nil
}

func setMissingValue(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func (iface *Interface) String() string {
	ifaceString := osString()
	return fmt.Sprintf(
		ifaceString,
		iface.Name,
		setMissingValue(iface.IPv4Addr),
		setMissingValue(iface.IPv4Mask),
		setMissingValue(iface.IPv4Network),
		setMissingValue(iface.MTU),
		setMissingValue(iface.HardwareAddr),
		setMissingValue(iface.IPv6Addr),
	)
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
