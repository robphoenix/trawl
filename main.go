package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/rdegges/go-ipify"
)

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
