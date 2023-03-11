package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB_SOURCE string `mapstructure:"DB_SOURCE"` //数据库Mysql连接来源

	REDIS_SOURCE   string `mapstructure:"REDIS_SOURCE"`   //redis连接端口
	REDIS_PASSWORD string `mapstructure:"REDIS_PASSWORD"` //redis连接密码

	ACCESS_KEYID     string `mapstructure:"ACCESS_KEYID"` //阿里云短信服务密钥
	ACCESS_KEYSECRET string `mapstructure:"ACCESS_KEYSECRET"`

	QN_ACCESS_KEY string `mapstructure:"QN_ACCESS_KEY"` //七牛云对象存储配置Kodo
	QN_SECRET_KEY string `mapstructure:"QN_SECRET_KEY"`
	QN_BUCKET     string `mapstructure:"QN_BUCKET"`

	Consul_Address string `mapstructure:"CONSUL_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("web")
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
