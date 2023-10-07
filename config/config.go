package config

import (
	"os"

	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Kafka struct {
		Acks                 string `yaml:"acks"`
		ProducerID           string `yaml:"producerId"`
		BootstrapServersHost string `yaml:"bootstrapServersHost"`
		BootstrapServersPort string `yaml:"bootstrapServersPort"`
		AutoOffsetReset      string `yaml:"autoOffsetReset"`
		GroupID              string `yaml:"groupId"`
	} `yaml:"Kafka"`

	Encryption struct {
		Algorithm     string `yaml:"algorithm"`
		EncryptionKey string `yaml:"encryptionKey"`
	} `yaml:"Encryption"`

	ProducerServer struct {
		BindAddr     string `yaml:"bindAddr"`
		ReadTimeout  string `yaml:"readTimeout"`
		WriteTimeout string `yaml:"writeTimeout"`
		IdleTimeout  string `yaml:"idleTimeout"`
	} `yaml:"ProducerServer"`

	ClientServer struct {
		BindAddr     string `yaml:"bindAddr"`
		ReadTimeout  string `yaml:"readTimeout"`
		WriteTimeout string `yaml:"writeTimeout"`
		IdleTimeout  string `yaml:"idleTimeout"`
	} `yaml:"ClientServer"`

	PostgresDB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbname"`
		SSLmode  string `yaml:"sslmode"`
		MaxConns string `yaml:"maxconns"`
	} `yaml:"PostgresDB"`

	Authentication struct {
		JWTSecretKey string `yaml:"jwtSecretKey"`
	} `yaml:"Authentication"`

	Session struct {
		SessionKey string `yaml:"sessionKey"`
	} `yaml:"Session"`

	Redis struct {
		Password string `yaml:"password"`
	} `yaml:"Redis"`
}

func ConfigureViper() *viper.Viper {
	logger := logger.GetLogger()

	v := viper.New()
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		logger.Error("Unable to read the configuration file.", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("Config loaded successfully.")

	return v
}

func ParseConfig(v *viper.Viper) *Config {
	logger := logger.GetLogger()

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		logger.Error("Unable to parse the configuration file.")
	}
	logger.Info("Configuration file parsed successfully.")

	return cfg
}
