package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"crypto/tls"
)

type Config struct {
	TLS struct {
		Enabled bool
		Cert    string
		Key     string
	} `json:"tls"`
	JWT struct {
		AccessTokenExpiry  time.Duration `json:"access_token_expiry"`
		RefreshTokenExpiry time.Duration `json:"refresh_token_expiry"`
		Issuer             string        `json:"issuer"`
		Audience           string        `json:"audience"`
		Secret             string        `json:"secret"`
		Expiry             time.Duration `json:"expiry"`
	} `json:"jwt"`
	LoginRateLimit int `json:"login_rate_limit"`
	RateLimit      int `json:"rate_limit"`
}

func Load() *Config {
	cfg := &Config{}

	// TLS config
	tlsEnabledStr := os.Getenv("ENABLE_TLS")
	cfg.TLS.Enabled = tlsEnabledStr == "true" || tlsEnabledStr == "1"
	flag.BoolVar(&cfg.TLS.Enabled, "tls", cfg.TLS.Enabled, "Enable TLS")

	cfg.TLS.Cert = os.Getenv("TLS_CERT_PATH")
	if cfg.TLS.Cert == "" {
		cfg.TLS.Cert = "certs/cert.pem"
	}
	flag.StringVar(&cfg.TLS.Cert, "tls-cert", cfg.TLS.Cert, "TLS certificate path")

	cfg.TLS.Key = os.Getenv("TLS_KEY_PATH")
	if cfg.TLS.Key == "" {
		cfg.TLS.Key = "certs/key.pem"
	}
	flag.StringVar(&cfg.TLS.Key, "tls-key", cfg.TLS.Key, "TLS private key path")

	// JWT config
	cfg.JWT.Issuer = os.Getenv("JWT_ISSUER")
	if cfg.JWT.Issuer == "" {
		cfg.JWT.Issuer = "rojgarsetu-backend"
	}
	cfg.JWT.Audience = os.Getenv("JWT_AUDIENCE")
	if cfg.JWT.Audience == "" {
		cfg.JWT.Audience = "rojgarsetu-api"
	}
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	if cfg.JWT.Secret == "" {
		panic("JWT_SECRET environment variable required")
	}
	if len(cfg.JWT.Secret) < 32 {
		panic("JWT_SECRET must be at least 32 characters long for security")
	}
	expiryStr := os.Getenv("JWT_EXPIRY_SECONDS")
	expiry, _ := strconv.Atoi(expiryStr)
	if expiry > 0 {
		cfg.JWT.Expiry = time.Duration(expiry) * time.Second
	} else {
		cfg.JWT.Expiry = 24 * time.Hour
	}
	// Access token expiry
	cfg.JWT.AccessTokenExpiry, _ = time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	if cfg.JWT.AccessTokenExpiry == 0 {
		cfg.JWT.AccessTokenExpiry = 15 * time.Minute
	}
	// Refresh token expiry
	cfg.JWT.RefreshTokenExpiry, _ = time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if cfg.JWT.RefreshTokenExpiry == 0 {
		cfg.JWT.RefreshTokenExpiry = 30 * 24 * time.Hour
	}

	// Rate limit
	// Rate limit
	rateStr := os.Getenv("RATE_LIMIT")
	cfg.RateLimit, _ = strconv.Atoi(rateStr)
	if cfg.RateLimit == 0 {
		cfg.RateLimit = 100
	}
	// Login rate limit
	loginRateStr := os.Getenv("LOGIN_RATE_LIMIT")
	cfg.LoginRateLimit, _ = strconv.Atoi(loginRateStr)
	if cfg.LoginRateLimit == 0 {
		cfg.LoginRateLimit = 5
	}

	flag.Parse()
	return cfg
}

// LoadTLSCert loads and returns the TLS certificate pair
func (c *Config) LoadTLSCert() tls.Certificate {
	cert, err := tls.LoadX509KeyPair(c.TLS.Cert, c.TLS.Key)
	if err != nil {
		panic("Failed to load TLS cert/key: " + err.Error())
	}
	return cert
}
