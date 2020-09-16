# tunnel
A simple tunnel (tcp + udp) writed with go.

# Installation
To install this package, you need to install Go and set your Go workspace first.

1. The first need Go installed (version 1.11+ is required), then you can use the below Go command to install tunnel.

```shell
$ go get -u github.com/jummyliu/tunnel@master
```

# Usage

```shell
$ tunnel -h

Usage: 
	tunnel [-h]
	tunnel [-t|-tcp <tunnelinfo>]... [-u|-udp <tunnelinfo>]... [-a|-all <tunnelinfo>]...
	tunnel 

Params:
	-t | -tcp <tunnelinfo>
		create a tcp tunnel with tunnelinfo
	-u | -udp <tunnelinfo>
		create a udp tunnel with tunnelinfo
	-a | -all <tunnelinfo>
		create both tcp and udp tunnel with tunnelinfo
	-h	help

	<tunnelinfo>
		sip:sport=dip:dport
		sip:sport=dport
		sport=dip:dport
		sport=dport
		port
		- default sip is <0.0.0.0>
		- default dip is <localhost>
```

```shell
$ tunnel -tcp 3306
procy (tcp) 0.0.0.0:3306 => localhost:3306

$ tunnel -udp 8080
procy (udp) 0.0.0.0:8080 => localhost:8080

$ tunnel -a 443
procy (udp) 0.0.0.0:443 => localhost:443
procy (tcp) 0.0.0.0:443 => localhost:443
```