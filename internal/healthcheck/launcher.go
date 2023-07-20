package healthcheck

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/rs/zerolog/log"
)

func Launch(port uint16) {
	handler := Client{}
	err := gnet.Run(
		&handler,
		fmt.Sprintf("%s://0.0.0.0:%d", "tcp", port),
		gnet.WithReusePort(true),
		gnet.WithLogLevel(logging.ErrorLevel),
	)
	if err != nil {
		log.Err(err).Msg("Failed to start healthcheck service.")
		return
	}
	log.Info().Int("port", int(port)).Str("mode", "tcp").Msg("Running healthcheck service.")
}
