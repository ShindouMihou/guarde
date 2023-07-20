package udp

import (
	"guarde/pkg/utils"
	"net"
	"time"
)

func Request(addr string, body []byte) ([]byte, error) {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(3000 * time.Millisecond)
	err = conn.SetReadDeadline(deadline)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(body)
	if err != nil {
		return nil, err
	}
	defer utils.EnsureClosure(conn.Close)
	response := make([]byte, 1024)
	offset, _, err := conn.ReadFrom(response)
	if err != nil {
		return nil, err
	}
	return response[:offset], nil
}
