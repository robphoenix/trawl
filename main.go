package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/rdegges/go-ipify"
)

// current release version
const (
	VERSION = "v0.1.1"
)

var version bool

func init() {
	flag.BoolVar(&version, "v", false, "print version and exit")
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

	var validIfaces []net.Interface

	for _, iface := range ifaces {
		if strings.Contains(iface.Flags.String(), "up") && !strings.Contains(iface.Flags.String(), "lo") {
			validIfaces = append(validIfaces, iface)
		}
	}

	for _, iface := range validIfaces {
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
	for range validIfaces {
		iface := <-c
		fmt.Println(iface.String())
	}

	fmt.Printf("public      %s\n", string(pubIP))
}
