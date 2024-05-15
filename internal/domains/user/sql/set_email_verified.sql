-- setEmailVerifiedQuery
-- $1 - user_id
UPDATE
    users
SET
    is_email_verified = true,
    updated_at = now()
WHERE
    user_id = $1
    AND deleted_at IS NULL;
