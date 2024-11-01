package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type DBConfig struct {
	DBHost     string `yaml:"db_host" env-required:"true"`
	DBPort     string `yaml:"db_port" env-required:"true"`
	DBUser     string `yaml:"db_user" env-required:"true"`
	DBPassword string `yaml:"db_pass" env-required:"true"`
	DBName     string `yaml:"db_name" env-required:"true"`
}

type Config struct {
	Env         string `yaml:"env" env:"Env" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	DBConfig    `yaml:"db_config"`
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		// if path empty check whether it is passed as arguments
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	// Check whether file exists at that location
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists at given path %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config file %s", err.Error())
	}
	return &cfg
}
