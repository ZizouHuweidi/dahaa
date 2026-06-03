-- name: CreateGame :one
INSERT INTO games (code, status, players, rounds, settings, last_activity, host_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, code, status, players, rounds, settings, created_at, updated_at, last_activity, host_id;

-- name: GetGameByCode :one
SELECT id, code, status, players, rounds, settings, created_at, updated_at, last_activity, host_id
FROM games
WHERE code = $1
LIMIT 1;

-- name: GetGameByID :one
SELECT id, code, status, players, rounds, settings, created_at, updated_at, last_activity, host_id
FROM games
WHERE id = $1
LIMIT 1;

-- name: UpdateGame :one
UPDATE games
SET status = $1,
    players = $2,
    rounds = $3,
    settings = $4,
    last_activity = $5,
    host_id = $6,
    updated_at = $7
WHERE code = $8
RETURNING id, code, status, players, rounds, settings, created_at, updated_at, last_activity, host_id;

-- name: DeleteGameByCode :exec
DELETE FROM games
WHERE code = $1;
