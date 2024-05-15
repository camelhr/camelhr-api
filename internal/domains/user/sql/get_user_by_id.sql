-- getUserByIDQuery
-- $1: user_id
SELECT
    user_id,
    organization_id,
    email,
    password_hash,
    api_token,
    is_owner,
    is_email_verified,
    disabled_at,
    comment,
    created_at,
    updated_at,
    deleted_at
FROM
    users
WHERE
    user_id = $1
    AND deleted_at IS NULL;
