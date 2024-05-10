-- getUserByAPITokenQuery
-- $1: api_token
SELECT
    user_id,
    organization_id,
    email,
    password_hash,
    api_token,
    is_owner,
    disabled_at,
    comment,
    created_at,
    updated_at,
    deleted_at
FROM
    users
WHERE
    api_token = $1
    AND deleted_at IS NULL;
