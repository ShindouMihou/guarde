// Package udp server client was originally built to run the UDP server until I discovered that
// a great library (gnet.host) already existed, so this was abandoned, but still kept for
// archival purposes.
package udp

import (
	"github.com/rs/zerolog/log"
	"net"
)

type Context struct {
	Connection    *net.UDPConn
	Message       []byte
	RemoteAddress *net.UDPAddr
}

func (context *Context) Reply(body []byte) error {
	_, err := context.Connection.WriteToUDP(body, context.RemoteAddress)
	if err != nil {
		return err
	}
	return nil
}

type Handler func(ctx *Context)

func Listen(port uint16, handler Handler) error {
	address := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: int(port),
	}
	server, err := net.ListenUDP("udp", &address)
	if err != nil {
		return err
	}
	for {
		buffer := make([]byte, 1024)
		offset, raddr, err := server.ReadFromUDP(buffer)
		if err != nil {
			log.Err(err).Msg("Failed to read from UDP server.")
			continue
		}

		message := buffer[:offset]
		go handler(&Context{Connection: server, Message: message, RemoteAddress: raddr})
	}
}
