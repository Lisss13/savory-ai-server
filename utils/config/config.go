package config

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// envVarRegex для поиска переменных окружения в формате ${VAR} или ${VAR:-default}
var envVarRegex = regexp.MustCompile(`\$\{([^}:]+)(?::-([^}]*))?\}`)

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

// expandEnvVars заменяет переменные окружения в строке.
// Поддерживает форматы: ${VAR} и ${VAR:-default}
func expandEnvVars(content []byte) []byte {
	return envVarRegex.ReplaceAllFunc(content, func(match []byte) []byte {
		groups := envVarRegex.FindSubmatch(match)
		if len(groups) < 2 {
			return match
		}

		varName := string(groups[1])
		value := os.Getenv(varName)

		// Если переменная не установлена и есть значение по умолчанию
		if value == "" && len(groups) > 2 && len(groups[2]) > 0 {
			value = string(groups[2])
		}

		return []byte(value)
	})
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

	// Подстановка переменных окружения
	file = expandEnvVars(file)

	err = toml.Unmarshal(file, &contents)

	return contents, err
}

func validateConfig(cfg *Config) {
	// add validation if needed
	if cfg.Anthropic.APIKey == "" {
		log.Panic().Msg("Anthropic API key is required")
	}

	if cfg.App.Port == "" {
		log.Panic().Msg("Port is required")
	}

	if cfg.DB.Postgres.DSN == "" {
		log.Panic().Msg("Postgres DSN is required")
	}
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

	validateConfig(config)

	return config
}
