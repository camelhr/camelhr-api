-- disableUserQuery
-- $1: user_id
-- $2: comment
UPDATE
    users
SET
    disabled_at = now(),
    comment = $2
WHERE
    user_id = $1
    AND deleted_at IS NULL
    AND disabled_at IS NULL
    AND NOT is_owner;
