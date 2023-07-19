package fileutils

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func Open(file string) (*os.File, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return f, nil
}

var homeDirectory = ""

func GetHomeDir() string {
	if homeDirectory == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Panic().Err(err).Msg("Failed to get home directory path")
		}
		homeDirectory = home
	}
	return homeDirectory
}

func JoinHomePath(paths ...string) string {
	return filepath.Join(GetHomeDir(), filepath.Join(paths...))
}
