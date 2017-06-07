package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/rdegges/go-ipify"
)

const (
	// Version of current release
	Version           = "v0.2.1"
	underlineChar     = "-"
	nameHeader        = "Name"
	ipv4AddrHeader    = "IPv4 Address"
	ipv4MaskHeader    = "IPv4 Mask"
	ipv4NetworkHeader = "IPv4 Network"
	mtuHeader         = "MTU"
	macHeader         = "MAC Address"
	ipv6AddrHeader    = "IPv6 Address"
	outputString      = "%-*s  %-15s  %-15s  %-18s  %-5s  %-17s  %s\n"
)

var (
	version     bool
	public      bool
	names       bool
	loopback    bool
	interfaces  bool
	filter      string
	ipv4address bool
	ipv4mask    bool
	ipv4network bool
	mtu         bool
	hwAddress   bool
	ipv6address bool
)

func init() {
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&public, "public", false, "print public IP address and exit")
	flag.BoolVar(&public, "p", false, "print public IP address and exit (shorthand)")
	flag.BoolVar(&names, "names", false, "print header names")
	flag.BoolVar(&names, "n", false, "print header names (shorthand)")
	flag.BoolVar(&loopback, "loopback", false, "include loopback interface in output")
	flag.BoolVar(&loopback, "l", false, "include loopback interface in output (shorthand)")
	flag.BoolVar(&interfaces, "interfaces", false, "list available interfaces")
	flag.BoolVar(&interfaces, "i", false, "list available interfaces (shorthand)")
	flag.StringVar(&filter, "filter", "", "filter interface names with a regular expression")
	flag.StringVar(&filter, "f", "", "filter interface names with a regular expression (shorthand)")
	flag.BoolVar(&ipv4address, "a", false, "print only IPv4 address, requires interface name")
	flag.BoolVar(&ipv4mask, "m", false, "print only IPv4 subnet mask, requires interface name")
	flag.BoolVar(&ipv4network, "s", false, "print only IPv4 network (subnet), requires interface name")
	flag.BoolVar(&mtu, "u", false, "print only MTU, requires interface name")
	flag.BoolVar(&hwAddress, "hw", false, "print only MAC address (hardware address), requires interface name")
	flag.BoolVar(&ipv6address, "6a", false, "print only IPv6 address, requires interface name")
	flag.Parse()
}

func main() {

	if version {
		fmt.Println(Version)
		return
	}

	if public {
		pubIP, err := ipify.GetIp()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(pubIP)
		return
	}

	if interfaces {
		fmt.Println(availableIfaces())
		return
	}

	args := flag.Args()
	if len(args) > 0 {
		for _, arg := range args {
			iface, err := net.InterfaceByName(arg)
			if err != nil {
				log.Fatal(err)
			}
			i, err := New(*iface)
			if err != nil {
				log.Fatal(err)
			}
			if ipv4address {
				fmt.Printf("%s\n", i.IPv4Addr)
				return
			}
			if ipv4mask {
				fmt.Printf("%s\n", i.IPv4Mask)
				return
			}
			if ipv4network {
				fmt.Printf("%s\n", i.IPv4Network)
				return
			}
			if mtu {
				fmt.Printf("%s\n", i.MTU)
				return
			}
			if hwAddress {
				fmt.Printf("%s\n", i.HardwareAddr)
				return
			}
			if ipv6address {
				fmt.Printf("%s\n", i.IPv6Addr)
				return
			}
			maxLen := len(i.Name)
			if names {
				fmt.Printf(headerString(maxLen))
			}
			fmt.Printf(ifaceString(maxLen, i))
		}
		return
	}

	var maxLen int
	var ifs []*Iface

	for _, iface := range getIfaces(loopback, filter) {
		i, err := New(iface)
		if err != nil {
			log.Fatal(err)
		}
		ifs = append(ifs, i)
		if nameLen := len(i.Name); nameLen > maxLen {
			maxLen = nameLen
		}
	}

	if names {
		fmt.Printf(headerString(maxLen))
	}

	for _, i := range ifs {
		fmt.Printf(ifaceString(maxLen, i))
	}
}

func underline(s string) string {
	return strings.Repeat(underlineChar, len(s))
}

func headerString(l int) string {
	var s string
	s += fmt.Sprintf(
		outputString,
		l,
		nameHeader,
		ipv4AddrHeader,
		ipv4MaskHeader,
		ipv4NetworkHeader,
		mtuHeader,
		macHeader,
		ipv6AddrHeader,
	)
	s += fmt.Sprintf(
		outputString,
		l,
		underline(nameHeader),
		underline(ipv4AddrHeader),
		underline(ipv4MaskHeader),
		underline(ipv4NetworkHeader),
		underline(mtuHeader),
		underline(macHeader),
		underline(ipv6AddrHeader),
	)
	return s
}

func ifaceString(l int, i *Iface) string {
	return fmt.Sprintf(
		outputString,
		l,
		setMissingValue(i.Name),
		setMissingValue(i.IPv4Addr),
		setMissingValue(i.IPv4Mask),
		setMissingValue(i.IPv4Network),
		setMissingValue(i.MTU),
		setMissingValue(i.HardwareAddr),
		setMissingValue(i.IPv6Addr),
	)
}

func getIfaces(loopback bool, filter string) []net.Interface {
	var ifaces []net.Interface
	allIfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range allIfaces {
		var l int
		// is it a loopback interface? do we want the loopback interface?
		if !loopback {
			l = int(iface.Flags & net.FlagLoopback)
		}
		// does the interface pass the filter?
		matched, err := regexp.MatchString(filter, iface.Name)
		if err != nil {
			log.Fatal(err)
		}
		// does the interface have any available addresses?
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		// is the interface up?
		up := iface.Flags & net.FlagUp
		if up != 0 && l == 0 && matched && len(addrs) > 0 {
			ifaces = append(ifaces, iface)
		}
	}
	return ifaces
}

func availableIfaces() string {
	var availIfaces []string
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range ifaces {
		availIfaces = append(availIfaces, iface.Name)
	}
	return strings.Join(availIfaces, ", ")
}

func setMissingValue(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
