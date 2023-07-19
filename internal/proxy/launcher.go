package proxy

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/rs/zerolog/log"
	"guarde/pkg/tcp"
	"guarde/pkg/udp"
)

func (connection *Connection) Launch(mode string, forwarder Forwarder, config *Configuration) {
	handler := Client{connection: connection, config: config, forwarder: forwarder, mode: mode}
	err := gnet.Run(
		&handler,
		fmt.Sprintf("%s://0.0.0.0:%d", mode, connection.Port),
		gnet.WithReusePort(true),
		gnet.WithLogLevel(logging.ErrorLevel),
	)
	if err != nil {
		log.Err(err).Str("mode", mode).Msg("Failed to start forwarding proxy.")
		return
	}
	log.Info().Int("port", int(connection.Port)).Str("mode", mode).Msg("Running forwarding proxy server.")
}

func (connection *Connection) LaunchUdp(config *Configuration) {
	connection.Launch("udp", udp.Request, config)
}

func (connection *Connection) LaunchTcp(config *Configuration) {
	connection.Launch("tcp", tcp.Request, config)
}
