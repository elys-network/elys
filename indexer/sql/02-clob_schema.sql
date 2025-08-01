-- CLOB Module Tables
-- Separate from Perpetual and TradeSheild modules

-- CLOB Order Types
CREATE TYPE clob_order_type AS ENUM ('LIMIT_BUY', 'LIMIT_SELL', 'MARKET_BUY', 'MARKET_SELL');
CREATE TYPE clob_order_status AS ENUM ('PENDING', 'PARTIALLY_FILLED', 'FILLED', 'CANCELLED', 'EXPIRED');
CREATE TYPE clob_position_side AS ENUM ('LONG', 'SHORT');

-- CLOB Markets
CREATE TABLE clob_markets (
    id SERIAL PRIMARY KEY,
    market_id BIGINT NOT NULL UNIQUE,
    ticker VARCHAR(50) NOT NULL,
    base_asset VARCHAR(128) NOT NULL,
    quote_asset VARCHAR(128) NOT NULL,
    tick_size NUMERIC NOT NULL,
    lot_size NUMERIC NOT NULL,
    min_order_size NUMERIC NOT NULL,
    max_order_size NUMERIC NOT NULL,
    max_leverage NUMERIC NOT NULL,
    initial_margin_fraction NUMERIC NOT NULL,
    maintenance_margin_fraction NUMERIC NOT NULL,
    funding_interval BIGINT NOT NULL,
    next_funding_time TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    block_height BIGINT NOT NULL
);

CREATE INDEX idx_clob_markets_ticker ON clob_markets (ticker);
CREATE INDEX idx_clob_markets_active ON clob_markets (is_active);

-- CLOB Orders
CREATE TABLE clob_orders (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    market_id BIGINT NOT NULL,
    counter BIGINT NOT NULL,
    owner VARCHAR(63) NOT NULL,
    sub_account_id BIGINT NOT NULL,
    order_type clob_order_type NOT NULL,
    price NUMERIC NOT NULL,
    amount NUMERIC NOT NULL,
    filled_amount NUMERIC NOT NULL DEFAULT 0,
    remaining_amount NUMERIC NOT NULL,
    status clob_order_status NOT NULL,
    time_in_force VARCHAR(20) NOT NULL DEFAULT 'GTC',
    post_only BOOLEAN NOT NULL DEFAULT FALSE,
    reduce_only BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    executed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    expires_at TIMESTAMP,
    block_height BIGINT NOT NULL,
    tx_hash VARCHAR(64),
    UNIQUE(market_id, counter)
);

CREATE INDEX idx_clob_orders_market_id ON clob_orders (market_id);
CREATE INDEX idx_clob_orders_owner ON clob_orders (owner);
CREATE INDEX idx_clob_orders_sub_account ON clob_orders (sub_account_id);
CREATE INDEX idx_clob_orders_status ON clob_orders (status);
CREATE INDEX idx_clob_orders_type ON clob_orders (order_type);
CREATE INDEX idx_clob_orders_created_at ON clob_orders (created_at DESC);
CREATE INDEX idx_clob_orders_market_status ON clob_orders (market_id, status) WHERE status IN ('PENDING', 'PARTIALLY_FILLED');
CREATE INDEX idx_clob_orders_market_type_price ON clob_orders (market_id, order_type, price) WHERE status IN ('PENDING', 'PARTIALLY_FILLED');

-- CLOB Trades
CREATE TABLE clob_trades (
    id SERIAL PRIMARY KEY,
    trade_id BIGSERIAL UNIQUE,
    market_id BIGINT NOT NULL,
    buyer VARCHAR(63) NOT NULL,
    buyer_sub_account_id BIGINT NOT NULL,
    seller VARCHAR(63) NOT NULL,
    seller_sub_account_id BIGINT NOT NULL,
    buyer_order_id BIGINT NOT NULL,
    seller_order_id BIGINT NOT NULL,
    price NUMERIC NOT NULL,
    quantity NUMERIC NOT NULL,
    trade_value NUMERIC NOT NULL,
    buyer_fee NUMERIC NOT NULL,
    seller_fee NUMERIC NOT NULL,
    is_buyer_taker BOOLEAN NOT NULL,
    is_buyer_liquidation BOOLEAN NOT NULL DEFAULT FALSE,
    is_seller_liquidation BOOLEAN NOT NULL DEFAULT FALSE,
    executed_at TIMESTAMP NOT NULL,
    block_height BIGINT NOT NULL,
    tx_hash VARCHAR(64) NOT NULL
);

