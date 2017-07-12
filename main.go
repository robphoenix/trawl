package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
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
			v4Mask = toDottedDec(network.Mask)
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

func main() {

	if version {
		fmt.Println(Version)
		return
	}

	if public {
		p, err := ipify.GetIp()
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(p)
		return
	}

	if interfaces {
		fmt.Println(availableInterfaces())
		return
	}

	args := flag.Args()
	if len(args) > 0 {
		iface, err := net.InterfaceByName(args[0])
		if err != nil {
			log.Fatal(err)
			return
		}
		i, err := New(*iface)
		if err != nil {
			log.Fatal(err)
			return
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
			fmt.Fprintln(w, tabbedNames())
		}
		fmt.Fprintln(w, i)
		w.Flush()
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)

	if names {
		fmt.Fprintln(w, tabbedNames())
	}

	for _, iface := range validInterfaces(loopback, filter) {
		i, err := New(iface)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Fprintln(w, i)
	}
	w.Flush()
}

func availableInterfaces() string {
	var ifs []string
	all, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range all {
		ifs = append(ifs, iface.Name)
	}
	return strings.Join(ifs, ", ")
}

func validInterfaces(loopback bool, filter string) []net.Interface {
	var valid []net.Interface
	all, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range all {
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
			valid = append(valid, iface)
		}
	}
	return valid
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

func setMissingValue(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func tabbedNames() string {
	ns := []string{
		"Name",
		"IPv4 Address",
		"IPv4 Mask",
		"IPv4 Network",
		"MTU",
		"MAC Address",
		"IPv6 Address",
	}
	var underlined []string
	for _, s := range ns {
		underlined = append(underlined, strings.Repeat("-", len(s)))
	}
	var buf bytes.Buffer
	buf.WriteString(strings.Join(ns, "\t"))
	buf.WriteString("\n")
	buf.WriteString(strings.Join(underlined, "\t"))
	return buf.String()
}
