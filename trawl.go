package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type mask net.IPMask

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

func (i *Iface) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s",
		setMissingValue(i.Name),
		setMissingValue(i.IPv4Addr),
		setMissingValue(i.IPv4Mask),
		setMissingValue(i.IPv4Network),
		setMissingValue(i.MTU),
		setMissingValue(i.HardwareAddr),
		setMissingValue(i.IPv6Addr),
	)
}

// NewIface instantiates an Iface object representing a device interface
func NewIface(netIface net.Interface) (*Iface, error) {
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

type headers struct {
	ipv4Addr      string
	ipv4Mask      string
	ipv4Network   string
	ipv6Addr      string
	mtu           string
	mac           string
	name          string
	underlineChar string
}

func (h *headers) String() string {
	s := "%s\t%s\t%s\t%s\t%s\t%s\t%s"
	return fmt.Sprintf(s+"\n"+s,
		h.name,
		h.ipv4Addr,
		h.ipv4Mask,
		h.ipv4Network,
		h.mtu,
		h.mac,
		h.ipv6Addr,
		underline(h.name, h.underlineChar),
		underline(h.ipv4Addr, h.underlineChar),
		underline(h.ipv4Mask, h.underlineChar),
		underline(h.ipv4Network, h.underlineChar),
		underline(h.mtu, h.underlineChar),
		underline(h.mac, h.underlineChar),
		underline(h.ipv6Addr, h.underlineChar),
	)
}

func newHeaders() *headers {
	return &headers{
		ipv4Addr:      "IPv4 Address",
		ipv6Addr:      "IPv6 Address",
		ipv4Mask:      "IPv4 Mask",
		ipv4Network:   "IPv4 Network",
		mac:           "MAC Address",
		mtu:           "MTU",
		name:          "Name",
		underlineChar: "-",
	}
}

func setMissingValue(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func underline(s, c string) string {
	return strings.Repeat(c, len(s))
}
