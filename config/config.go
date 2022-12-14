package config

import "github.com/spf13/viper"

var Cfg *Config

type Config struct {
	Server Server `mapstructure:"server"`
}

type Server struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
}

func LoadConfig(path string) {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(Cfg)
	if err != nil {
		panic(err)
	}
}
