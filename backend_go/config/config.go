package config

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
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
	Database       struct {
		MaxOpenConns    int           `json:"max_open_conns"`
		MaxIdleConns    int           `json:"max_idle_conns"`
		ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
		ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
	} `json:"database"`
}

// readSecret reads a secret from file if available, otherwise returns env var
func readSecret(envVar, fileVar string) string {
	// Try file first (Docker secrets)
	if fileContent := os.Getenv(fileVar); fileContent != "" {
		if data, err := ioutil.ReadFile(fileContent); err == nil {
			return strings.TrimSpace(string(data))
		}
	}
	// Fall back to environment variable
	return os.Getenv(envVar)
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

	// Read JWT secret from file or environment
	cfg.JWT.Secret = readSecret("JWT_SECRET", "JWT_SECRET_FILE")
	if cfg.JWT.Secret == "" {
		panic("JWT_SECRET environment variable or file required")
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

	// Database connection pool settings
	cfg.Database.MaxOpenConns, _ = strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 25
	}
	cfg.Database.MaxIdleConns, _ = strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 5
	}
	cfg.Database.ConnMaxLifetime, _ = time.ParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME"))
	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = 5 * time.Minute
	}
	cfg.Database.ConnMaxIdleTime, _ = time.ParseDuration(os.Getenv("DB_CONN_MAX_IDLE_TIME"))
	if cfg.Database.ConnMaxIdleTime == 0 {
		cfg.Database.ConnMaxIdleTime = 1 * time.Minute
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
