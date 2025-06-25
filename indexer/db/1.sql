CREATE SCHEMA IF NOT EXISTS chain AUTHORIZATION elys;
CREATE SCHEMA IF NOT EXISTS tradeshield AUTHORIZATION elys;

ALTER ROLE elys SET search_path TO chain;
ALTER ROLE elys SET search_path TO tradeshield;

CREATE TABLE IF NOT EXISTS chain.block
(
    id                VARCHAR(255) NOT NULL,
    last_block_height BIGINT       NOT NULL    DEFAULT 0,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tradeshield.perpetual_orders
(
    owner_address     VARCHAR(255) NOT NULL,
    pool_id           BIGINT       NOT NULL,
    order_id          BIGINT       NOT NULL,
    order_type        SMALLINT     NOT NULL,
    is_long           bool         NOT NULL,
    collateral_amount NUMERIC      NOT NULL,
    collateral_denom  VARCHAR(128) NOT NULL,
    price             NUMERIC      NOT NULL,
    take_profit_price NUMERIC      NOT NULL,
    stop_loss_price   NUMERIC      NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (owner_address, pool_id, order_id)
);

CREATE INDEX tradeshield_perpetual_order_book ON tradeshield.perpetual_orders (pool_id, order_type, is_long, price);

CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_ts_perpetual_orders_updated_at
    BEFORE UPDATE
    ON tradeshield.perpetual_orders
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_chain_block_updated_at
    BEFORE UPDATE
    ON chain.block
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

GRANT ALL ON chain.block TO elys;
GRANT ALL ON tradeshield.perpetual_orders TO elys;
