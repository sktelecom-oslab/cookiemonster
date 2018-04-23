package domain

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Resource struct {
	Kind     string
	Name     string
	Target   int64
}

type Namespace struct {
	Name     string
	Resource []Resource
}

type Config struct {
	Namespace []Namespace
	Interval int64
	Duration int64
	Slack    bool
}

var config *Config

func init() {
	path := "config"
	config = &Config{}
	if err := config.ReadConfig(path); err != nil {
		log.Println(err)
	}
}

func GetConfig() *Config {
	return config
}

func (c *Config) ReadConfig(path string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if path != "" {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Can't read config file: %s \n", err)
	}

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("config file format error: %s \n", err)
	}

	return nil
}
