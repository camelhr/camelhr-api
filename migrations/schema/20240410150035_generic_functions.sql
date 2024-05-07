-- +goose Up
-- +goose StatementBegin
-- create a generic function that prevents an operation
CREATE OR REPLACE FUNCTION operation_not_allowed()
RETURNS TRIGGER AS $$
BEGIN
  RAISE EXCEPTION '% operation on table % is not allowed: %', TG_OP, TG_TABLE_NAME, TG_NAME::TEXT;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS operation_not_allowed();
-- +goose StatementEnd
