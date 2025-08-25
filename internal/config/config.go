package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port         int           `mapstructure:"port"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		RpmLimit     int           `mapstructure:"rpmLimit"`
	} `mapstructure:"server"`

	Routes []Route `mapstructure:"routes"`
}

type Route struct {
	ServiceName string `mapstructure:"serviceName"`
	Path        string `mapstructure:"path"`
	Endpoint    string `mapstructure:"endpoint"`
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

func (c Config) Print() {
	println("Server:")
	println("  Port:", c.Server.Port)
	println("  ReadTimeout:", c.Server.ReadTimeout.String())
	println("  WriteTimeout:", c.Server.WriteTimeout.String())
	println("Routes:")
	for i, route := range c.Routes {
		println("  Route", i+1, ":")
		println("    Path:", route.Path)
		println("    ServiceName:", route.ServiceName)
		println("    Endpoints:")
		for _, endpoint := range route.Endpoint {
			println("      -", endpoint)
		}
	}
}
