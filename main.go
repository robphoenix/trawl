package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/rdegges/go-ipify"
)

const (
	VERSION = "v0.1.0"
)

// Interface provides the information for a device interface
type Interface struct {
	Name        string
	IPv4Address string
	IPv4Mask    string
	IPv4Network string
	IPv6Address string
}

var version bool

func init() {
	// parse flags
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.Parse()
}

func main() {

	if version {
		fmt.Println(VERSION)
		return
	}

	c := make(chan *Interface)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, iface := range ifaces {
		go func(iface net.Interface) {
			i, err := New(iface)
			if err != nil {
				log.Fatal(err)
			}
			c <- i
		}(iface)
	}

	// get public IP address
	pubIP, err := ipify.GetIp()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n")
	for range ifaces {
		iface := <-c
		fmt.Println(iface.String())
	}

	fmt.Printf("public      %s\n", string(pubIP))
}

// New instantiates an Interface object for the passed in net.Interface type
// representing a device interface
func New(iface net.Interface) (i *Interface, err error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return i, err
	}

	// get IPv4 address
	ipv4 := addrs[0].String()
	ipv4Split := strings.Split(ipv4, "/")
	ipv4Address := ipv4Split[0]

	// get IPv4 mask and convert to dotted decimal
	ipv4Cidr := ipv4Split[1]
	ipv4Mask, err := toDottedDec(ipv4Cidr)
	if err != nil {
		return i, err
	}

	// get IPv4 network
	_, ipnet, err := net.ParseCIDR(addrs[0].String())
	if err != nil {
		return i, err
	}

	// get IPv6 address & mask
	var ipv6Address string
	if len(addrs) > 1 {
		ipv6Address = addrs[1].String()
	}

	return &Interface{
		Name:        iface.Name,
		IPv4Address: ipv4Address,
		IPv4Mask:    ipv4Mask,
		IPv4Network: ipnet.String(),
		IPv6Address: ipv6Address,
	}, nil
}

func (iface *Interface) String() string {
	ifaceString := "%-10s  %-15s  %-15s  %-18s %s"
	return fmt.Sprintf(
		ifaceString,
		iface.Name,
		iface.IPv4Address,
		iface.IPv4Mask,
		iface.IPv4Network,
		iface.IPv6Address,
	)
}

func toDottedDec(cidr string) (s string, err error) {
	maskBits := []string{"", "128", "192", "224", "240", "248", "252", "254", "255"}
	n, err := strconv.Atoi(cidr)
	if err != nil {
		return s, err
	}

	if n > 32 || n < 0 {
		return s, fmt.Errorf("Not a valid network mask: %s", cidr)
	}

	allOnes := n / 8
	someOnes := n % 8
	mask := make([]string, 4)

	for i := 0; i < allOnes; i++ {
		mask[i] = "255"
	}

	if maskBits[someOnes] != "" {
		mask[allOnes] = maskBits[someOnes]
	}

	for i, octet := range mask {
		if octet == "" {
			mask[i] = "0"
		}
	}

	return strings.Join(mask, "."), nil
}
