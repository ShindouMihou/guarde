package tcp

import (
	"guarde/pkg/utils"
	"net"
	"time"
)

func Request(addr string, body []byte) ([]byte, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	defer utils.EnsureClosure(conn.Close)
	err = conn.SetReadDeadline(time.Now().Add(3000 * time.Millisecond))
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(body)
	if err != nil {
		return nil, err
	}
	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
