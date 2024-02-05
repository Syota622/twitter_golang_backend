-- name: ConfirmUserEmail :exec
UPDATE users
SET is_confirmed = TRUE
WHERE confirmation_token = $1 AND is_confirmed = FALSE;
