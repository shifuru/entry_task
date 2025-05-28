package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Jwt   JwtConfig   `mapstructure:"jwt"`
	Db    DbConfig    `mapstructure:"db"`
	Redis RedisConfig `mapstructure:"redis"`
}

type JwtConfig struct {
	Key  string `mapstructure:"key"`
	Mode uint   `mapstructure:"mode"`
}
type DbConfig struct {
	Dsn         string `mapstructure:"dsn"`
	MaxIdle     int    `mapstructure:"max_idle"`
	MaxIdleTime int    `mapstructure:"max_idle_time"`
	MaxOpen     int    `mapstructure:"max_open"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

var ProjectConfig *Config

func init() {
	viper.SetConfigName("config.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %v \n", err))
	}

	ProjectConfig = &Config{}
	err = viper.Unmarshal(ProjectConfig)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %v \n", err))
	}
}
