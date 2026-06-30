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
		UsePathStyle    bool   `mapstructure:"use_path_style"`
		// PublicRead controls whether the app grants the bucket a public-read
		// policy itself. Leave true for local/MinIO dev where nothing else
		// fronts the bucket. Set to false once a real CloudFront distribution
		// with Origin Access Control is provisioned in front of the bucket in
		// AWS — at that point the bucket must stay private and only
		// CloudFront's OAC principal should be granted access (configured in
		// AWS, not by this app).
		PublicRead bool `mapstructure:"public_read"`
	} `mapstructure:"s3"`

	// CDN holds the public-facing domain (e.g. a CloudFront distribution)
	// that every file URL in the app is generated against. The S3 bucket
	// above is only ever a storage origin — see internal/infra/cdn.
	CDN struct {
		Domain string `mapstructure:"domain"`
	} `mapstructure:"cdn"`

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
	// No default here on purpose: an empty endpoint means "use AWS's real
	// endpoints" (production). Viper's AutomaticEnv treats an env var set
	// to "" identically to one that's unset, so any non-empty default here
	// would make it impossible to ever clear S3_ENDPOINT via .env once set.
	// Local/MinIO dev gets its endpoint from .env.example's explicit
	// S3_ENDPOINT=http://localhost:9000, not from a code-level fallback.
	v.SetDefault("s3.region", "us-east-1")
	v.SetDefault("s3.access_key_id", "minioadmin")
	v.SetDefault("s3.secret_access_key", "minioadmin")
	v.SetDefault("s3.upload_bucket", "thumbnailiq-uploads")
	v.SetDefault("s3.use_path_style", true)
	v.SetDefault("s3.public_read", true)
	v.SetDefault("cdn.domain", "http://localhost:9000/thumbnailiq-uploads")
	v.SetDefault("cv_service.url", "http://localhost:8001")
	v.SetDefault("gemini.api_key", "")
	v.SetDefault("gemini.model", "gemini-2.0-flash")
	v.SetDefault("payment.provider", "razorpay")
	v.SetDefault("payment.currency", "USD")
	v.SetDefault("razorpay.key_id", "defaultID")
	v.SetDefault("razorpay.key_secret", "key_secret")
	v.SetDefault("youTube.api_key", "youtubeAPI")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	_ = v.BindEnv("cdn.domain", "CLOUDFRONT_DOMAIN")

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
