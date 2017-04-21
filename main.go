package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"runtime"
	"strings"

	"github.com/rdegges/go-ipify"
)

const (
	// Version of current release
	Version           = "v0.1.4"
	opSys             = runtime.GOOS
	win               = "windows"
	linux             = "linux"
	darwin            = "darwin"
	underlineChar     = "-"
	nameHeader        = "Name"
	ipv4AddrHeader    = "IPv4 Address"
	ipv4MaskHeader    = "IPv4 Mask"
	ipv4NetworkHeader = "IPv4 Network"
	mtuHeader         = "MTU"
	macHeader         = "MAC Address"
	ipv6AddrHeader    = "IPv6 Address"
	windowsString     = "%-35s  %-15s  %-15s  %-18s  %-5s  %-17s  %s\n"
	linuxString       = "%-10s  %-15s  %-15s  %-18s  %-5s  %-17s  %s\n"
	darwinString      = "%-10s  %-15s  %-15s  %-18s  %-5s  %-17s  %s\n"
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
	flag.BoolVar(&ipv4address, "a", false, "print only IPv4 address, requires interface")
	flag.BoolVar(&ipv4mask, "m", false, "print only IPv4 subnet mask, requires interface")
	flag.BoolVar(&ipv4network, "s", false, "print only IPv4 network (subnet), requires interface")
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
			if names {
				printHeaders()
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
			fmt.Printf(i.String())
		}
		return
	}

	if names {
		printHeaders()
	}

	for _, iface := range getIfaces(loopback, filter) {
		i, err := New(iface)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(i.String())
	}
}

func underline(s string) string {
	return strings.Repeat(underlineChar, len(s))
}

func osString() (s string) {
	switch opSys {
	case win:
		s = windowsString
	case linux:
		s = linuxString
	case darwin:
		s = darwinString
	}
	return
}

func printHeaders() {
	headersString := osString()
	fmt.Printf(
		headersString,
		nameHeader,
		ipv4AddrHeader,
		ipv4MaskHeader,
		ipv4NetworkHeader,
		mtuHeader,
		macHeader,
		ipv6AddrHeader,
	)
	fmt.Printf(
		headersString,
		underline(nameHeader),
		underline(ipv4AddrHeader),
		underline(ipv4MaskHeader),
		underline(ipv4NetworkHeader),
		underline(mtuHeader),
		underline(macHeader),
		underline(ipv6AddrHeader),
	)
}

func getIfaces(loopback bool, filter string) (ifaces []net.Interface) {
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
	return
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
