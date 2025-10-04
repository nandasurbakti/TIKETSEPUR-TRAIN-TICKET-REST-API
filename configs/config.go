package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	URL 		  string
	Host          string
	Port          string
	Password      string
	DB            int
	SessionExpiry time.Duration
}

type RabbitMQConfig struct {
	URL 	  string	
	URLLokal  string 
	QueueName string 
}

type JWTConfig struct {
	Secret string
	Exp time.Duration
}

type PaymentConfig struct {
	MidtransServerKey string
	MidtransClientKey string
	MidtransEnv       string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	JWT      JWTConfig
	Payment  PaymentConfig
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	sessionExpiryStr := viper.GetString("redis.session_expiry")
    duration, err := time.ParseDuration(sessionExpiryStr)
    if err != nil {
        return nil, fmt.Errorf("invalid session_expiry: %w", err)
    }
    config.Redis.SessionExpiry = duration

	ConnMaxLifetimeStr := viper.GetString("database.conn_max_lifetime")
    duration1, err := time.ParseDuration(ConnMaxLifetimeStr)
    if err != nil {
        return nil, fmt.Errorf("invalid session_expiry: %w", err)
    }
    config.Database.ConnMaxLifetime = duration1

	config.Payment.MidtransClientKey = viper.GetString("MIDTRANS_CLIENT_KEY")
	config.Payment.MidtransServerKey = viper.GetString("MIDTRANS_SERVER_KEY")

	return &config, nil
}