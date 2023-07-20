package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken     string 
	GptToken 		  string
	AssemblyToken     string

	Messages 		  Messages
}

type Messages struct {
	Responses
}


type Responses struct {
	Start             string `mapstructure:"start"`
	UnknownCommand    string `mapstructure:"unknown_command"`
}

func InitConfig() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err 
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err 
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err 
	}


	if err := parseEnv(&cfg); err != nil {
		return nil, err 
	}

	return &cfg, nil 
}

func parseEnv(cfg *Config) error {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	if err := viper.BindEnv("tgApiKey"); err != nil {
		return err 
	}

	if err := viper.BindEnv("gptApiKey"); err != nil {
		return err 
	}

	if err := viper.BindEnv("assemblyApiKey"); err != nil {
		return err 
	}

	cfg.TelegramToken = viper.GetString("tgApiKey")
	cfg.GptToken = viper.GetString("gptApiKey")
	cfg.AssemblyToken = viper.GetString("assemblyApiKey")

	return nil 
}