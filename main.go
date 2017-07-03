package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/rdegges/go-ipify"
)

const (
	// Version of current release
	Version = "v0.2.1"
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
		iface, err := net.InterfaceByName(args[0])
		if err != nil {
			log.Fatal(err)
		}
		i, err := NewIface(*iface)
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
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
		if names {
			fmt.Fprintln(w, newHeaders())
		}
		fmt.Fprintln(w, i)
		w.Flush()
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)

	if names {
		fmt.Fprintln(w, newHeaders())
	}

	for _, iface := range getIfaces(loopback, filter) {
		i, err := NewIface(iface)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(w, i)
	}

	w.Flush()
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
