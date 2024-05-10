-- generateAPITokenQuery
-- $1: user_id
UPDATE
    users
SET
    api_token = random_token(),
    updated_at = now()
WHERE
    user_id = $1
    AND api_token IS NULL
    AND deleted_at IS NULL
    AND disabled_at IS NULL;
