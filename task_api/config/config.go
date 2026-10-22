package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server  ServerConfig  `mapstructure:"server"`
		DB      DBConfig      `mapstructure:"db"`
		Redis   RedisConfig   `mapstructure:"redis"`
		Log     LogConfig     `mapstructure:"log"`
		Timeout TimeoutConfig `mapstructure:"timeout"`
	}

	ServerConfig struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
		Env  string `mapstructure:"env"`
	}

	DBConfig struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	RedisConfig struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		TTL  int    `mapstructure:"ttl"`
	}

	LogConfig struct {
		Level string `mapstructure:"level"`
	}

	TimeoutConfig struct {
		RequestSeconds int `mapstructure:"request_seconds"`
	}
)

func (db DBConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode,
	)
}

func (redis RedisConfig) GetRedis() string {
	return fmt.Sprintf("%s:%d", redis.Host, redis.Port)
}

var (
	once           sync.Once
	configInstance *Config
	configErr      error
)

func GetConfig() (*Config, error) {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			configErr = err
			return
		}

		if err := viper.Unmarshal(&configInstance); err != nil {
			configErr = err
			return
		}
	})

	return configInstance, configErr
}
