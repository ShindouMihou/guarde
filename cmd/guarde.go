package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"guarde/internal/proxy"
	"guarde/pkg/fileutils"
	"os"
)

var ConfigDirectory = fileutils.JoinHomePath(".guarde", "config.yml")

func main() {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Stack()
	log.Logger = logger.Logger().Level(zerolog.InfoLevel)

	log.Info().Str("path", ConfigDirectory).Msg("Reading for configuration file")
	config, err := proxy.New(ConfigDirectory)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read configuration file.")
	}
	if config.Verbose {
		log.Logger = logger.Logger().Level(zerolog.DebugLevel)
	}
	if config.Proxy.Udp != nil {
		log.Info().Msg("Starting proxy server for UDP.")
		go config.Proxy.Udp.LaunchUdp(config)
	}
	if config.Proxy.Tcp != nil {
		log.Info().Msg("Starting proxy server for TCP.")
		go config.Proxy.Tcp.LaunchTcp(config)
	}
	select {}
}
