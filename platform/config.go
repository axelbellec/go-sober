package platform

import (
	"fmt"

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

type SwaggerConfig struct {
	Host string `env:"SWAGGER_HOST" envDefault:"localhost:8080"`
}

type DatabaseConfig struct {
	SQL struct {
		FilePath string `env:"DB_FILE_PATH" envDefault:"db/sober.db"`
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
