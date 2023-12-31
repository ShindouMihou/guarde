package udp

import (
	"guarde/internal/global"
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
	defer utils.EnsureClosure(conn.Close)
	_, err = conn.Write(body)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(time.Duration(global.ReadDeadline.GetDefault(1024)+256) * time.Millisecond)
	err = conn.SetReadDeadline(deadline)
	if err != nil {
		return nil, err
	}
	response := make([]byte, global.BufferSize.GetDefault(1024))
	offset, _, err := conn.ReadFrom(response)
	if err != nil {
		return nil, err
	}
	return response[:offset], nil
}
