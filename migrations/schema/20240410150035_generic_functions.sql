-- +goose Up
-- +goose StatementBegin

-- create a function that prevents an operation
CREATE OR REPLACE FUNCTION operation_not_allowed()
RETURNS TRIGGER AS $$
BEGIN
  RAISE EXCEPTION '% operation on table % is not allowed: %', TG_OP, TG_TABLE_NAME, TG_NAME::TEXT;
END;
$$ LANGUAGE plpgsql;

-- create a function that generates a random token
CREATE OR REPLACE FUNCTION random_token() RETURNS text AS $$
BEGIN
    RETURN md5(random()::text || clock_timestamp()::text);
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS operation_not_allowed();
-- +goose StatementEnd