CREATE INDEX idx_clob_trades_market_id ON clob_trades (market_id);
CREATE INDEX idx_clob_trades_buyer ON clob_trades (buyer);
CREATE INDEX idx_clob_trades_seller ON clob_trades (seller);
CREATE INDEX idx_clob_trades_executed_at ON clob_trades (executed_at DESC);
CREATE INDEX idx_clob_trades_buyer_order ON clob_trades (buyer_order_id);
CREATE INDEX idx_clob_trades_seller_order ON clob_trades (seller_order_id);

-- CLOB Positions (Perpetual Positions)
CREATE TABLE clob_positions (
    id SERIAL PRIMARY KEY,
    position_id BIGSERIAL UNIQUE,
    market_id BIGINT NOT NULL,
    owner VARCHAR(63) NOT NULL,
    sub_account_id BIGINT NOT NULL,
    side clob_position_side NOT NULL,
    size NUMERIC NOT NULL,
    notional NUMERIC NOT NULL,
    entry_price NUMERIC NOT NULL,
    mark_price NUMERIC NOT NULL,
    liquidation_price NUMERIC NOT NULL,
    margin NUMERIC NOT NULL,
    margin_ratio NUMERIC NOT NULL,
    unrealized_pnl NUMERIC NOT NULL,
    realized_pnl NUMERIC NOT NULL DEFAULT 0,
    cumulative_funding NUMERIC NOT NULL DEFAULT 0,
    last_funding_payment NUMERIC NOT NULL DEFAULT 0,
    last_funding_time TIMESTAMP,
    opened_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    closed_at TIMESTAMP,
    close_price NUMERIC,
    is_liquidated BOOLEAN NOT NULL DEFAULT FALSE,
    block_height BIGINT NOT NULL,
    tx_hash VARCHAR(64)
);

CREATE INDEX idx_clob_positions_market_id ON clob_positions (market_id);
CREATE INDEX idx_clob_positions_owner ON clob_positions (owner);
CREATE INDEX idx_clob_positions_sub_account ON clob_positions (sub_account_id);
CREATE INDEX idx_clob_positions_opened_at ON clob_positions (opened_at DESC);
CREATE INDEX idx_clob_positions_active ON clob_positions (closed_at) WHERE closed_at IS NULL;
CREATE INDEX idx_clob_positions_liquidation ON clob_positions (market_id, liquidation_price) WHERE closed_at IS NULL;

-- CLOB Order Book Snapshots
CREATE TABLE clob_order_book_snapshots (
    id SERIAL PRIMARY KEY,
    market_id BIGINT NOT NULL,
    bids JSONB NOT NULL, -- Array of {price, quantity, orders}
    asks JSONB NOT NULL, -- Array of {price, quantity, orders}
    best_bid NUMERIC,
    best_ask NUMERIC,
    mid_price NUMERIC,
    spread NUMERIC,
    total_bid_volume NUMERIC NOT NULL,
    total_ask_volume NUMERIC NOT NULL,
    snapshot_time TIMESTAMP NOT NULL,
    block_height BIGINT NOT NULL
);

CREATE INDEX idx_clob_snapshots_market_id ON clob_order_book_snapshots (market_id);
CREATE INDEX idx_clob_snapshots_time ON clob_order_book_snapshots (snapshot_time DESC);
CREATE INDEX idx_clob_snapshots_market_time ON clob_order_book_snapshots (market_id, snapshot_time DESC);

-- CLOB Funding Rates
CREATE TABLE clob_funding_rates (
    id SERIAL PRIMARY KEY,
    market_id BIGINT NOT NULL,
    funding_rate NUMERIC NOT NULL,
    premium_rate NUMERIC NOT NULL,
    mark_price NUMERIC NOT NULL,
    index_price NUMERIC NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    next_funding_time TIMESTAMP NOT NULL,
    block_height BIGINT NOT NULL
);

