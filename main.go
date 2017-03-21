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

type Interface struct {
	Name        string
	IPv4Address string
	IPv4Mask    string
	IPv4Network string
	IPv6Address string
}

type ExternalIP string

func main() {

	flag.Parse()
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n")

	for _, iface := range ifaces {
		fmt.Println(New(iface).String())
	}

	// public
	pubIP, err := ipify.GetIp()
	if err != nil {
		fmt.Println("Couldn't get my IP address:", err)
	}
	fmt.Printf("public      %s\n", string(pubIP))
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

func New(iface net.Interface) *Interface {
	addrs, err := iface.Addrs()
	if err != nil {
		log.Fatal(err)
	}

	ipv4 := addrs[0].String()

	ipv4Cidr := strings.Split(ipv4, "/")[1]
	ipv4Mask, err := toDottedDec(ipv4Cidr)
	if err != nil {
		log.Fatal(err)
	}
	ip, ipnet, err := net.ParseCIDR(addrs[0].String())
	if err != nil {
		log.Fatal(err)
	}
	var ipv6Address string
	if len(addrs) > 1 {
		ipv6Address = addrs[1].String()
	}

	return &Interface{
		Name:        iface.Name,
		IPv4Address: ip.String(),
		IPv4Mask:    ipv4Mask,
		IPv4Network: ipnet.String(),
		IPv6Address: ipv6Address,
	}
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
