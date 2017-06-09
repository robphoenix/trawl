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
	outputString      = "%-*s  %-*s  %-*s  %-*s  %-*s  %-*s  %s\n"
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

type fieldLengths struct {
	name int
	addr int
	mask int
	net  int
	mtu  int
	mac  int
}

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
			fl := fieldLengths{
				len(i.Name),
				len(i.IPv4Addr),
				len(i.IPv4Mask),
				len(i.IPv4Network),
				len(i.MTU),
				len(i.HardwareAddr),
			}
			if names {
				fmt.Printf(headerString(fl))
			}
			fmt.Printf(ifaceString(fl, i))
		}
		return
	}

	var fl fieldLengths
	var ifs []*Iface

	for _, iface := range getIfaces(loopback, filter) {
		i, err := New(iface)
		if err != nil {
			log.Fatal(err)
		}
		ifs = append(ifs, i)
		fl.name = maxLen(fl.name, len(i.Name))
		fl.addr = maxLen(fl.addr, len(i.IPv4Addr))
		fl.mask = maxLen(fl.mask, len(i.IPv4Mask))
		fl.net = maxLen(fl.net, len(i.IPv4Network))
		fl.mtu = maxLen(fl.mtu, len(i.MTU))
		fl.mac = maxLen(fl.mac, len(i.HardwareAddr))
	}

	if names {
		fmt.Printf(headerString(fl))
	}

	for _, i := range ifs {
		fmt.Printf(ifaceString(fl, i))
	}
}

func underline(s string) string {
	return strings.Repeat(underlineChar, len(s))
}

func maxLen(m, l int) int {
	if l > m {
		return l
	}
	return m
}

func headerString(fl fieldLengths) string {
	var s string
	s += fmt.Sprintf(
		outputString,
		fl.name,
		nameHeader,
		fl.addr,
		ipv4AddrHeader,
		fl.mask,
		ipv4MaskHeader,
		fl.net,
		ipv4NetworkHeader,
		fl.mtu,
		mtuHeader,
		fl.mac,
		macHeader,
		ipv6AddrHeader,
	)
	s += fmt.Sprintf(
		outputString,
		fl.name,
		underline(nameHeader),
		fl.addr,
		underline(ipv4AddrHeader),
		fl.mask,
		underline(ipv4MaskHeader),
		fl.net,
		underline(ipv4NetworkHeader),
		fl.mtu,
		underline(mtuHeader),
		fl.mac,
		underline(macHeader),
		underline(ipv6AddrHeader),
	)
	return s
}

func ifaceString(fl fieldLengths, i *Iface) string {
	return fmt.Sprintf(
		outputString,
		fl.name,
		setMissingValue(i.Name),
		fl.addr,
		setMissingValue(i.IPv4Addr),
		fl.mask,
		setMissingValue(i.IPv4Mask),
		fl.net,
		setMissingValue(i.IPv4Network),
		fl.mtu,
		setMissingValue(i.MTU),
		fl.mac,
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
