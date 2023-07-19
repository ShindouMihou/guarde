package proxy

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/rs/zerolog/log"
	"guarde/pkg/utils"
	"net"
	"strings"
)

type Client struct {
	gnet.BuiltinEventEngine
	config     *Configuration
	connection *Connection
	mode       string
	forwarder  Forwarder
}

type Forwarder func(addr string, body []byte) ([]byte, error)

func (s *Client) OnTraffic(conn gnet.Conn) gnet.Action {
	logger := log.With().Str("mode", s.mode).Logger()
	logger.Info().Str("addr", conn.RemoteAddr().String()).Msg("received connection")

	addr := strings.SplitN(conn.RemoteAddr().String(), ":", 2)[0]
	ip := net.ParseIP(addr)

	logger = logger.With().Str("addr", ip.String()).Logger()
	if !ip.IsLoopback() {
		allow := s.config.IsAllowed(ip.String())
		if !allow {
			logger.Warn().Msg("Not permitted address.")
			return gnet.Close
		}
	}
	message, err := conn.Next(-1)
	if err != nil {
		logger.Err(err).Msg("Failed to read request body.")
		return gnet.Close
	}
	if s.config.Verbose {
		logger.Debug().Str("t", "qry").Msg(utils.Simplify(string(message)))
	}
	response, err := s.forwarder(s.connection.Forward, message)
	if err != nil {
		logger.Err(err).Msg("Failed to request from forward address.")
		return gnet.Close
	}
	if s.config.Verbose {
		logger.Debug().Str("t", "qry").Msg(utils.Simplify(string(response)))
	}
	_, err = conn.Write(response)
	if err != nil {
		log.Err(err).Msg("Failed to reply to client.")
		return gnet.Close
	}
	logger.Info().Msg("forwarded connection")
	return gnet.None
}
