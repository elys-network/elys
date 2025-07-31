-- Database schema for Elys Chain Indexer

-- Create database if not exists
-- CREATE DATABASE elys_indexer;

-- Extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enum types
CREATE TYPE order_type AS ENUM ('LIMIT_BUY', 'LIMIT_SELL', 'MARKET_BUY', 'STOP_LOSS');
CREATE TYPE order_status AS ENUM ('PENDING', 'EXECUTED', 'CANCELLED', 'CLOSED');
CREATE TYPE position_type AS ENUM ('LONG', 'SHORT');
CREATE TYPE perpetual_order_type AS ENUM ('LIMIT_OPEN', 'LIMIT_CLOSE', 'MARKET', 'STOP_LOSS', 'TAKE_PROFIT');

-- TradeSheild Spot Orders
CREATE TABLE spot_orders (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE,
    order_type order_type NOT NULL,
    owner_address VARCHAR(63) NOT NULL,
    order_target_denom VARCHAR(128) NOT NULL,
    order_price JSONB NOT NULL,
    order_amount NUMERIC NOT NULL,
    status order_status NOT NULL,
    created_at TIMESTAMP NOT NULL,
    executed_at TIMESTAMP,
    closed_at TIMESTAMP,
    block_height BIGINT NOT NULL,
    tx_hash VARCHAR(64)
);

CREATE INDEX idx_spot_orders_owner ON spot_orders (owner_address);
CREATE INDEX idx_spot_orders_status ON spot_orders (status);
CREATE INDEX idx_spot_orders_created_at ON spot_orders (created_at DESC);
CREATE INDEX idx_spot_orders_target_denom ON spot_orders (order_target_denom);

-- TradeSheild Perpetual Orders
CREATE TABLE perpetual_orders (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE,
    owner_address VARCHAR(63) NOT NULL,
    perpetual_order_type perpetual_order_type NOT NULL,
    position position_type NOT NULL,
    trigger_price JSONB NOT NULL,
    collateral NUMERIC NOT NULL,
    leverage NUMERIC NOT NULL,
    take_profit_price NUMERIC NOT NULL,
    stop_loss_price NUMERIC NOT NULL,
    position_id BIGINT,
    status order_status NOT NULL,
    created_at TIMESTAMP NOT NULL,
    executed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    block_height BIGINT NOT NULL,
    tx_hash VARCHAR(64)
);

CREATE INDEX idx_perpetual_orders_owner ON perpetual_orders (owner_address);
CREATE INDEX idx_perpetual_orders_status ON perpetual_orders (status);
CREATE INDEX idx_perpetual_orders_position_id ON perpetual_orders (position_id);
CREATE INDEX idx_perpetual_orders_created_at ON perpetual_orders (created_at DESC);

-- Perpetual Module Positions (MTP - Margin Trading Positions)
CREATE TABLE perpetual_positions (
    id SERIAL PRIMARY KEY,
    mtp_id BIGINT NOT NULL UNIQUE,
    owner_address VARCHAR(63) NOT NULL,
    amm_pool_id BIGINT NOT NULL,
    position position_type NOT NULL,
    collateral_asset VARCHAR(128) NOT NULL,
    collateral NUMERIC NOT NULL,
    liabilities NUMERIC NOT NULL,
    custody NUMERIC NOT NULL,
    mtp_health NUMERIC NOT NULL,
    open_price NUMERIC NOT NULL,
    closing_price NUMERIC,
    net_pnl NUMERIC,
    stop_loss_price NUMERIC NOT NULL,
    take_profit_price NUMERIC NOT NULL,
    opened_at TIMESTAMP NOT NULL,
    closed_at TIMESTAMP,
    closed_by VARCHAR(63),
    close_trigger VARCHAR(50),
    block_height BIGINT NOT NULL,
    tx_hash VARCHAR(64)
);

CREATE INDEX idx_perpetual_positions_owner ON perpetual_positions (owner_address);
CREATE INDEX idx_perpetual_positions_pool ON perpetual_positions (amm_pool_id);
CREATE INDEX idx_perpetual_positions_opened_at ON perpetual_positions (opened_at DESC);
CREATE INDEX idx_perpetual_positions_status ON perpetual_positions (closed_at);

-- Trade History (for both spot and perpetual)
CREATE TABLE trades (
    id SERIAL PRIMARY KEY,
    trade_type VARCHAR(20) NOT NULL,
    reference_id BIGINT NOT NULL,
    owner_address VARCHAR(63) NOT NULL,
    asset VARCHAR(128) NOT NULL,
    amount NUMERIC NOT NULL,
    price NUMERIC NOT NULL,
    fees JSONB NOT NULL,
    executed_at TIMESTAMP NOT NULL,
    block_height BIGINT NOT NULL,
    tx_hash VARCHAR(64) NOT NULL,
    event_type VARCHAR(100) NOT NULL
);

CREATE INDEX idx_trades_owner ON trades (owner_address);
CREATE INDEX idx_trades_asset ON trades (asset);
CREATE INDEX idx_trades_executed_at ON trades (executed_at DESC);
CREATE INDEX idx_trades_reference ON trades (trade_type, reference_id);

-- Order Book Snapshots
CREATE TABLE order_book_snapshots (
    id SERIAL PRIMARY KEY,
    asset_pair VARCHAR(255) NOT NULL,
    bids JSONB NOT NULL,
    asks JSONB NOT NULL,
    snapshot_time TIMESTAMP NOT NULL,
    block_height BIGINT NOT NULL
);

CREATE INDEX idx_order_book_asset_pair ON order_book_snapshots (asset_pair);
CREATE INDEX idx_order_book_time ON order_book_snapshots (snapshot_time DESC);

-- WebSocket Subscriptions
CREATE TABLE websocket_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id VARCHAR(255) NOT NULL,
    subscription_type VARCHAR(50) NOT NULL,
    filters JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_ping TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ws_client ON websocket_subscriptions (client_id);
CREATE INDEX idx_ws_last_ping ON websocket_subscriptions (last_ping);

-- Indexer State
CREATE TABLE indexer_state (
    id SERIAL PRIMARY KEY,
    last_processed_height BIGINT NOT NULL,
    last_processed_time TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Initialize indexer state
INSERT INTO indexer_state (last_processed_height, last_processed_time, updated_at) 
VALUES (0, NOW(), NOW()) 
ON CONFLICT DO NOTHING;

-- Include CLOB schema
\i clob_schema.sql