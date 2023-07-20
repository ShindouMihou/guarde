package udp

import (
	"bytes"
	"guarde/internal/global"
	"guarde/pkg/utils"
	"io"
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
	deadline := time.Now().Add(time.Duration(global.ReadDeadline.GetDefault(1024)) * time.Millisecond)
	err = conn.SetReadDeadline(deadline)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(body)
	if err != nil {
		return nil, err
	}
	defer utils.EnsureClosure(conn.Close)
	var buffer bytes.Buffer
	size, err := io.Copy(&buffer, conn)
	if err != nil {
		return nil, err
	}
	resp := make([]byte, size)
	copy(resp, buffer.Bytes())
	return resp, nil
}
