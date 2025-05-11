-- +goose Up
-- +goose StatementBegin
CREATE IF NOT EXISTS TABLE tokens (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    network_id VARCHAR(100) NOT NULL,
    network VARCHAR(100) NOT NULL
);
CREATE IF NOT EXISTS TABLE pairs (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    numerator UUID NOT NULL REFERENCES tokens(uuid) ON DELETE CASCADE,
    denominator UUID NOT NULL REFERENCES tokens(uuid) ON DELETE CASCADE,
    timeframe INTERVAL NOT NULL,
    is_active BOOL DEFAULT FALSE NOT NULL,
    UNIQUE (numerator, denominator, timeframe)
);
CREATE IF NOT EXISTS TABLE prices (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ts TIMESTAMP NOT NULL,
    price NUMERIC(20, 12) NOT NULL,
    pair_uuid UUID NOT NULL REFERENCES pairs(uuid) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS prices;
DROP TABLE IF EXISTS pairs;
DROP TABLE IF EXISTS tokens;
-- +goose StatementEnd