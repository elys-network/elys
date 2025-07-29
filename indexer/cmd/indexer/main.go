package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/elys-network/elys/indexer/internal/cache"
	"github.com/elys-network/elys/indexer/internal/config"
	"github.com/elys-network/elys/indexer/internal/database"
	"github.com/elys-network/elys/indexer/internal/indexer"
	"github.com/elys-network/elys/indexer/internal/websocket"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	configPath string
	logLevel   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "elys-indexer",
		Short: "Elys chain indexer for TradeSheild and Perpetual modules",
		RunE:  run,
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to configuration file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Setup logger
	logger, err := setupLogger(logLevel)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create database connection
	db, err := database.New(&cfg.Database, logger)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Create repository
	repo := database.NewRepository(db, logger)

	// Create Redis cache
	cache, err := cache.New(&cfg.Redis, logger)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	defer cache.Close()

	// Create WebSocket server
	wsServer := websocket.NewServer(&cfg.WebSocket, cache, repo, logger)
	if err := wsServer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start WebSocket server: %w", err)
	}
	defer wsServer.Stop(context.Background())

	// Create and start indexer
	indexer, err := indexer.New(cfg, repo, cache, logger)
	if err != nil {
		return fmt.Errorf("failed to create indexer: %w", err)
	}

	if err := indexer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start indexer: %w", err)
	}
	defer indexer.Stop()

	logger.Info("Indexer started successfully",
		zap.String("chain_id", cfg.Chain.ChainID),
		zap.String("rpc_endpoint", cfg.Chain.RPCEndpoint),
		zap.String("ws_addr", cfg.WebSocket.ListenAddr),
	)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down indexer...")
	return nil
}

func setupLogger(level string) (*zap.Logger, error) {
	// Parse log level
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// Create logger configuration
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Build logger
	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return logger, nil
}
