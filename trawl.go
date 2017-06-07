package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// Iface provides the information for a device interface
type Iface struct {
	HardwareAddr string
	IPv4Addr     string
	IPv4Mask     string
	IPv4Network  string
	IPv6Addr     string
	MTU          string
	Name         string
}

type mask net.IPMask

// New instantiates an Iface object representing a device interface
func New(netIface net.Interface) (*Iface, error) {
	addrs, err := netIface.Addrs()
	if err != nil {
		return &Iface{}, err
	}

	// we can't rely on the order of the addresses in the addrs array
	ipv4, ipv6 := extractAddrs(addrs)

	// get IPv4 address and network
	var v4Addr, v4Mask, v4Net string
	if len(ipv4) > 0 {
		addr, network, err := net.ParseCIDR(ipv4)
		if err != nil {
			log.Fatal(err)
		}
		if addr != nil {
			v4Addr = addr.String()
			v4Mask = mask(network.Mask).toDottedDec()
			v4Net = network.String()
		}
	}

	return &Iface{
		HardwareAddr: netIface.HardwareAddr.String(),
		IPv4Addr:     v4Addr,
		IPv4Mask:     v4Mask,
		IPv4Network:  v4Net,
		IPv6Addr:     ipv6,
		MTU:          strconv.Itoa(netIface.MTU),
		Name:         netIface.Name,
	}, nil
}

func setMissingValue(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func (iface *Iface) String() string {
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

func (m mask) toDottedDec() string {
	parts := make([]string, len(m))
	for i, part := range m {
		parts[i] = strconv.FormatUint(uint64(part), 10)
	}
	return strings.Join(parts, ".")
}
