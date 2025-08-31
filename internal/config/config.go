package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port     int `mapstructure:"port"`
		RpmLimit int `mapstructure:"rpmLimit"`
	} `mapstructure:"server"`

	Routes []Route `mapstructure:"routes"`

	JwtSecret string `mapstructure:"jwt-secret"`
}

type Route struct {
	ServiceName string `mapstructure:"serviceName"`
	Path        string `mapstructure:"path"`
	Endpoint    string `mapstructure:"endpoint"`
	Method      string `mapstructure:"method"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	fmt.Println(viper.GetString("server.name"))

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
