package app

import (
	"embed"
	"errors"
	"io/fs"
	"path"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed embed
	embedFS      embed.FS
	unwrapFSOnce sync.Once
	unwrappedFS  fs.FS
)

func FS() fs.FS {
	unwrapFSOnce.Do(func() {
		fsys, err := fs.Sub(embedFS, "embed")
		if err != nil {
			panic(err)
		}
		unwrappedFS = fsys
	})
	return unwrappedFS
}

// Config holds data for application configuration
type AppConfig struct {
	Server *Server   `yaml:"server,omitempty"`
	DB     *Database `yaml:"database,omitempty"`

	Debug bool `yaml:"debug,omitempty"`
	Env   bool `yaml:"env,omitempty"`
}

// Database holds data for database configuration
type Database struct {
	DSN string `yaml:"dsn,omitempty"`
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
func LoadConfigFile(fsys fs.FS, service, env string) (*AppConfig, error) {
	// default config
	var c AppConfig

	// load from YAML config file
	if rawcfg, err := fs.ReadFile(fsys, path.Join("config", env+".yaml")); err == nil {
		if err := yaml.Unmarshal(rawcfg, &c); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// if dsn still empty, throw error
	if c.DB.DSN == "" {
		return nil, errors.New("database configuration is missing")
	}

	return &c, nil
}
