package config

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
)

// Config holds data for application configuration
type AppConfig struct {
	Server *Server   `yaml:"server,omitempty"`
	DB     *Database `yaml:"database,omitempty"`
}

// Database holds data for database configuration
type Database struct {
	Dsn string `yaml:"dsn,omitempty" envconfig: `
}

// Server holds data for server configuration
type Server struct {
	Network         string        `yaml:"network,omitempty" default:"3000"`
	Host            string        `yaml:"host,omitempty" default:"3000"`
	Port            string        `yaml:"port,omitempty" default:"3000"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout,omitempty" default:"10s"`
	ReadTimeout     time.Duration `yaml:"read_timeout,omitempty" default:"5s"`
	WriteTimeout    time.Duration `yaml:"write_timeout,omitempty" default:"10s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout,omitempty" default:"60s"`
}

// Load returns config from yaml and environment variables.
func Load(file string) (*AppConfig, error) {
	log.Infof("loading config file : %s \n", file)

	// default config
	var c AppConfig

	// load from YAML config file
	if rawcfg, err := ioutil.ReadFile(file); err == nil {
		if err := yaml.Unmarshal(rawcfg, &c); err != nil {
			log.Errorf("error on json marshall of config file : %s", file)
			return nil, err
		}
	} else {
		log.Errorf("error reading config file : %s", file)
		return nil, err
	}

	// if dsn still empty, throw error
	if c.DB.Dsn == "" {
		return nil, errors.New("database configuration is missing")
	}

	return &c, nil
}
