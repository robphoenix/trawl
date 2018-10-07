package main

import (
	"errors"
	"flag"
	"fmt"
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
	Version = "v0.4.0"
)

var (
	version   bool
	public    bool
	names     bool
	loopback  bool
	ifaces    bool
	filter    string
	v4addr    bool
	v4mask    bool
	v4net     bool
	mtu       bool
	mac       bool
	v6addr    bool
	v4compl   bool
	v6compl   bool
	addrN     bool
	usageText = `
Trawl %s

A strong fishing net for dragging along the sea bottom
to collect IP addresses and similar flotsam & jetsam.

trawl [options...] <interface>

Options:
  -version, -v	     print version and exit

  -interfaces, -i    list available interfaces

  -loopback, -l      include loopback interface in output,
                     loopback is not included by default

  -names, -n         include output column names,
                     column names are not printed by default

  -public, -p        print only the public IP address

  -filter, -f	     filter, by name, which interfaces are displayed,
                     using a regular expression

  -a	             print only the IPv4 address, requires interface name
  -m	             print only the IPv4 subnet mask, requires interface name
  -s	             print only the IPv4 network (subnet), requires interface name
  -hw                print only the MAC address (hardware address),
		     requires interface name
  -u	             print only the MTU, requires interface name
  -6a                print only the IPv6 address, requires interface name
  -4c                print the complete list of IPv4 addresses an interface has,
		     includes subnet mask, requires interface name
  -6c                print the complete list of IPv6 addresses an interface has,
		     includes subnet mask, requires interface name
  -c                 print the number of IPv4 & IPv6 addresses associated with an interface,
                     requires interface name
`
)

func init() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf(usageText, Version))
	}

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&public, "public", false, "print public IP address and exit")
	flag.BoolVar(&public, "p", false, "print public IP address and exit (shorthand)")
	flag.BoolVar(&names, "names", false, "print header names")
	flag.BoolVar(&names, "n", false, "print header names (shorthand)")
	flag.BoolVar(&loopback, "loopback", false, "include loopback interface in output")
	flag.BoolVar(&loopback, "l", false, "include loopback interface in output (shorthand)")
	flag.BoolVar(&ifaces, "interfaces", false, "list available interfaces")
	flag.BoolVar(&ifaces, "i", false, "list available interfaces (shorthand)")
	flag.StringVar(&filter, "filter", "", "filter interface names with a regular expression")
	flag.StringVar(&filter, "f", "", "filter interface names with a regular expression (shorthand)")
	flag.BoolVar(&v4addr, "a", false, "print only IPv4 address, requires interface name")
	flag.BoolVar(&v4mask, "m", false, "print only IPv4 subnet mask, requires interface name")
	flag.BoolVar(&v4net, "s", false, "print only IPv4 network (subnet), requires interface name")
	flag.BoolVar(&mtu, "u", false, "print only MTU, requires interface name")
	flag.BoolVar(&mac, "hw", false, "print only MAC address (hardware address), requires interface name")
	flag.BoolVar(&v6addr, "6a", false, "print only IPv6 address, requires interface name")
	flag.BoolVar(&v4compl, "4c", false, "print all IPv4 address, requires interface name")
	flag.BoolVar(&v6compl, "6c", false, "print all IPv6 address, requires interface name")
	flag.BoolVar(&addrN, "c", false, "print count of all IPv4 & IPv6 addresses, requires interface name")
	flag.Parse()
}

// Iface provides the information for a device interface
type Iface struct {
	AddressCount string
	HardwareAddr string
	IPv4Addr     string
	IPv4Mask     string
	IPv4Network  string
	IPv6Addr     string
	MTU          string
	Name         string
}

func (i *Iface) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t",
		check(i.Name),
		check(i.IPv4Addr),
		check(i.IPv4Mask),
		check(i.IPv4Network),
		check(i.MTU),
		check(i.HardwareAddr),
		check(i.IPv6Addr),
		check(i.AddressCount),
	)
}

