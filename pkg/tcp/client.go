package tcp

import (
	"guarde/internal/global"
	"guarde/pkg/utils"
	"io"
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
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(global.ReadDeadline.GetDefault(1024)) * time.Millisecond))
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(body)
	if err != nil {
		return nil, err
	}
	response, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	return response, nil
}
