package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver           string `mapstructure:"DB_DRIVER"`
	DBSource           string `mapstructure:"DB_SOURCE"`
	ServerAddress      string `mapstructure:"SERVER_ADDRESS"`
	GRPCServerAddress  string `mapstructure:"GRPC_SERVER_ADDRESS"`
	SymmetricKey       string `mapstructure:"SYMMETRIC_KEY"`
	RedisServerAddress string `mapstructure:"REDIS_SERVER_ADDRESS"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	er := viper.ReadInConfig()
	if er != nil {
		return nil, er
	}

	errr := viper.Unmarshal(&config)
	if errr != nil {
		return nil, errr
	}

	return config, nil

}
