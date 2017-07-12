# ![trawl](/logo.png)

> A strong fishing net for dragging along the sea bottom to collect IP addresses
> and similar flotsam & jetsam.

Prints out network interface information to the console, much like
`ifconfig`/`ipconfig`/`ip`/`ifdata`

[![Issue Count](https://codeclimate.com/github/robphoenix/trawl/badges/issue_count.svg?style=flat-square)](https://codeclimate.com/github/robphoenix/trawl)
[![Go Report Card](https://goreportcard.com/badge/github.com/robphoenix/trawl)](https://goreportcard.com/report/github.com/robphoenix/trawl)

```sh
$ trawl
wlp1s0   192.168.1.78  255.255.255.0    192.168.1.0/24  1500  10:02:b5:e4:de:8c  fe80::defe:3c33:4335:e669/64
docker0  172.17.0.1    255.255.0.0      172.17.0.0/16   1500  02:42:78:9a:c1:36  -
tun0     10.47.10.6    255.255.255.255  10.47.10.6/32   1500  -                  -
```

## Flags

Include column names

```sh
$ trawl -n
Name     IPv4 Address  IPv4 Mask        IPv4 Network    MTU   MAC Address        IPv6 Address
----     ------------  ---------        ------------    ---   -----------        ------------
wlp1s0   192.168.1.78  255.255.255.0    192.168.1.0/24  1500  10:02:b5:e4:de:8c  fe80::defe:3c33:4335:e669/64
docker0  172.17.0.1    255.255.0.0      172.17.0.0/16   1500  02:42:78:9a:c1:36  -
tun0     10.47.10.6    255.255.255.255  10.47.10.6/32   1500  -                  -
```

Filter interface names using a case insensitive regular expression

```sh
$ trawl -f w
wlp1s0  192.168.1.78  255.255.255.0  192.168.1.0/24  1500  10:02:b5:e4:de:8c  fe80::defe:3c33:4335:e669/64

$ trawl -n -f dock
Name     IPv4 Address  IPv4 Mask    IPv4 Network   MTU   MAC Address        IPv6 Address
----     ------------  ---------    ------------   ---   -----------        ------------
docker0  172.17.0.1    255.255.0.0  172.17.0.0/16  1500  02:42:78:9a:c1:36  -
```

Get a list of available interfaces

```sh
$ trawl -i
lo, wlp1s0, docker0, tun0
```

The loopback interface is ignored by default, but you can include it if you want

```sh
$ trawl -l
lo       127.0.0.1     255.0.0.0        127.0.0.0/8     65536  -                  ::1/128
wlp1s0   192.168.1.78  255.255.255.0    192.168.1.0/24  1500   10:02:b5:e4:de:8c  fe80::defe:3c33:4335:e669/64
docker0  172.17.0.1    255.255.0.0      172.17.0.0/16   1500   02:42:78:9a:c1:36  -
tun0     10.47.10.6    255.255.255.255  10.47.10.6/32   1500   -                  -
```

Specify the particular interface you want to know about

```sh
$ trawl wlp1s0
wlp1s0  192.168.1.78  255.255.255.0  192.168.1.0/24  1500  10:02:b5:e4:de:8c  fe80::defe:3c33:4335:e669/64
```

Get only the specific information you want, requires an interface name be provided

```sh
# IPv4 Address
$ trawl -a wlp1s0
192.168.1.78
# IPv4 Subnet Mask
$ trawl -m wlp1s0
255.255.255.0
# IPv4 Network
$ trawl -s wlp1s0
192.168.1.0/24
# IPv4 MTU
$ trawl -u wlp1s0
1500
# MAC Address
$ trawl -hw wlp1s0
10:02:b5:e4:de:8c
# IPv6 Address & Mask
$ trawl -6a wlp1s0
fe80::defe:3c33:4335:e669/64
```

You can also get your public IP address

```sh
$ trawl -p
104.238.169.73
```

All the same functionality is available in Windows

```cmd
% trawl
Local Area Connection 4              169.254.17.182   255.255.0.0      169.254.0.0/16      1500   02:00:4c:4f:4f:50        fe80::6cd7:885:5ae5:11b6/64
Wireless Network Connection          10.26.101.28     255.255.255.0    10.26.101.0/24      1500   24:77:03:c1:7e:2c        fe80::48e8:96c3:7457:8a3d/64
VirtualBox Host-Only Network         192.168.56.1     255.255.255.0    192.168.56.0/24     1500   0a:00:27:00:00:1a        fe80::31ac:de12:1d27:fbc9/64
VirtualBox Host-Only Network #2      10.0.0.1         255.255.0.0      10.0.0.0/16         1500   0a:00:27:00:00:1c        fe80::701e:c603:1aee:597e/64
Teredo Tunneling Pseudo-Interface    -                -                -                   1280   00:00:00:00:00:00:00:e0  fe80::1cea:232a:c110:463d/64
```

## Installation

If you don't have the Go programming language installed you can download the
appropriate binary for your system from the [releases page](https://github.com/robphoenix/trawl/releases),
rename it as `trawl`, and put it in your path ([howto ubuntu](https://askubuntu.com/questions/440691/add-a-binary-to-my-path)/[howto windows](https://uk.mathworks.com/matlabcentral/answers/94933-how-do-i-edit-my-system-path-in-windows?requestedDomain=www.mathworks.com)).

> I have not tried all of the binaries, so if there's a problem with one let me know, thanks.

If you do have Go installed...

```
go get -u github.com/robphoenix/trawl
```

## Acknowledgements

I totally used the awesome @jessfraz's [battery](https://github.com/jessfraz/battery)
as a starting point and continual touchstone for how to build this. Trawl is far
from perfect, but I've learnt from it, which was more the point.

Boat graphic by <a href="http://www.flaticon.com/authors/freepik">Freepik</a> from <a href="http://www.flaticon.com/">Flaticon</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0">CC BY 3.0</a>. Made with <a href="http://logomakr.com" title="Logo Maker">Logo Maker</a>
