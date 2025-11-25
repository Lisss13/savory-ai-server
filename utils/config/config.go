package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// app struct config
type app = struct {
	Name        string        `toml:"name"`
	Host        string        `toml:"host"`
	Port        string        `toml:"port"`
	PrintRoutes bool          `toml:"print-routes"`
	Prefork     bool          `toml:"prefork"`
	Production  bool          `toml:"production"`
	IdleTimeout time.Duration `toml:"idle-timeout"`
	TLS         struct {
		Enable   bool   `toml:"enable"`
		CertFile string `toml:"cert-file"`
		KeyFile  string `toml:"key-file"`
	}
	ChatServiceIrl string `toml:"chat_service_url"`
}

// db struct config
type db = struct {
	Postgres struct {
		DSN string `toml:"dsn"`
	}
}

// log struct config
type logger = struct {
	TimeFormat string        `toml:"time-format"`
	Level      zerolog.Level `toml:"level"`
	Prettier   bool          `toml:"prettier"`
}

// anthropic config
type anthropic = struct {
	APIKey string `toml:"api_key"`
}

// middleware
type middleware = struct {
	Compress struct {
		Enable bool
		Level  compress.Level
	}

	Recover struct {
		Enable bool
	}

	Monitor struct {
		Enable bool
		Path   string
	}

	Pprof struct {
		Enable bool
	}

	Limiter struct {
		Enable     bool
		Max        int
		Expiration time.Duration `toml:"expiration_seconds"`
	}

	FileSystem struct {
		Enable bool
		Browse bool
		MaxAge int `toml:"max_age"`
		Index  string
		Root   string
	}

	Jwt struct {
		Secret     string        `toml:"secret"`
		Expiration time.Duration `toml:"expiration_seconds"`
	}
}

type Config struct {
	App        app
	DB         db
	Logger     logger
	Middleware middleware
	Anthropic  anthropic
}

// func to parse config
func ParseConfig(name string, debug ...bool) (*Config, error) {
	var (
		contents *Config
		file     []byte
		err      error
	)

	if len(debug) > 0 {
		file, err = os.ReadFile(name)
	} else {
		_, b, _, _ := runtime.Caller(0)
		// get base path
		path := filepath.Dir(filepath.Dir(filepath.Dir(b)))
		file, err = os.ReadFile(filepath.Join(path, "./config/", name+".toml"))
	}

	if err != nil {
		return &Config{}, err
	}

	err = toml.Unmarshal(file, &contents)

	return contents, err
}

// initialize config
func NewConfig() *Config {
	// Check if CONFIG_FILE environment variable is set
	configFile := os.Getenv("CONFIG_FILE")
	var config *Config
	var err error

	if configFile != "" {
		// If CONFIG_FILE is set, use it directly
		config, err = ParseConfig(configFile, true)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse config from CONFIG_FILE")
		}
	} else {
		// Otherwise, use the default config file
		config, err = ParseConfig("config")
	}

	if err != nil && !fiber.IsChild() {
		// panic if config is not found
		log.Panic().Err(err).Msg("config not found")
	}

	return config
}
