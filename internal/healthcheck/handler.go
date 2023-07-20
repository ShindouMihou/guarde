package healthcheck

import (
	"encoding/json"
	"github.com/panjf2000/gnet/v2"
	"github.com/rs/zerolog/log"
)

type Client struct {
	gnet.BuiltinEventEngine
}

func (s *Client) OnTraffic(conn gnet.Conn) gnet.Action {
	logger := log.With().Str("mode", "tcp").Str("addr", conn.RemoteAddr().String()).Logger()
	logger.Info().Msg("received healthcheck connection")

	report := report()
	bytes, err := json.Marshal(report)
	if err != nil {
		logger.Err(err).Msg("Failed to marshal healthcheck report.")
		_, _ = conn.Write([]byte("failed to marshal healthcheck report."))
		return gnet.Close
	}
	_, err = conn.Write(bytes)
	if err != nil {
		logger.Err(err).Msg("Failed to write to healthcheck.")
	}
	return gnet.Close
}
