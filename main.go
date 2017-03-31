package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"

	"github.com/rdegges/go-ipify"
)

const (
	// Version of current release
	Version           = "v0.1.4"
	os                = runtime.GOOS
	win               = "windows"
	linux             = "linux"
	underlineChar     = "-"
	nameHeader        = "Name"
	ipv4AddrHeader    = "IPv4 Address"
	ipv4MaskHeader    = "IPv4 Mask"
	ipv4NetworkHeader = "IPv4 Network"
	mtuHeader         = "MTU"
	macHeader         = "MAC Address"
	ipv6AddrHeader    = "IPv6 Address"
	windowsString     = "%-35s  %-15s  %-15s  %-18s  %-4s  %-17s  %s\n"
	linuxString       = "%-10s  %-15s  %-15s  %-18s  %-4s  %-17s  %s\n"
)

var (
	version bool
	public  bool
	names   bool
)

func init() {
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&public, "public", false, "print public IP address and exit")
	flag.BoolVar(&public, "p", false, "print public IP address and exit (shorthand)")
	flag.BoolVar(&names, "names", false, "print header names")
	flag.BoolVar(&names, "n", false, "print header names (shorthand)")
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

	if names {
		printHeaders()
	}

	ifaces := getIfaces()
	for _, iface := range ifaces {
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
	switch os {
	case win:
		s = windowsString
	case linux:
		s = linuxString
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

func getIfaces() (ifaces []net.Interface) {
	allIfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range allIfaces {
		up := iface.Flags & net.FlagUp
		loopback := iface.Flags & net.FlagLoopback
		if up != 0 && loopback == 0 {
			ifaces = append(ifaces, iface)
		}
	}
	return
}
