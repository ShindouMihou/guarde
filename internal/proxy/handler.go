package proxy

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/rs/zerolog/log"
	"guarde/internal/healthcheck"
	"guarde/pkg/utils"
	"net"
	"strings"
	"time"
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
	start := time.Now()

	healthcheck.RequestsReceived.Increment()
	logger := log.With().Str("mode", s.mode).Logger()
	logger.Info().Str("addr", conn.RemoteAddr().String()).Msg("received connection")

	remoteAddr := conn.RemoteAddr().String()
	addr := strings.SplitN(remoteAddr, ":", 2)[0]
	ip := net.ParseIP(addr)

	logger = logger.With().Str("addr", remoteAddr).Logger()
	if !ip.IsLoopback() {
		allow := s.config.IsAllowed(ip.String())
		if !allow {
			healthcheck.RequestsRejected.Increment()
			logger.Warn().Str("lt", time.Since(start).String()).Msg("Not permitted address.")
			return gnet.Close
		}
	}
	message, err := conn.Next(-1)
	if err != nil {
		healthcheck.RequestsErrored.Increment()
		logger.Err(err).Msg("Failed to read request body.")
		return gnet.Close
	}
	healthcheck.AverageRequestSize.Add(len(message))
	if s.config.Verbose {
		logger.Debug().Str("t", "qry").Msg(utils.Simplify(string(message)))
	}
	forwardStart := time.Now()
	response, err := s.forwarder(s.connection.Forward, message)
	if err != nil {
		healthcheck.RequestsErrored.Increment()

		forwardLatency := time.Since(forwardStart)
		healthcheck.AverageForwardLatency.Add(int(forwardLatency.Milliseconds()))

		logger.Err(err).Str("lt", forwardLatency.String()).Msg("Failed to request from forward address.")

		if s.connection.Fallback != nil {
			ok := false
			for _, fallback := range s.connection.Fallback.Addresses {
				fallbackStart := time.Now()
				logger.Debug().Str("fb", fallback).Msg("Requesting from fallback address.")
				response, err = s.forwarder(fallback, message)

				if err != nil {
					fallbackLatency := time.Since(fallbackStart)
					healthcheck.AverageFallbackLatency.Add(int(fallbackLatency.Milliseconds()))
					healthcheck.FallbackFailures.Increment()

					logger.Err(err).Str("lt", fallbackLatency.String()).Str("fb", fallback).Msg("Failed to request from fallback address.")
					continue
				}

				latency := time.Since(forwardStart)
				healthcheck.AverageFallbackLatency.Add(int(latency.Milliseconds()))

				ok = true
				break
			}
			if !ok {
				return gnet.Close
			}
		} else {
			return gnet.Close
		}
	}

	forwardLatency := time.Since(forwardStart)
	healthcheck.AverageForwardLatency.Add(int(forwardLatency.Milliseconds()))

	if s.config.Verbose {
		logger.Debug().Str("t", "qry").Str("lt", forwardLatency.String()).Msg(utils.Simplify(string(response)))
	}
	healthcheck.AverageResponseSize.Add(len(response))
	_, err = conn.Write(response)
	if err != nil {
		healthcheck.RequestsErrored.Increment()
		log.Err(err).Msg("Failed to reply to client.")
		return gnet.Close
	}

	requestLatency := time.Since(start)
	healthcheck.AverageRequestLatency.Add(int(requestLatency.Milliseconds()))
	healthcheck.RequestsForwarded.Increment()

	logger.Info().Str("lt", requestLatency.String()).Msg("forwarded connection")
	return gnet.None
}
