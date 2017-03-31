package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/rdegges/go-ipify"
)

// current release version
const (
	VERSION = "v0.1.3"
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
		fmt.Println(pubIP)
		return
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	var validIfaces []net.Interface

	for _, iface := range ifaces {
		up := iface.Flags & net.FlagUp
		loopback := iface.Flags & net.FlagLoopback
		if up != 0 && loopback == 0 {
			validIfaces = append(validIfaces, iface)
		}
	}

	for _, iface := range validIfaces {
		i, err := New(iface)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(i)
	}
}
