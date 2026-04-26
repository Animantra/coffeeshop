package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Log      `yaml:"logger"`
		Postgres `yaml:"postgres"`
		JWT      `yaml:"jwt"`
	}

	App struct {
		Name    string `yaml:"name"    env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Host string `yaml:"host" env:"HTTP_HOST"`
		Port int    `yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `yaml:"log_level" env:"LOG_LEVEL"`
	}

	Postgres struct {
		URL string `yaml:"url" env:"PG_URL"`
	}

	JWT struct {
		Secret     string `yaml:"secret"      env:"JWT_SECRET"`
		ExpireHours int   `yaml:"expire_hours" env:"JWT_EXPIRE_HOURS"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(dir)

	err = cleanenv.ReadConfig(dir+"/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
