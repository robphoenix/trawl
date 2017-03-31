package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/rdegges/go-ipify"
)

// current release version
const (
	OS           = runtime.GOOS
	Version      = "v0.1.3"
	linuxHeaders = `Name        IPv4 Address     IPv4 Mask        IPv4 Network        MTU   MAC Address        IPv6 Address` + "\n" +
		`----        ------------     ----------       ------------        ---   -----------        ------------`
	windowsHeaders = `Name                                 IPv4 Address     IPv4 Mask        IPv4 Network        MTU   MAC Address        IPv6 Address` + "\n" +
		`----                                 ------------     ----------       ------------        ---   -----------        ------------`
)

var (
	version bool
	public  bool
	names   bool
)

func init() {
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&public, "public", false, "print public IP address and exit")
	flag.BoolVar(&public, "p", false, "print public IP address and exit (shorthand)")
	flag.BoolVar(&names, "names", false, "print header names")
	flag.BoolVar(&names, "n", false, "print header names (shorthand)")
	flag.Parse()
}

func main() {
	if version {
		fmt.Println(Version)
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

	if names {
		switch OS {
		case "windows":
			fmt.Println(windowsHeaders)
		case "linux":
			fmt.Println(linuxHeaders)
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
