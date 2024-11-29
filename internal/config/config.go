package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// Config struct holds the application configuration
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	HTTP     HTTPConfig     `yaml:"http"`
	Logging  LoggerConfig   `yaml:"logging"`
}

type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         string        `yaml:"sslmode"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_life_time"`
}

type LoggerConfig struct {
	FilePath string `yaml:"filePath"`
	Encoding string `yaml:"encoding"`
	Level    string `yaml:"level"`
	Logger   string `yaml:"logger"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

var (
	once   sync.Once
	config *Config
)

func Load() (*Config, error) {
	var loadErr error
	once.Do(func() {
		fpath := os.Getenv("CONFIG_FILE")
		if fpath == "" {
			_, b, _, _ := runtime.Caller(0)
			root := filepath.Join(filepath.Dir(b), "../..")
			fpath = filepath.Join(root, "config.yml")
		}

		absPath, err := filepath.Abs(fpath)

		if err != nil {
			loadErr = fmt.Errorf("failed to resolve config file path: %w", err)
			return
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			loadErr = fmt.Errorf("config file does not exist: %s", absPath)
			return
		}

		file, err := os.Open(absPath)
		if err != nil {
			loadErr = fmt.Errorf("failed to open config file: %w", err)
			return
		}
		defer file.Close()

		var cfg Config
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&cfg); err != nil {
			loadErr = fmt.Errorf("failed to decode config file: %w", err)
			return
		}

		if err := validateConfig(&cfg); err != nil {
			loadErr = fmt.Errorf("config validation failed: %w", err)
			return
		}

		config = &cfg

		http_port, err := strconv.Atoi(os.Getenv("PORT"))
		if http_port != 0 && err == nil {
			config.HTTP.Port = http_port

		}
	})

	return config, loadErr
}

func validateConfig(cfg *Config) error {
	if cfg.Database.Host == "" {
		return fmt.Errorf("database.host is missing")
	}
	if cfg.Database.Port <= 0 {
		return fmt.Errorf("database.port is invalid or missing")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("database.user is missing")
	}
	if cfg.Database.Password == "" {
		return fmt.Errorf("database.password is missing")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("database.dbname is missing")
	}
	if cfg.Database.SSLMode == "" {
		return fmt.Errorf("database.sslmode is missing")
	}

	if cfg.HTTP.Host == "" {
		return fmt.Errorf("http.host is missing")
	}
	if cfg.HTTP.Port <= 0 {
		return fmt.Errorf("http.port is invalid or missing")
	}

	if cfg.Logging.Logger == "" {
		return fmt.Errorf("Logging.Logger is missing")
	}
	if cfg.Logging.Encoding == "" {
		return fmt.Errorf("Logging.Encoding is missing")
	}

	if cfg.Logging.FilePath == "" {
		return fmt.Errorf("Logging.FilePath is missing")
	}
	if cfg.Logging.Level == "" {
		return fmt.Errorf("Logging.Level is missing")
	}
	return nil
}
