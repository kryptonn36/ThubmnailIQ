package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Env  string `mapstructure:"env"`
	} `mapstructure:"server"`

	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"database"`

	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	JWT struct {
		AccessSecret  string        `mapstructure:"access_secret"`
		RefreshSecret string        `mapstructure:"refresh_secret"`
		AccessTTL     time.Duration `mapstructure:"access_ttl"`
		RefreshTTL    time.Duration `mapstructure:"refresh_ttl"`
	} `mapstructure:"jwt"`

	S3 struct {
		Endpoint        string `mapstructure:"endpoint"`
		Region          string `mapstructure:"region"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		UploadBucket    string `mapstructure:"upload_bucket"`
		PublicBaseURL   string `mapstructure:"public_base_url"`
		UsePathStyle    bool   `mapstructure:"use_path_style"`
	} `mapstructure:"s3"`

	YouTube struct {
		APIKey string `mapstructure:"api_key"`
	} `mapstructure:"youtube"`

	Payment struct {
		Provider string `mapstructure:"provider"`
		Currency string `mapstructure:"currency"`
	} `mapstructure:"payment"`

	Razorpay struct {
		KeyID     string `mapstructure:"key_id"`
		KeySecret string `mapstructure:"key_secret"`
	} `mapstructure:"razorpay"`

	Stripe struct {
		SecretKey string `mapstructure:"secret_key"`
	} `mapstructure:"stripe"`

	CVService struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"cv_service"`

	Gemini struct {
		APIKey string `mapstructure:"api_key"`
		Model  string `mapstructure:"model"`
	} `mapstructure:"gemini"`
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error in loading env file: %v",err)
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	v.SetDefault("server.port", 8080)
	v.SetDefault("server.env", "development")
	v.SetDefault("database.url", "postgres://thumbnailiq:thumbnailiq@localhost:5432/thumbnailiq?sslmode=disable")
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.db", 0)
	v.SetDefault("jwt.access_secret", "dev-access-secret-change-me")
	v.SetDefault("jwt.refresh_secret", "dev-refresh-secret-change-me")
	v.SetDefault("jwt.access_ttl", "15m")
	v.SetDefault("jwt.refresh_ttl", "168h")
	v.SetDefault("s3.endpoint", "http://localhost:9000")
	v.SetDefault("s3.region", "us-east-1")
	v.SetDefault("s3.access_key_id", "minioadmin")
	v.SetDefault("s3.secret_access_key", "minioadmin")
	v.SetDefault("s3.upload_bucket", "thumbnailiq-uploads")
	v.SetDefault("s3.public_base_url", "http://localhost:9000/thumbnailiq-uploads")
	v.SetDefault("s3.use_path_style", true)
	v.SetDefault("cv_service.url", "http://localhost:8001")
	v.SetDefault("gemini.api", "geminiAPI")
	v.SetDefault("gemini.model", "gemini-2.0-flash")
	v.SetDefault("payment.provider", "razorpay")
	v.SetDefault("payment.currency", "INR")
	v.SetDefault("razorpay.key_id", "defaultID")
	v.SetDefault("razorpay.key_secret", "key_secret")
	v.SetDefault("youTube.api_key", "youtubeAPI")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}
	return &cfg, nil
}
