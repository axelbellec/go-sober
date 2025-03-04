package platform

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

const (
	TEST_ENV  = "test"
	LOCAL_ENV = "local"
	PROD_ENV  = "prod"
)

type LoggerConfig struct {
	Level  string `env:"LOG_LEVEL" envDefault:"INFO"`
	Format string `env:"LOG_FORMAT" envDefault:"json"`
}

type LLMConfig struct {
	Groq struct {
		BaseURL string `env:"GROQ_BASE_URL" envDefault:"https://api.groq.com/openai/v1"`
		Model   string `env:"GROQ_MODEL" envDefault:"llama-3.2-1b-preview"`
		APIKey  string `env:"GROQ_API_KEY"`
	}
}

type DatabaseConfig struct {
	SQL struct {
		FilePath        string        `env:"DB_FILE_PATH" envDefault:"db/sober.db"`
		MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
		MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
		ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"1h"`
	}
}

type AuthConfig struct {
	JWT struct {
		Secret string `env:"JWT_SECRET"`
	}
}

type Config struct {
	AppName     string `env:"APP_NAME"`
	AppVersion  string `env:"APP_VERSION" envDefault:"unknown"`
	Environment string `env:"ENVIRONMENT"`
	Port        string `env:"PORT" envDefault:"8080"`
	Logger      LoggerConfig
	Database    DatabaseConfig
	Auth        AuthConfig
	LLM         LLMConfig
}

var AppConfig *Config = nil

func initConfig() {
	cfg := Config{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		panic(fmt.Sprintf("could not parse config: %v", err))
	}

	AppConfig = &cfg
}