CREATE INDEX idx_clob_funding_market_id ON clob_funding_rates (market_id);
CREATE INDEX idx_clob_funding_timestamp ON clob_funding_rates (timestamp DESC);
CREATE INDEX idx_clob_funding_market_time ON clob_funding_rates (market_id, timestamp DESC);

-- CLOB Liquidations
CREATE TABLE clob_liquidations (
    id SERIAL PRIMARY KEY,
    liquidation_id BIGSERIAL UNIQUE,
    market_id BIGINT NOT NULL,
    position_id BIGINT NOT NULL,
    owner VARCHAR(63) NOT NULL,
    sub_account_id BIGINT NOT NULL,
    liquidator VARCHAR(63),
    side clob_position_side NOT NULL,
    size NUMERIC NOT NULL,
    price NUMERIC NOT NULL,
    liquidation_fee NUMERIC NOT NULL,
    insurance_fund_contribution NUMERIC NOT NULL,
    is_adl BOOLEAN NOT NULL DEFAULT FALSE,
    liquidated_at TIMESTAMP NOT NULL,
    block_height BIGINT NOT NULL,
    tx_hash VARCHAR(64) NOT NULL
);

CREATE INDEX idx_clob_liquidations_market_id ON clob_liquidations (market_id);
CREATE INDEX idx_clob_liquidations_owner ON clob_liquidations (owner);
CREATE INDEX idx_clob_liquidations_position ON clob_liquidations (position_id);
CREATE INDEX idx_clob_liquidations_time ON clob_liquidations (liquidated_at DESC);

-- CLOB Market Statistics (for quick access)
CREATE TABLE clob_market_stats (
    id SERIAL PRIMARY KEY,
    market_id BIGINT NOT NULL UNIQUE,
    volume_24h NUMERIC NOT NULL DEFAULT 0,
    trades_24h BIGINT NOT NULL DEFAULT 0,
    high_24h NUMERIC,
    low_24h NUMERIC,
    open_interest NUMERIC NOT NULL DEFAULT 0,
    open_interest_notional NUMERIC NOT NULL DEFAULT 0,
    last_price NUMERIC,
    last_trade_time TIMESTAMP,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_clob_stats_market_id ON clob_market_stats (market_id);

-- Function to update market statistics
CREATE OR REPLACE FUNCTION update_clob_market_stats()
RETURNS TRIGGER AS $$
BEGIN
    -- Update market stats when a trade occurs
    UPDATE clob_market_stats
    SET 
        volume_24h = (
            SELECT COALESCE(SUM(trade_value), 0)
            FROM clob_trades
            WHERE market_id = NEW.market_id
            AND executed_at > NOW() - INTERVAL '24 hours'
        ),
        trades_24h = (
            SELECT COUNT(*)
            FROM clob_trades
            WHERE market_id = NEW.market_id
            AND executed_at > NOW() - INTERVAL '24 hours'
        ),
        last_price = NEW.price,
        last_trade_time = NEW.executed_at,
        updated_at = NOW()
    WHERE market_id = NEW.market_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for updating market stats on new trades
CREATE TRIGGER update_market_stats_on_trade
AFTER INSERT ON clob_trades
FOR EACH ROW
EXECUTE FUNCTION update_clob_market_stats();

-- Views for easier querying

-- Active orders view
CREATE VIEW clob_active_orders AS
SELECT * FROM clob_orders
WHERE status IN ('PENDING', 'PARTIALLY_FILLED')
ORDER BY created_at DESC;

-- Open positions view
CREATE VIEW clob_open_positions AS
SELECT * FROM clob_positions
WHERE closed_at IS NULL
ORDER BY opened_at DESC;

-- Recent trades view
CREATE VIEW clob_recent_trades AS
SELECT * FROM clob_trades
WHERE executed_at > NOW() - INTERVAL '24 hours'
ORDER BY executed_at DESC;