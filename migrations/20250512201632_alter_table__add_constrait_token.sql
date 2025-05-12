-- +goose Up
-- +goose StatementBegin
ALTER TABLE tokens 
ADD CONSTRAINT tokens_name_symbol_network_unique 
UNIQUE (name, symbol, network_id, network);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tokens 
DROP CONSTRAINT IF EXISTS tokens_name_symbol_network_unique;
-- +goose StatementEnd
