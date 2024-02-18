package config

import (
	"log"
	"log/slog"
	"os"
	"time"
	"url-shortener/internal/lib/logger/handlers/slogpretty"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Env         string `yaml:"env" env-defaul:"local" env-required:"true"`
		StoragePath string `yaml:"storage_path" env-required:"true"`
		LoggerPath  string `yaml:"logger_path"`
		Log         Log    `yaml:"log"`
		HttpServer  `yaml:"http_server" `
		Clients     ClientConfig `yaml:"clients"`
		AppSecret   string       `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
	}

	HttpServer struct {
		Address         string        `yaml:"address" env-default:"localhost:8080"`
		Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
		User            string        `yaml:"user" env-required:"true"`
		Password        string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
	}

	Client struct {
		Address      string        `yaml:"address"`
		Timeout      time.Duration `yaml:"timeout"`
		RetriesCount int           `yaml:"retries_count"`
		Insecure     bool          `yaml:"insecure"`
	}
	ClientConfig struct {
		SSO Client `yaml:"sso"`
	}

	Log struct {
		Slog Slog `yaml:"slog"`
	}
	Slog struct {
		Level     slog.Level              `yaml:"level"`
		AddSource bool                    `yaml:"add_source"`
		Format    slogpretty.FieldsFormat `yaml:"format"` // json, text or pretty
		Pretty    PrettyLog               `yaml:"pretty"`
	}
	PrettyLog struct {
		FieldsFormat slogpretty.FieldsFormat `yaml:"fields_format"` // json, json-indent or yaml
		Emoji        bool                    `yaml:"emoji"`
		TimeLayout   string                  `yaml:"time_layout"`
	}
)

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
