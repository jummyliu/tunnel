package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Tunnel struct
type Tunnel struct {
	Sip   string
	Sport int
	Dip   string
	Dport int
	Type  TunnelType
}

// TunnelType the type of tunnel
type TunnelType uint8

// Support tunnel's type
const (
	_ = iota
	TunnelTypeALL
	TunnelTypeTCP
	TunnelTypeUDP
)

// Build tunnel
func (t *Tunnel) Build() {
	var wg sync.WaitGroup
	var networks []string
	switch t.Type {
	case TunnelTypeTCP:
		networks = append(networks, "tcp")
	case TunnelTypeUDP:
		networks = append(networks, "udp")
	case TunnelTypeALL:
		networks = append(networks, "tcp", "udp")
	}

	for _, network := range networks {
		go func(network string) {
			wg.Add(1)
			defer wg.Done()
			sAddr := fmt.Sprintf("%s:%d", t.Sip, t.Sport)
			dAddr := fmt.Sprintf("%s:%d", t.Dip, t.Dport)

			switch network {
			case "tcp":
				ln, err := net.Listen(network, sAddr)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer ln.Close()
				fmt.Printf("proxy (%s) %s => %s\n", network, sAddr, dAddr)
				for {
					conn, err := ln.Accept()
					if err != nil {
						fmt.Println(err)
						continue
					}
					go func(network string) {
						defer conn.Close()
						tunnelConn, err := net.Dial(network, dAddr)
						if err != nil {
							fmt.Println(err)
							return
						}
						defer tunnelConn.Close()
						relay(conn, tunnelConn)
					}(network)
				}
			case "udp":
				conn, err := net.ListenPacket(network, sAddr)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer conn.Close()
				fmt.Printf("proxy (%s) %s => %s\n", network, sAddr, dAddr)
				buf := make([]byte, 1024)
				for {
					_, raddr, err := conn.ReadFrom(buf)
					if err != nil {
						fmt.Println(err)
						continue
					}
					pc, err := net.ListenPacket(network, "")
					pc.SetReadDeadline(time.Now().Add(5 * time.Second))
					tAddr, err := net.ResolveUDPAddr(network, dAddr)
					pc.WriteTo(buf, tAddr)
					go func() {
						buf := make([]byte, 1024)
						for {
							_, pcAddr, err := pc.ReadFrom(buf)
							if err != nil {
								continue
							}
							if pcAddr.String() != tAddr.String() {
								continue
							}
							pc.WriteTo(buf, raddr)
						}
					}()
				}
			}
		}(network)
	}
	wg.Wait()
}

func relay(left, right net.Conn) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)
	go func() {
		n, err := io.Copy(right, left)
		right.SetDeadline(time.Now())
		left.SetDeadline(time.Now())
		ch <- res{n, err}
	}()

	n, err := io.Copy(left, right)
	right.SetDeadline(time.Now())
	left.SetDeadline(time.Now())
	rs := <-ch
	if err == nil {
		err = rs.Err
	}
	return n, rs.N, err
}
