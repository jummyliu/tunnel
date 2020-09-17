package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/jummyliu/tunnel/flagtype"
)

var (
	tcpTunnel flagtype.FlagArr
	udpTunnel flagtype.FlagArr
	allTunnel flagtype.FlagArr
	h         bool

	tunnels []*Tunnel
)

var usageStr = `
Easy to create (tcp | udp) tunnel.
Usage: 
	tunnel [-h]
	tunnel [-t|-tcp <tunnelinfo>]... [-u|-udp <tunnelinfo>]... [-a|-all <tunnelinfo>]...

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
`

func usage() {
	fmt.Println(usageStr)
}

func init() {
	flag.Var(&tcpTunnel, "tcp", "tcp tunnel")
	flag.Var(&tcpTunnel, "t", "tcp tunnel")
	flag.Var(&udpTunnel, "udp", "udp tunnel")
	flag.Var(&udpTunnel, "u", "udp tunnel")
	flag.Var(&allTunnel, "all", "tcp & udp tunnel")
	flag.Var(&allTunnel, "a", "tcp & udp tunnel")
	flag.BoolVar(&h, "h", false, "help")

	flag.Parse()
	flag.Usage = usage
	if h {
		flag.Usage()
		os.Exit(0)
	}

	// tcp
	tunnels = append(tunnels, parseTunnel([]string(tcpTunnel), TunnelTypeTCP)...)
	// udp
	tunnels = append(tunnels, parseTunnel([]string(udpTunnel), TunnelTypeUDP)...)
	// all
	tunnels = append(tunnels, parseTunnel([]string(allTunnel), TunnelTypeALL)...)

	if len(tunnels) == 0 {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	var wg sync.WaitGroup
	for _, tunnel := range tunnels {
		wg.Add(1)
		go func(tunnel *Tunnel) {
			tunnel.Build()
			wg.Wait()
		}(tunnel)
	}
	wg.Wait()
}

const (
	defaultSip = "0.0.0.0"
	defaultDip = "localhost"
)

// error
var (
	ErrPortParseFailure = errors.New("port parse failure")
	ErrIPParseFailure   = errors.New("ip parse failure")
)

func parseTunnel(arr []string, t TunnelType) []*Tunnel {
	tunnels := make([]*Tunnel, 0)
	for _, tunnelStr := range arr {
		tunnelSplitArr := strings.Split(tunnelStr, "=")

		switch len(tunnelSplitArr) {
		case 1:
			// port
			_, port, err := parseIP(tunnelSplitArr[0], defaultSip)
			if err != nil {
				fmt.Println("Error on", tunnelStr, ":", err)
				flag.Usage()
				os.Exit(-1)
			}
			tunnels = append(tunnels, &Tunnel{
				Sip:   defaultSip,
				Sport: port,
				Dip:   defaultDip,
				Dport: port,
				Type:  t,
			})
		case 2:
			// sip:sport=dip:dport
			// sport=dip:dport
			// sip:sport=dport
			// sport=dport
			sip, sport, err := parseIP(tunnelSplitArr[0], defaultSip)
			if err != nil {
				fmt.Println("Error on", tunnelStr, ":", err)
				flag.Usage()
				os.Exit(-1)
			}
			dip, dport, err := parseIP(tunnelSplitArr[1], defaultDip)
			if err != nil {
				fmt.Println("Error on", tunnelStr, ":", err)
				flag.Usage()
				os.Exit(-1)
			}
			tunnels = append(tunnels, &Tunnel{
				Sip:   sip,
				Sport: sport,
				Dip:   dip,
				Dport: dport,
				Type:  t,
			})
		}
	}
	return tunnels
}


// parseIP
// 		ip:port
// 		port
func parseIP(ipstr string, defaultIP string) (ip string, port int, err error) {
	ipSplitArr := strings.Split(ipstr, ":")
	switch len(ipSplitArr) {
	case 1:
		// port
		p, err := strconv.ParseInt(ipSplitArr[0], 10, 64)
		if err != nil {
			return "", 0, ErrPortParseFailure
		}
		return defaultIP, int(p), nil
	case 2:
		// ip:port
		ip = ipSplitArr[0]
		if ipSplitArr[0] != defaultDip {
			i := net.ParseIP(ipSplitArr[0])
			if i == nil {
				return "", 0, ErrIPParseFailure
			}
			ip = i.String()
		}
		p, err := strconv.ParseInt(ipSplitArr[1], 10, 64)
		if err != nil {
			return "", 0, ErrPortParseFailure
		}
		return ip, int(p), nil
	}
	return "", 0, ErrIPParseFailure
}
