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

	// AdminJWT is a deliberately separate secret/TTL pair from JWT above, so
	// admin session tokens can never be confused with — or forged from —
	// customer session tokens, even though both are signed with the same
	// pkg/jwt.Service implementation.
	AdminJWT struct {
		AccessSecret  string        `mapstructure:"access_secret"`
		RefreshSecret string        `mapstructure:"refresh_secret"`
		AccessTTL     time.Duration `mapstructure:"access_ttl"`
		RefreshTTL    time.Duration `mapstructure:"refresh_ttl"`
	} `mapstructure:"admin_jwt"`

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

	// CORS.AllowedOrigins is a comma-separated allowlist of browser origins
	// permitted to make cross-origin calls (e.g. the web app). We never reflect
	// an arbitrary Origin back; only entries on this list are echoed.
	CORS struct {
		AllowedOrigins string `mapstructure:"allowed_origins"`
	} `mapstructure:"cors"`

	YouTube struct {
		APIKey string `mapstructure:"api_key"`
	} `mapstructure:"youtube"`

	// SMTP delivers transactional email (e.g. email-verification codes). Leave
	// host/from empty to run without email in local dev — codes are still
	// generated, just not sent. For Gmail use host smtp.gmail.com, port 587,
	// and a 16-char App Password (not your normal password).
	SMTP struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		From     string `mapstructure:"from"`
		FromName string `mapstructure:"from_name"`
	} `mapstructure:"smtp"`

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
		return nil, fmt.Errorf("error in loading env file: %v", err)
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
	v.BindEnv("jwt.access_secret", "JWT_ACCESS_SECRET")
	v.BindEnv("jwt.refresh_secret", "JWT_REFRESH_SECRET")
	v.SetDefault("jwt.access_ttl", "15m")
	v.SetDefault("jwt.refresh_ttl", "168h")
	v.BindEnv("admin_jwt.access_secret", "ADMIN_JWT_ACCESS_SECRET")
	v.BindEnv("admin_jwt.refresh_secret", "ADMIN_JWT_REFRESH_SECRET")
	v.SetDefault("admin_jwt.access_ttl", "15m")
	v.SetDefault("admin_jwt.refresh_ttl", "168h")
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
	v.SetDefault("cors.allowed_origins", "http://localhost:3000")
	v.SetDefault("smtp.port", 587)
	v.SetDefault("smtp.from_name", "ThumbnailIQ")
	v.SetDefault("cv_service.url", "http://localhost:8001")
	v.BindEnv("gemini.api_key", "GEMINI_API_KEY")
	v.SetDefault("gemini.model", "gemini-2.0-flash")
	v.SetDefault("payment.provider", "razorpay")
	v.SetDefault("payment.currency", "USD")
	v.SetDefault("razorpay.key_id", "defaultID")
	v.BindEnv("razorpay.key_secret", "RAZORPAY_KEY_SECRET")
	v.SetDefault("youtube.api_key", "youtubeAPI")

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
	if err := cfg.validateSecrets(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// devDefaultSecrets are the placeholder signing secrets shipped as viper
// defaults for local development. They are public (they live in source), so
// running with any of them signs tokens anyone can forge. In production the
// app must refuse to boot until real secrets are injected via the environment.
var devDefaultSecrets = map[string]bool{
	"dev-access-secret-change-me":        true,
	"dev-refresh-secret-change-me":       true,
	"dev-admin-access-secret-change-me":  true,
	"dev-admin-refresh-secret-change-me": true,
}

// validateSecrets fails fast when a signing secret is empty or still set to a
// known dev default while running in production, so a misconfigured deploy can
// never silently come up with forgeable JWTs.
func (c *Config) validateSecrets() error {
	if !strings.EqualFold(c.Server.Env, "production") {
		return nil
	}
	secrets := map[string]string{
		"JWT_ACCESS_SECRET":        c.JWT.AccessSecret,
		"JWT_REFRESH_SECRET":       c.JWT.RefreshSecret,
		"ADMIN_JWT_ACCESS_SECRET":  c.AdminJWT.AccessSecret,
		"ADMIN_JWT_REFRESH_SECRET": c.AdminJWT.RefreshSecret,
	}
	for name, val := range secrets {
		if val == "" {
			return fmt.Errorf("%s must be set in production", name)
		}
		if devDefaultSecrets[val] {
			return fmt.Errorf("%s is still set to its insecure development default; set a strong random secret in production", name)
		}
	}
	return nil
}
