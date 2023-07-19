package tcp

import (
	"guarde/pkg/utils"
	"net"
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
