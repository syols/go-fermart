package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"strconv"
)

type Option func(s *Config)

type Config struct {
	ServerAddress            Address `yaml:"address"`
	AccuralAddress           Address `yaml:"accural"`
	DatabaseConnectionString string  `yaml:"database"`
	Key                      *string `yaml:"key" env:"KEY"`
}

type Address struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

func NewConfig() (settings Config, err error) {
	if err := settings.setDefault("/Users/s.olshanskiy/Documents/github/projects/go-fermart/develop.yml"); err != nil {
		return Config{}, err
	}
	return settings.setFromOptions(NewEnvironmentVariables().Options()...), nil
}

func (s *Config) Address() string {
	return fmt.Sprintf("%s:%d", s.ServerAddress.Host, s.ServerAddress.Port)
}

func (s *Config) setFromOptions(options ...Option) Config {
	for _, fn := range options {
		fn(s)
	}
	return *s
}

func (s *Config) setDefault(configPath string) error {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(file, s); err != nil {
		return err
	}
	return nil
}

func (s *Config) String() (result string) {
	if marshal, err := yaml.Marshal(s); err != nil {
		result = string(marshal)
	}
	return
}

func withServerAddress(address string) Option {
	return func(s *Config) {
		if host, port, err := net.SplitHostPort(address); err == nil {
			if port, err := strconv.ParseUint(port, 0, 16); err == nil {
				s.ServerAddress.Host = host
				s.ServerAddress.Port = uint16(port)
			}
		}
	}
}

func withAccuralAddress(address string) Option {
	return func(s *Config) {
		if host, port, err := net.SplitHostPort(address); err == nil {
			if port, err := strconv.ParseUint(port, 0, 16); err == nil {
				s.AccuralAddress.Host = host
				s.AccuralAddress.Port = uint16(port)
			}
		}
	}
}

func withKey(value string) Option {
	return func(s *Config) {
		s.Key = &value
	}
}

func withDatabase(value string) Option {
	return func(s *Config) {
		s.DatabaseConnectionString = value
	}
}
