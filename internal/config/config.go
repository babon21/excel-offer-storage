package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Username string
		Password string
		Host     string
		Port     string
		DbName   string
	}
}

func Init() Config {
	var config Config

	viper.AutomaticEnv()

	config.Server.Port = viper.GetString("SERVER_PORT")
	config.Database.Username = viper.GetString("DATABASE_USER")
	config.Database.Password = viper.GetString("DATABASE_PASSWORD")
	config.Database.Host = viper.GetString("DATABASE_HOST")
	config.Database.Port = viper.GetString("DATABASE_PORT")
	config.Database.DbName = viper.GetString("DATABASE_NAME")

	return config
}
