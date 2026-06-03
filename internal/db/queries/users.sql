-- name: CreateUser :one
INSERT INTO users (id, username, email, password_hash, display_name, games_played, games_won, total_points, last_login_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, username, email, password_hash, display_name, games_played, games_won, total_points, last_login_at, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, email, password_hash, display_name, games_played, games_won, total_points, last_login_at, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, display_name, games_played, games_won, total_points, last_login_at, created_at, updated_at
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, display_name, games_played, games_won, total_points, last_login_at, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET username = $1,
    email = $2,
    password_hash = $3,
    display_name = $4,
    games_played = $5,
    games_won = $6,
    total_points = $7,
    last_login_at = $8,
    updated_at = $9
WHERE id = $10
RETURNING id, username, email, password_hash, display_name, games_played, games_won, total_points, last_login_at, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserStats :one
UPDATE users
SET games_played = $1,
    games_won = $2,
    total_points = $3,
    updated_at = $4
WHERE id = $5
RETURNING id, username, email, password_hash, display_name, games_played, games_won, total_points, last_login_at, created_at, updated_at;
