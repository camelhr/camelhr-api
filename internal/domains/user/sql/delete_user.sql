-- deleteUserQuery
-- $1: user_id
UPDATE
    users
SET
    deleted_at = now()
WHERE
    user_id = $1
    AND deleted_at IS NULL;
