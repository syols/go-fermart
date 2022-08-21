package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Option func(s *Config)

type Config struct {
	ServerAddress  Address `yaml:"address"`
	AccrualAddress Address `yaml:"accrual"`
	DatabaseURL    string  `yaml:"database"`
	Sign           string  `yaml:"sign"`
}

type Address struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

func NewConfig() (settings Config, err error) {
	if err := settings.setDefault("develop.yml"); err != nil {
		return Config{}, err
	}
	return settings.set(NewEnvironmentVariables().Options()...), nil
}

func (s *Config) set(options ...Option) Config {
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

func withAccrualAddress(address string) Option {
	return func(s *Config) {
		if host, port, err := net.SplitHostPort(address); err == nil {
			if port, err := strconv.ParseUint(port, 0, 16); err == nil {
				s.AccrualAddress.Host = host
				s.AccrualAddress.Port = uint16(port)
			}
		}
	}
}

func withDatabase(value string) Option {
	return func(s *Config) {
		s.DatabaseURL = value
	}
}

func withSign(value string) Option {
	return func(s *Config) {
		s.Sign = value
	}
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
