-- Drop indexes
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_game_invites_game_id;
DROP INDEX IF EXISTS idx_game_invites_to_user;
DROP INDEX IF EXISTS idx_game_invites_status;

-- Drop tables
DROP TABLE IF EXISTS game_invites;
DROP TABLE IF EXISTS users; 