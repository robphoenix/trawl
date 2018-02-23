# ![trawl](/logo.png)

> A strong fishing net for dragging along the sea bottom to collect IP addresses
> and similar flotsam & jetsam.

Prints out network interface information to the console, much like
`ifconfig`/`ipconfig`/`ip`/`ifdata`.

[![Issue Count](https://codeclimate.com/github/robphoenix/trawl/badges/issue_count.svg?style=flat-square)](https://codeclimate.com/github/robphoenix/trawl)
[![Go Report Card](https://goreportcard.com/badge/github.com/robphoenix/trawl)](https://goreportcard.com/report/github.com/robphoenix/trawl)

```console
❯ trawl
eth1   10.62.10.6     255.255.255.252  10.62.10.4/30   1500  00:ff:28:31:d0:72  fe80::d824:2e8d:bf80:69c9/64  2
wifi0  192.168.1.242  255.255.255.0    192.168.1.0/24  1500  10:02:b5:e4:de:8c  fe80::ed51:1db6:b32:ad90/64   2
```

## Flags

Show column names.

```console
❯ trawl -n                                                                                                                                                                    [1]
Name   IPv4 Address   IPv4 Mask        IPv4 Network    MTU   MAC Address        IPv6 Address                  Address Count
----   ------------   ---------        ------------    ---   -----------        ------------                  -------------
eth1   10.62.10.6     255.255.255.252  10.62.10.4/30   1500  00:ff:28:31:d0:72  fe80::d824:2e8d:bf80:69c9/64  2
wifi0  192.168.1.242  255.255.255.0    192.168.1.0/24  1500  10:02:b5:e4:de:8c  fe80::ed51:1db6:b32:ad90/64   2
```

Filter interface names using a case insensitive regular expression.

```console
❯ trawl -f wi
wifi0  192.168.1.242  255.255.255.0  192.168.1.0/24  1500  10:02:b5:e4:de:8c  fe80::ed51:1db6:b32:ad90/64  2

❯ trawl -n -f eth
Name  IPv4 Address  IPv4 Mask        IPv4 Network   MTU   MAC Address        IPv6 Address                  Address Count
----  ------------  ---------        ------------   ---   -----------        ------------                  -------------
eth1  10.62.10.6    255.255.255.252  10.62.10.4/30  1500  00:ff:28:31:d0:72  fe80::d824:2e8d:bf80:69c9/64  2
```

Get a list of available interfaces. Without any flags `trawl` only prints out interfaces which are _up_.

```console
❯ trawl -i
eth0, eth1, lo, wifi0, wifi1, eth2
```

The loopback interface is ignored by default, but you can include it if you like.

```console
❯ trawl -l
eth1   10.62.10.6     255.255.255.252  10.62.10.4/30   1500  00:ff:28:31:d0:72  fe80::d824:2e8d:bf80:69c9/64  2
lo     127.0.0.1      255.0.0.0        127.0.0.0/8     1500  -                  ::1/128                       2
wifi0  192.168.1.242  255.255.255.0    192.168.1.0/24  1500  10:02:b5:e4:de:8c  fe80::ed51:1db6:b32:ad90/64   2
```

Specify the particular interface you want to know about.

```console
❯ trawl wifi0
wifi0  192.168.1.242  255.255.255.0  192.168.1.0/24  1500  10:02:b5:e4:de:8c  fe80::ed51:1db6:b32:ad90/64  2
```

Show only the specific information you want, requires an interface name be provided.

```console
# IPv4 Address
❯ trawl -a wifi0
192.168.1.242
# IPv4 Subnet Mask
❯ trawl -m wifi0
255.255.255.0
# IPv4 Network
❯ trawl -s wifi0
192.168.1.0/24
# IPv4 MTU
❯ trawl -u wifi0
1500
# MAC Address
❯ trawl -hw wifi0
10:02:b5:e4:de:8c
# IPv6 Address & Mask
❯ trawl -6a wifi0
fe80::ed51:1db6:b32:ad90/64
```

Print a complete list of addresses for an interface.

```console
❯ trawl -4c wifi0
192.168.0.100/24
10.90.0.18/16

❯ trawl -6c wifi0
fe80::defe:3c33:4335:e669/64
fe80::/10
```

You can also get your public IP address.

```console
❯ trawl -p
104.238.169.73
```

All the same functionality is available in Windows.

```cmd
C:\Users\robphoenix>trawl -l -n
Name                         IPv4 Address   IPv4 Mask        IPv4 Network    MTU   MAC Address        IPv6 Address                  Address Count
----                         ------------   ---------        ------------    ---   -----------        ------------                  -------------
Ethernet                     10.62.10.6     255.255.255.252  10.62.10.4/30   1500  00:ff:28:31:d0:72  fe80::d824:2e8d:bf80:69c9/64  2
Wi-Fi                        192.168.1.242  255.255.255.0    192.168.1.0/24  1500  10:02:b5:e4:de:8c  fe80::ed51:1db6:b32:ad90/64   2
Loopback Pseudo-Interface 1  127.0.0.1      255.0.0.0        127.0.0.0/8     -1    -                  ::1/128                       2
```

## Installation

If you don't have the Go programming language installed you can download the
appropriate binary for your system from the [releases page](https://github.com/robphoenix/trawl/releases),
rename it as `trawl`, and put it in your path ([howto ubuntu](https://askubuntu.com/questions/440691/add-a-binary-to-my-path)/[howto windows](https://uk.mathworks.com/matlabcentral/answers/94933-how-do-i-edit-my-system-path-in-windows?requestedDomain=www.mathworks.com)).

> I have not tried all of the binaries, so if there's a problem with one let me know, thanks.

If you do have Go installed...

```console
go get -u github.com/robphoenix/trawl
```

## Acknowledgements

I totally used the awesome @jessfraz's [battery](https://github.com/jessfraz/battery)
as a starting point and continual touchstone for how to build this. Trawl is far
from perfect, but I've learnt from it, which was more the point.

Boat graphic by <a href="http://www.flaticon.com/authors/freepik">Freepik</a> from <a href="http://www.flaticon.com/">Flaticon</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0">CC BY 3.0</a>. Made with <a href="http://logomakr.com" title="Logo Maker">Logo Maker</a>
