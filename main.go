package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Interface struct {
	Name        string
	IPv4Address string
	IPv4Mask    string
	IPv4Network string
	IPv6Address string
	MACAddress  string
	MTU         int
}

type ExternalIP string

func main() {

	flag.Parse()
	wantIfaces := flag.Args()

	ifaces, err := getIfaces(wantIfaces)
	if err != nil {
		log.Fatal(err)
	}

	for _, iface := range ifaces {
		i := New(iface)
		// TODO print only if not empty
		fmt.Println(i.String())
	}
	resp, err := http.Get("https://api.ipify.org/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("\nExternal IP:\t%s\n", string(body))
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

func checkIface(iface string, got []net.Interface) (ifg net.Interface, err error) {
	for _, ifg := range got {
		if iface == ifg.Name {
			return ifg, nil
		}
	}
	// TODO list available interfaces
	return ifg, fmt.Errorf("sorry, interface [%s] is not available", iface)
}

func getIfaces(want []string) ([]net.Interface, error) {

	gotIfaces, err := net.Interfaces()
	if err != nil {
		return []net.Interface{}, err
		// log.Fatal(err)
	}

	if len(want) == 0 {
		return gotIfaces, nil
	}

	var ifaces []net.Interface

	// check interfaces are available
	for _, wif := range want {
		iface, err := checkIface(wif, gotIfaces)
		// TODO allow valid interfaces even when given with invalid interfaces
		if err != nil {
			return []net.Interface{}, err
			// log.Fatal(err)
		}
		ifaces = append(ifaces, iface)
	}

	return ifaces, nil
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
		MTU:         iface.MTU,
		MACAddress:  iface.HardwareAddr.String(),
	}
}

func (iface *Interface) String() string {
	ifaceString := `Interface %s
IPv4 address:	%s
IPv4 mask:	%s
IPv4 network:	%s
IPv6 address:	%s
MTU:		%d
MAC Address:	%s
`
	return fmt.Sprintf(
		ifaceString,
		iface.Name,
		iface.IPv4Address,
		iface.IPv4Mask,
		iface.IPv4Network,
		iface.IPv6Address,
		iface.MTU,
		iface.MACAddress,
	)
}
