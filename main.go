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
	VERSION = "v0.1.2"
)

var (
	version bool
	public  bool
)

func init() {
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&public, "public", false, "print public IP address and exit")
	flag.BoolVar(&public, "p", false, "print public IP address and exit (shorthand)")
	flag.Parse()
}

func main() {

	if version {
		fmt.Println(VERSION)
		return
	}

	if public {
		pubIP, err := ipify.GetIp()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(pubIP))
		return
	}

	c := make(chan string)
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
			c <- i.String()
		}(iface)
	}

	for range validIfaces {
		iface := <-c
		fmt.Println(iface)
	}
}
