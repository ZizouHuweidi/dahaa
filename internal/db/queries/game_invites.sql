-- name: CreateGameInvite :one
INSERT INTO game_invites (id, game_id, from_user, to_user, status, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, game_id, from_user, to_user, status, created_at, expires_at;

-- name: GetGameInviteByID :one
SELECT id, game_id, from_user, to_user, status, created_at, expires_at
FROM game_invites
WHERE id = $1
LIMIT 1;

-- name: GetPendingGameInvites :many
SELECT id, game_id, from_user, to_user, status, created_at, expires_at
FROM game_invites
WHERE to_user = $1
  AND status = 'pending'
  AND expires_at > $2
ORDER BY created_at DESC;

-- name: UpdateGameInviteStatus :one
UPDATE game_invites
SET status = $1
WHERE id = $2
RETURNING id, game_id, from_user, to_user, status, created_at, expires_at;

-- name: DeleteGameInvite :exec
DELETE FROM game_invites
WHERE id = $1;
