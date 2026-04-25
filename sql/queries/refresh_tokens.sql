-- name: AddRefreshToken :exec
INSERT INTO refresh_tokens(token) VALUES($1);

-- name: GetRefreshToken :one
SELECT token FROM refresh_tokens WHERE token=$1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE token=$1;
