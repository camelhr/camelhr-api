-- resetPasswordQuery
-- $1 - user_id
-- $2 - password_hash
UPDATE users
SET
    password_hash = $2,
    updated_at = now()
WHERE
    user_id = $1
    AND deleted_at IS NULL
    AND disabled_at IS NULL;
