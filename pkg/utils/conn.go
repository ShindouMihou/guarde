package utils

import "github.com/rs/zerolog/log"

func EnsureClosure(f func() error) {
	err := f()
	if err != nil {
		log.Err(err).Stack().Msg("Failed to close a closeable.")
	}
}
