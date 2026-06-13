package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Auth      AuthConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
	Logging   LoggingConfig
	Agent     AgentConfig
	Monitoring MonitoringConfig
	Alerting  AlertingConfig
	Encryption EncryptionConfig
	SMTP      SMTPConfig
	Swagger   SwaggerConfig
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
	BaseURL     string `mapstructure:"base_url"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	MaxBodyMB       int           `mapstructure:"max_request_body_mb"`
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type DatabaseConfig struct {
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	Name                string        `mapstructure:"name"`
	User                string        `mapstructure:"user"`
	Password            string        `mapstructure:"password"`
	SSLMode             string        `mapstructure:"ssl_mode"`
	MaxOpenConns        int           `mapstructure:"max_open_conns"`
	MaxIdleConns        int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime     time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime     time.Duration `mapstructure:"conn_max_idle_time"`
	LogSlowQueries      bool          `mapstructure:"log_slow_queries"`
	SlowQueryThreshold  time.Duration `mapstructure:"slow_query_threshold"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		d.Host, d.Port, d.Name, d.User, d.Password, d.SSLMode,
	)
}

type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	ConnTimeout  time.Duration `mapstructure:"conn_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type AuthConfig struct {
	JWTAccessSecret  string        `mapstructure:"jwt_access_secret"`
	JWTRefreshSecret string        `mapstructure:"jwt_refresh_secret"`
	AccessTokenTTL   time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL  time.Duration `mapstructure:"refresh_token_ttl"`
	MaxFailedLogins  int           `mapstructure:"max_failed_logins"`
	LockoutDuration  time.Duration `mapstructure:"lockout_duration"`
	BcryptCost       int           `mapstructure:"bcrypt_cost"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	ExposeHeaders  []string `mapstructure:"expose_headers"`
	MaxAge         int      `mapstructure:"max_age"`
}

type RateLimitConfig struct {
	Enabled                bool `mapstructure:"enabled"`
	RequestsPerMinute      int  `mapstructure:"requests_per_minute"`
	Burst                  int  `mapstructure:"burst"`
	AuthRequestsPerMinute  int  `mapstructure:"auth_requests_per_minute"`
}

type LoggingConfig struct {
	Level       string `mapstructure:"level"`
	Format      string `mapstructure:"format"`
	Output      string `mapstructure:"output"`
	FilePath    string `mapstructure:"file_path"`
	MaxSizeMB   int    `mapstructure:"max_size_mb"`
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxAgeDays  int    `mapstructure:"max_age_days"`
}

type AgentConfig struct {
	HeartbeatInterval  time.Duration `mapstructure:"heartbeat_interval"`
	OfflineThreshold   time.Duration `mapstructure:"offline_threshold"`
	TaskExpiry         time.Duration `mapstructure:"task_expiry"`
	GRPCPort           int           `mapstructure:"grpc_port"`
	GRPCMaxRecvSizeMB  int           `mapstructure:"grpc_max_recv_size_mb"`
	GRPCMaxSendSizeMB  int           `mapstructure:"grpc_max_send_size_mb"`
}

type MonitoringConfig struct {
	DefaultCheckInterval time.Duration `mapstructure:"default_check_interval"`
	DefaultRetryInterval time.Duration `mapstructure:"default_retry_interval"`
	DefaultMaxRetries    int           `mapstructure:"default_max_retries"`
	MetricBatchSize      int           `mapstructure:"metric_batch_size"`
	MetricFlushInterval  time.Duration `mapstructure:"metric_flush_interval"`
}

type AlertingConfig struct {
	EvaluationInterval  time.Duration `mapstructure:"evaluation_interval"`
	NotificationWorkers int           `mapstructure:"notification_workers"`
	MaxRetryAttempts    int           `mapstructure:"max_retry_attempts"`
	RetryBackoffBase    time.Duration `mapstructure:"retry_backoff_base"`
}

type EncryptionConfig struct {
	SecretKey string `mapstructure:"secret_key"`
}

type SMTPConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	FromAddress string `mapstructure:"from_address"`
	FromName    string `mapstructure:"from_name"`
	UseTLS      bool   `mapstructure:"use_tls"`
}

type SwaggerConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Host        string `mapstructure:"host"`
	BasePath    string `mapstructure:"base_path"`
	Title       string `mapstructure:"title"`
	Description string `mapstructure:"description"`
	Version     string `mapstructure:"version"`
}

// Load reads configuration from file and environment variables.
// Environment variables take precedence: INFRACORE_SERVER_PORT, etc.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	// Environment variable binding: INFRACORE_DATABASE_HOST → database.host
	v.SetEnvPrefix("INFRACORE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}
