package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"guarde/internal/global"
	"guarde/internal/healthcheck"
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
		log.Info().Msg("Verbose mode is enabled.")
		log.Logger = logger.Logger().Level(zerolog.DebugLevel)
	}
	if len(config.Options) > 0 {
		log.Info().Int("len", len(config.Options)).Msg("Loaded additional options")
		global.Options = config.Options
	}
	log.Info().Int("rulesets", len(config.Ruleset)).Msg("Loaded rulesets.")
	if config.Proxy.Udp != nil {
		log.Info().Msg("Starting proxy server for UDP.")
		if config.Proxy.Udp.Fallback != nil {
			log.Info().Int("fallbacks", len(config.Proxy.Udp.Fallback.Addresses)).Msg("Loaded UDP fallback addresses.")
		}
		go config.Proxy.Udp.LaunchUdp(config)
	}
	if config.Proxy.Tcp != nil {
		log.Info().Msg("Starting proxy server for TCP.")
		if config.Proxy.Tcp.Fallback != nil {
			log.Info().Int("fallbacks", len(config.Proxy.Tcp.Fallback.Addresses)).Msg("Loaded TCP fallback addresses.")
		}
		go config.Proxy.Tcp.LaunchTcp(config)
	}
	if config.Healthcheck != nil {
		log.Info().Msg("Starting server for healthcheck.")
		go healthcheck.Launch(config.Healthcheck.Port)
	}
	select {}
}
