package proxy

import (
	"gopkg.in/yaml.v3"
	"guarde/pkg/fileutils"
	"guarde/pkg/utils"
	"io"
)

type Configuration struct {
	Proxy   Proxy  `yaml:"proxy"`
	Ruleset []Rule `yaml:"ruleset"`
	Timeout uint16 `yaml:"timeout"`
	Verbose bool   `yaml:"verbose"`
	Allow   Allow  `yaml:"allow"`
}

type Rule map[string]string

type Proxy struct {
	Udp *Connection `yaml:"udp,omitempty"`
	Tcp *Connection `yaml:"tcp,omitempty"`
}

type Allow struct {
	PropertyNotFound bool `yaml:"property_not_found"`
}

type Connection struct {
	Forward string `yaml:"forward"`
	Port    uint16 `yaml:"port"`
}

func New(dir string) (*Configuration, error) {
	f, err := fileutils.Open(dir)
	if err != nil {
		return nil, err
	}
	defer utils.EnsureClosure(f.Close)
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var config Configuration
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
