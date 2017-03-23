# ![trawl](/logo.png)

> A strong fishing net for dragging along the sea bottom to collect IP addresses
> and similar flotsam & jetsam.

Prints out network interface information to the console, much like
ifconfig/ipconfig/ip/ifdata, of the form:

```
Name  IPv4 Address  IPv4 Mask  IPv4 Network  MTU  MAC Address  IPv6 Address
```

Or it can print out your external public IP address.

## Linux

```sh
$ trawl
tun0        10.6.10.6        255.255.255.255  10.6.10.6/32        1500
wlp1s0      192.168.1.78     255.255.255.0    192.168.1.0/24      1500  10:02:b5:e4:de:8c  fe80::defe:3c33:4335:e669/64
docker0     172.17.0.1       255.255.0.0      172.17.0.0/16       1500  02:42:58:31:a9:87
```

```sh
$ trawl -p
104.238.169.73
```

## Windows

```sh
% trawl
VirtualBox Host-Only Network         192.168.56.1     255.255.255.0    192.168.56.0/24    fe80::31ac:de12:1d27:fbc9/64
VirtualBox Host-Only Network #2      10.0.0.1         255.255.0.0      10.0.0.0/16        fe80::701e:c603:1aee:597e/64
Local Area Connection 3              10.48.10.10      255.255.255.252  10.48.10.8/30      fe80::989:e670:8216:5528/64
Local Area Connection 4              169.254.17.182   255.255.0.0      169.254.0.0/16     fe80::6cd7:885:5ae5:11b6/64
Wireless Network Connection          192.168.1.224    255.255.255.0    192.168.1.0/24     fe80::48e8:96c3:7457:8a3d/64
Public                               46.19.140.62
```

## Installation

If you don't have the Go programming language installed you can download the binary from the
[releases page](https://github.com/robphoenix/trawl/releases) and put it in your
path ([howto ubuntu](https://askubuntu.com/questions/440691/add-a-binary-to-my-path)/[howto windows](https://uk.mathworks.com/matlabcentral/answers/94933-how-do-i-edit-my-system-path-in-windows?requestedDomain=www.mathworks.com)).

If you do have Go installed, get it like you do.

```
go get -u github.com/robphoenix/trawl
```

Or if you want to install from source you can clone with git:

```
git clone https://github.com/robphoenix/trawl.git
cd trawl
go install
```

## Acknowledgements

I totally used the awesome @jessfraz's [battery](https://github.com/jessfraz/battery)
as a starting point and continual touchstone for how to build this. Trawl is far
from perfect, but I've learnt from it, which was more the point.

Boat graphic by <a href="http://www.flaticon.com/authors/freepik">Freepik</a> from <a href="http://www.flaticon.com/">Flaticon</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0">CC BY 3.0</a>. Made with <a href="http://logomakr.com" title="Logo Maker">Logo Maker</a>
