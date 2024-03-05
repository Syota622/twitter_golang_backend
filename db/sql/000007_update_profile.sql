-- name: UpdateUserProfile :exec
UPDATE users
SET username = $2, email = $3, bio = $4, profile_image_url = $5, background_image_url = $6
WHERE id = $1;
