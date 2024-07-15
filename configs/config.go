package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	SecretKey      string         `mapstructure:"SECRET_KEY"`
	Host           string         `mapstructure:"HOST"`
	GinServerPort  string         `mapstructure:"GIN_SERVER_PORT"`
	GrpcServerPort string         `mapstructure:"GRPC_SERVER_PORT"`
	DatabaseConfig DatabaseConfig `mapstructure:",squash"`
}

type DatabaseConfig struct {
	Port     string `mapstructure:"DB_PORT"`
	Host     string `mapstructure:"DB_HOST"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

func InitConfig(path string) (*Config, error) {
	var config Config
	if err := LoadConfig(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadConfig(path string, config *Config) error {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	viper.SetDefault("SECRET_KEY", "secret-key")
	viper.SetDefault("HOST", "localhost")
	viper.SetDefault("GIN_SERVER_PORT", "8081")
	viper.SetDefault("GRPC_SERVER_PORT", "50051")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "1702")
	viper.SetDefault("DB_NAME", "authentication")

	err = viper.Unmarshal(config)
	if err != nil {
		return fmt.Errorf("unable to decode into struct: %v", err)
	}

	return nil
}

func GetDatabaseCongig(path string) (*DatabaseConfig, error) {
	config, err := InitConfig(path)
	if err != nil {
		return nil, err
	}
	return &config.DatabaseConfig, nil
}
