package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	
)

type HTTPServer struct {
	Addr string
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config {

	var configPath string

	configPath = os.Getenv("CONFIGPATH")

	if configPath == "" {
		flags:= flag.String("config", "", "path to the configure file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err){

		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cnf Config

	err:= cleanenv.ReadConfig(configPath, &cnf)
	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}

	return &cnf



}