// New instantiates an Iface object representing a device interface
func New(netIface net.Interface) (*Iface, error) {

	ipv4s, ipv6s, n := expand(&netIface)

	// get IPv4 address and network
	var v4Addr, v4Mask, v4Net string
	if len(ipv4s) > 0 {
		addr, network, err := net.ParseCIDR(ipv4s[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if addr != nil {
			v4Addr = addr.String()
			v4Mask = dotted(network.Mask)
			v4Net = network.String()
		}
	}
	var ipv6 string
	if len(ipv6s) > 0 {
		ipv6 = ipv6s[0]
	}

	return &Iface{
		AddressCount: n,
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
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(p)
		return
	}

	if ifaces {
		fmt.Println(available())
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	args := flag.Args()
	if len(args) > 0 {
		iface, err := net.InterfaceByName(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		ipv4s, ipv6s, _ := expand(iface)

		i, err := New(*iface)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		if v4addr {
			fmt.Println(i.IPv4Addr)
			return
		}
		if v4mask {
			fmt.Println(i.IPv4Mask)
			return
		}
		if v4net {
			fmt.Println(i.IPv4Network)
			return
		}
		if mtu {
			fmt.Println(i.MTU)
			return
		}
		if mac {
			fmt.Println(i.HardwareAddr)
			return
		}
		if v6addr {
			fmt.Println(i.IPv6Addr)
			return
		}
		if v4compl {
			for _, v := range ipv4s {
				fmt.Println(v)
			}
			return
		}
		if v6compl {
			for _, v := range ipv6s {
				fmt.Println(v)
			}
			return
		}
		if addrN {
			fmt.Println(i.AddressCount)
			return
		}
		if names {
			fmt.Fprintln(w, header())
		}
		fmt.Fprintln(w, i)
		return
	}

	if v4addr || v4mask || v4net || mtu || mac || v6addr || v4compl || v6compl || addrN {
		fmt.Fprintln(os.Stderr, errors.New("requires interface name"))
		return
	}

	if names {
		fmt.Fprintln(w, header())
	}

	for _, iface := range usable(loopback, filter) {
		i, err := New(iface)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Fprintln(w, i)
	}
}

func available() string {
	var ifs []string
	all, err := net.Interfaces()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	for _, iface := range all {
		ifs = append(ifs, iface.Name)
	}
	return strings.Join(ifs, ", ")
}

func usable(loopback bool, filter string) []net.Interface {
	var valid []net.Interface
	all, err := net.Interfaces()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
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
			fmt.Fprintln(os.Stderr, err.Error())
		}
		// does the interface have any available addresses?
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		// is the interface up?
		up := iface.Flags & net.FlagUp
		if up != 0 && l == 0 && matched && len(addrs) > 0 {
			valid = append(valid, iface)
		}
	}
	return valid
}

func expand(iface *net.Interface) (ipv4s, ipv6s []string, num string) {
	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	for _, addr := range addrs {
		a := addr.String()
		switch {
		case strings.Contains(a, ":"):
			ipv6s = append(ipv6s, a)
		case strings.Contains(a, "."):
			ipv4s = append(ipv4s, a)
		}
	}
	return ipv4s, ipv6s, strconv.Itoa(len(addrs))
}

func dotted(mask net.IPMask) string {
	parts := make([]string, len(mask))
	for i, part := range mask {
		parts[i] = strconv.FormatUint(uint64(part), 10)
	}
	return strings.Join(parts, ".")
}

func check(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func header() string {
	ns := []string{
		"Name",
		"IPv4 Address",
		"IPv4 Mask",
		"IPv4 Network",
		"MTU",
		"MAC Address",
		"IPv6 Address",
		"Address Count",
	}
	var underlined []string
	for _, s := range ns {
		underlined = append(underlined, strings.Repeat("-", len(s)))
	}
	n := strings.Join(ns, "\t")
	u := strings.Join(underlined, "\t")
	return n + "\n" + u + "\t"
}
