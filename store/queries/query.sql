-- name: GetAddress :one
SELECT *
FROM addresses
WHERE id = ?;
