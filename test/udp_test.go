package test

import (
	"net"
	"testing"
)

// $ tunnel -u 10081=10082
//
// Then you can access udp://localhost:10082 using udp://ip:10081

// TestUDPServe
//
// Command: go test -v -run TestUDPServe$
func TestUDPServe(t *testing.T) {
	conn, err := net.ListenPacket("udp", "localhost:10082")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		_, raddr, err := conn.ReadFrom(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(raddr, string(buf))
		conn.WriteTo(buf, raddr)
	}
}

// TestUDPClient 
// 
// Command: go test -v -run TestUDPClient$
func TestUDPClient(t *testing.T) {
	pc, err := net.ListenPacket("udp", "")
	if err != nil {
		t.Fatal(err)
	}
	defer pc.Close()

	buf := make([]byte, 1024)
	for {
		addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:10081")
		pc.WriteTo([]byte("hello world"), addr)

		_, addr2, _ := pc.ReadFrom(buf)

		t.Log(addr2, string(buf))
	}
}