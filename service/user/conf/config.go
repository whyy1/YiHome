package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB_SOURCE        string `mapstructure:"DB_SOURCE"`
	REDIS_SOURCE     string `mapstructure:"REDIS_SOURCE"`
	REDIS_PASSWORD   string `mapstructure:"REDIS_PASSWORD"`
	ACCESS_KEYID     string `mapstructure:"ACCESS_KEYID"`
	ACCESS_KEYSECRET string `mapstructure:"ACCESS_KEYSECRET"`
	CONSUL_SOURCE    string `mapstructure:"CONSUL_SOURCE"`
	CONSUL_ADDRESS   string `mapstructure:"CONSUL_ADDRESS"`
	CONSUL_VERSION   string `mapstructure:"CONSUL_VERSION"`
	CONSUL_SERVICE   string `mapstructure:"CONSUL_SERVICE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("user")
	viper.SetConfigType("env")

	// read variable from env variable
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
