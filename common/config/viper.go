package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	AppName                  string `mapstructure:"APP_NAME"`
	AppMode                  string `mapstructure:"APP_MODE"`
	AppPort                  string `mapstructure:"APP_PORT"`
	DBHostWrite              string `mapstructure:"DB_HOST_WRITE"`
	DBNameWrite              string `mapstructure:"DB_NAME_WRITE"`
	DBUserWrite              string `mapstructure:"DB_USER_WRITE"`
	DBPortWrite              string `mapstructure:"DB_PORT_WRITE"`
	DBPasswordWrite          string `mapstructure:"DB_PASSWORD_WRITE"`
	DBHostRead               string `mapstructure:"DB_HOST_READ"`
	DBNameRead               string `mapstructure:"DB_NAME_READ"`
	DBUserRead               string `mapstructure:"DB_USER_READ"`
	DBPortRead               string `mapstructure:"DB_PORT_READ"`
	DBPasswordRead           string `mapstructure:"DB_PASSWORD_READ"`
	JWTSecret                string `mapstructure:"JWT_SECRET"`
	ErrorLink                string `mapstructure:"ERROR_LINK_INFO"`
	TokenExpiredMinuet       int64  `mapstructure:"TOKEN_EXPIRED_MINUET"`
	RereshTokenExpiredMinuet int64  `mapstructure:"REFRESH_TOKEN_EXPIRED_MINUET"`
	CorsOrigins              string `mapstructure:"CORS_ALLOWED_ORIGINS"`
}

var envs = []string{
	"APP_NAME",
	"APP_MODE",
	"APP_PORT",
	"DB_HOST_WRITE",
	"DB_NAME_WRITE",
	"DB_USER_WRITE",
	"DB_PORT_WRITE",
	"DB_PASSWORD_WRITE",
	"DB_HOST_READ",
	"DB_NAME_READ",
	"DB_USER_READ",
	"DB_PORT_READ",
	"DB_PASSWORD_READ",
	"JWT_SECRET",
	"ERROR_LINK_INFO",
	"REFRESH_TOKEN_EXPIRED_MINUET",
	"TOKEN_EXPIRED_MINUET",
	"CORS_ALLOWED_ORIGINS",
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	return config, nil
}

func LoadConfigNostuct() (map[string]interface{}, error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	config := viper.AllSettings()
	if err := validator.New().Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}
