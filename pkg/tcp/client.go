package tcp

import (
	"guarde/internal/global"
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
	_, err = conn.Write(body)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(time.Duration(global.ReadDeadline.GetDefault(1024)+256) * time.Millisecond)
	err = conn.SetReadDeadline(deadline)
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
