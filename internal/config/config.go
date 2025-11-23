package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host     string         `yaml:"host" env-default:"localhost"`
	Port     int            `yaml:"port"`
	Postgres PostgresConfig `yaml:"postgres"`
	Timeout  time.Duration  `yaml:"timeout" env-default:"300ms"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	MaxConns int32  `yaml:"max_conns"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func (c *Config) LoadEnv() {
	if c.Postgres.User == "" {
		c.Postgres.User = os.Getenv("POSTGRES_USER")
	}
	if c.Postgres.Password == "" {
		c.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	}
	if c.Postgres.DBName == "" {
		c.Postgres.DBName = os.Getenv("POSTGRES_DB")
	}
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
