package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Chain     ChainConfig     `mapstructure:"chain"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
	Indexer   IndexerConfig   `mapstructure:"indexer"`
}

type ChainConfig struct {
	RPCEndpoint  string        `mapstructure:"rpc_endpoint"`
	GRPCEndpoint string        `mapstructure:"grpc_endpoint"`
	ChainID      string        `mapstructure:"chain_id"`
	StartHeight  int64         `mapstructure:"start_height"`
	PollInterval time.Duration `mapstructure:"poll_interval"`
	MaxRetries   int           `mapstructure:"max_retries"`
	RetryDelay   time.Duration `mapstructure:"retry_delay"`
	BatchSize    int           `mapstructure:"batch_size"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxConnAge   time.Duration `mapstructure:"max_conn_age"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type WebSocketConfig struct {
	ListenAddr      string        `mapstructure:"listen_addr"`
	ReadBufferSize  int           `mapstructure:"read_buffer_size"`
	WriteBufferSize int           `mapstructure:"write_buffer_size"`
	MaxMessageSize  int64         `mapstructure:"max_message_size"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	PongTimeout     time.Duration `mapstructure:"pong_timeout"`
	PingInterval    time.Duration `mapstructure:"ping_interval"`
	MaxConnections  int           `mapstructure:"max_connections"`
}

type IndexerConfig struct {
	WorkerCount             int           `mapstructure:"worker_count"`
	EventBufferSize         int           `mapstructure:"event_buffer_size"`
	OrderBookUpdateInterval time.Duration `mapstructure:"order_book_update_interval"`
	MetricsEnabled          bool          `mapstructure:"metrics_enabled"`
	MetricsPort             int           `mapstructure:"metrics_port"`
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Chain defaults
	viper.SetDefault("chain.rpc_endpoint", "http://localhost:26657")
	viper.SetDefault("chain.grpc_endpoint", "localhost:9090")
	viper.SetDefault("chain.chain_id", "elys-1")
	viper.SetDefault("chain.start_height", 0)
	viper.SetDefault("chain.poll_interval", "1s")
	viper.SetDefault("chain.max_retries", 5)
	viper.SetDefault("chain.retry_delay", "5s")
	viper.SetDefault("chain.batch_size", 100)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "indexer")
	viper.SetDefault("database.database", "elys_indexer")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.conn_max_lifetime", "1h")

	// Redis defaults
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.dial_timeout", "5s")
	viper.SetDefault("redis.read_timeout", "3s")
	viper.SetDefault("redis.write_timeout", "3s")
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)
	viper.SetDefault("redis.max_conn_age", "0")
	viper.SetDefault("redis.pool_timeout", "4s")
	viper.SetDefault("redis.idle_timeout", "5m")

	// WebSocket defaults
	viper.SetDefault("websocket.listen_addr", ":8080")
	viper.SetDefault("websocket.read_buffer_size", 1024)
	viper.SetDefault("websocket.write_buffer_size", 1024)
	viper.SetDefault("websocket.max_message_size", 512000)
	viper.SetDefault("websocket.write_timeout", "10s")
	viper.SetDefault("websocket.pong_timeout", "60s")
	viper.SetDefault("websocket.ping_interval", "54s")
	viper.SetDefault("websocket.max_connections", 1000)

	// Indexer defaults
	viper.SetDefault("indexer.worker_count", 4)
	viper.SetDefault("indexer.event_buffer_size", 1000)
	viper.SetDefault("indexer.order_book_update_interval", "1s")
	viper.SetDefault("indexer.metrics_enabled", true)
	viper.SetDefault("indexer.metrics_port", 9090)
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}
