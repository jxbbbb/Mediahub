package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Http struct {
		IP   string
		Port int
		Mode string
	}
	Log struct {
		Level   string
		LogPath string `mapstructure:"logPath"`
	} `mapstructure:"log"`
	DependOn struct {
		ShortUrl struct {
			Address     string
			AccessToken string
		}
	}
}

var conf *Config

func InitConfig(filePath string, typ ...string) {
	v := viper.New()
	v.SetConfigFile(filePath)
	if len(typ) > 0 {
		v.SetConfigType(typ[0])
	}
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	conf = &Config{}
	err = v.Unmarshal(conf)
	if err != nil {
		log.Fatal(err)
	}

}

func GetConfig() *Config {
	return conf
}
