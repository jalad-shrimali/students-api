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

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"` //struct embedding
	MongoURI        string
	MongoDatabase   string
	MongoCollection string
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse() //parse the flags

		configPath = *flags //flag is a pointer
		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) { //check if the file exists
		log.Fatalf("config file not found at %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg) //load the config file
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}
	return &cfg
}
