package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", iface.Name)

		ipv4 := addrs[0].String()

		ipv4Mask := strings.Split(ipv4, "/")[1]
		ipv4Mask, err = toDottedDec(ipv4Mask)
		if err != nil {
			log.Fatal(err)
		}
		ip, ipnet, err := net.ParseCIDR(addrs[0].String())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\tIPv4 address:\t%s\n", ip)
		fmt.Printf("\tIPv4 network:\t%s\n", ipnet)
		fmt.Printf("\tIPv4 mask:\t%s\n", ipv4Mask)
		if len(addrs) > 1 {
			fmt.Printf("\tIPv6 address:\t%s\n", addrs[1])
		}
		fmt.Printf("\tMTU:\t\t%d\n", iface.MTU)
		if string(iface.HardwareAddr) != "" {
			fmt.Printf("\tMAC Address:\t%s\n", iface.HardwareAddr)
		}
	}
	resp, err := http.Get("https://api.ipify.org/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("\nExternal IP:\t%s\n", string(body))

}

func toDottedDec(cidr string) (string, error) {
	integer, err := strconv.Atoi(cidr)
	if err != nil {
		return "", err
	}
	if integer > 32 || integer < 0 {
		return "", fmt.Errorf("Not a valid network mask: %s", cidr)
	}

	maskBits := []string{"", "128", "192", "224", "240", "248", "252", "254", "255"}
	allOnes := integer / 8
	someOnes := integer % 8
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

	dottedDec := strings.Join(mask, ".")
	return dottedDec, nil
}
