# ![trawl](/logo.png)

> A strong fishing net for dragging along the sea bottom to collect IP addresses
> and similar flotsam & jetsam.

Prints out network interface information to the console, much like
`ifconfig`/`ipconfig`/`ip`/`ifdata`

```sh
$ trawl -h
Usage of trawl:
  -6a
    	print only IPv6 address, requires interface name
  -a	print only IPv4 address, requires interface name
  -f string
    	filter interface names with a regular expression (shorthand)
  -filter string
    	filter interface names with a regular expression
  -hw
    	print only MAC address (hardware address), requires interface name
  -i	list available interfaces (shorthand)
  -interfaces
    	list available interfaces
  -l	include loopback interface in output (shorthand)
  -loopback
    	include loopback interface in output
  -m	print only IPv4 subnet mask, requires interface name
  -n	print header names (shorthand)
  -names
    	print header names
  -p	print public IP address and exit (shorthand)
  -public
    	print public IP address and exit
  -s	print only IPv4 network (subnet), requires interface name
  -u	print only MTU, requires interface name
  -v	print version and exit (shorthand)
  -version
    	print version and exit
```

[![Issue Count](https://codeclimate.com/github/robphoenix/trawl/badges/issue_count.svg?style=flat-square)](https://codeclimate.com/github/robphoenix/trawl)

## Example Output

### Linux/macOS

```sh
$ trawl
wlp1s0      192.168.1.78     255.255.255.0    192.168.1.0/24      1500   10:02:b5:e4:de:8c  fe80::defe:3c33:4335:e669/64
docker0     172.17.0.1       255.255.0.0      172.17.0.0/16       1500   02:42:47:11:03:c2  -
tun0        10.8.8.25        255.255.255.0    10.8.8.0/24         1500   -                  -
```

The output ignores the loopback interface by default, but using flags you can include it if
you want, along with column names:

```sh
$ trawl -l -n
Name        IPv4 Address     IPv4 Mask        IPv4 Network        MTU    MAC Address        IPv6 Address
----        ------------     ---------        ------------        ---    -----------        ------------
lo          127.0.0.1        255.0.0.0        127.0.0.0/8         65536  -                  ::1/128
wlp1s0      192.168.1.78     255.255.255.0    192.168.1.0/24      1500   10:02:b5:e4:de:8c  fe80::defe:3c33:4335:e669/64
docker0     172.17.0.1       255.255.0.0      172.17.0.0/16       1500   02:42:47:11:03:c2  -
tun0        10.8.8.25        255.255.255.0    10.8.8.0/24         1500   -                  -
```

Trawl has a number of flags if you want to be more specific:

You can filter the interface names with a regular expression.

```sh
$ trawl -f w
wlp1s0      192.168.1.78     255.255.255.0    192.168.1.0/24      1500   10:02:b5:e4:de:8c  fe80::defe:3c33:4335:e669/64
$ trawl -f dock
docker0     172.17.0.1       255.255.0.0      172.17.0.0/16       1500   02:42:47:11:03:c2  -
```

Get a list of available interfaces.

```sh
$ trawl -i
lo, wlp1s0, docker0, tun0
```

Specify the particular interface you want to know about.

```sh
$ trawl wlp1s0
wlp1s0      192.168.1.78     255.255.255.0    192.168.1.0/24      1500   10:02:b5:e4:de:8c  fe80::defe:3c33:4335:e669/64
```

Get only the specific information you want for an interface. This does require
an interface name be provided.

```sh
$ trawl -a wlp1s0
192.168.1.78
$ trawl -m wlp1s0
255.255.255.0
$ trawl -s wlp1s0
192.168.1.0/24
$ trawl -u wlp1s0
1500
$ trawl -hw wlp1s0
10:02:b5:e4:de:8c
$ trawl -6a wlp1s0
fe80::defe:3c33:4335:e669/64
```

And you can also get your public IP address

```sh
$ trawl -p
104.238.169.73
```

### Windows

All the same functionality is available in Windows.

```sh
% trawl
Local Area Connection 4              169.254.17.182   255.255.0.0      169.254.0.0/16      1500  02:00:3d:5c:5c:50  fe80::6cd7:885:5ae5:11b6/64
Teredo Tunneling Pseudo-Interface                                                             0                     fe80::101e:24fb:c110:462c/64
VirtualBox Host-Only Network         192.168.56.1     255.255.255.0    192.168.56.0/24     1500  0a:00:32:00:00:2b  fe80::31ac:de12:1d27:fbc9/64
VirtualBox Host-Only Network #2      10.0.0.1         255.255.0.0      10.0.0.0/16         1500  0a:00:32:00:00:2b  fe80::701e:c603:1aee:597e/64
Local Area Connection                10.90.128.3      255.255.0.0      10.90.0.0/16        1500  d5:be:c4:70:34:f5  fe80::a4f5:c0bf:b0ca:5551/64
Wireless Network Connection          10.26.101.64     255.255.255.0    10.26.101.0/24      1500  87:77:a3:d1:7e:2c  fe80::48e8:96c3:7457:8a3d/64
```

```sh
% trawl -l -n
Name                                 IPv4 Address     IPv4 Mask        IPv4 Network        MTU   MAC Address        IPv6 Address
----                                 ------------     ----------       ------------        ---   -----------        ------------
Local Area Connection 4              169.254.17.182   255.255.0.0      169.254.0.0/16      1500  02:00:3d:5c:5c:50  fe80::6cd7:885:5ae5:11b6/64
Teredo Tunneling Pseudo-Interface                                                             0                     fe80::101e:24fb:c110:462c/64
VirtualBox Host-Only Network         192.168.56.1     255.255.255.0    192.168.56.0/24     1500  0a:00:32:00:00:2b  fe80::31ac:de12:1d27:fbc9/64
VirtualBox Host-Only Network #2      10.0.0.1         255.255.0.0      10.0.0.0/16         1500  0a:00:32:00:00:2b  fe80::701e:c603:1aee:597e/64
Local Area Connection                10.90.128.3      255.255.0.0      10.90.0.0/16        1500  d5:be:c4:70:34:f5  fe80::a4f5:c0bf:b0ca:5551/64
Wireless Network Connection          10.26.101.64     255.255.255.0    10.26.101.0/24      1500  87:77:a3:d1:7e:2c  fe80::48e8:96c3:7457:8a3d/64
Loopback Pseudo-Interface 1          127.0.0.1        255.0.0.0        127.0.0.0/8         -1                        ::1/128
```

```
% trawl -p
62.239.185.211
```

## Installation

If you don't have the Go programming language installed you can download the
appropriate binary for your system from the [releases page](https://github.com/robphoenix/trawl/releases),
rename it as `trawl`, and put it in your path ([howto ubuntu](https://askubuntu.com/questions/440691/add-a-binary-to-my-path)/[howto windows](https://uk.mathworks.com/matlabcentral/answers/94933-how-do-i-edit-my-system-path-in-windows?requestedDomain=www.mathworks.com)).

> I have not tried all of the binaries, so if there's a problem with one let
> me know, thanks.

If you do have Go installed...

```
go get -u github.com/robphoenix/trawl
```

## Acknowledgements

I totally used the awesome @jessfraz's [battery](https://github.com/jessfraz/battery)
as a starting point and continual touchstone for how to build this. Trawl is far
from perfect, but I've learnt from it, which was more the point.

Boat graphic by <a href="http://www.flaticon.com/authors/freepik">Freepik</a> from <a href="http://www.flaticon.com/">Flaticon</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0">CC BY 3.0</a>. Made with <a href="http://logomakr.com" title="Logo Maker">Logo Maker</a>
