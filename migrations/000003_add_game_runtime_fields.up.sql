-- +goose Up
-- +goose StatementBegin
ALTER TABLE games
ADD COLUMN settings JSONB NOT NULL DEFAULT '{}',
ADD COLUMN last_activity TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
ADD COLUMN host_id TEXT NOT NULL DEFAULT '';

ALTER TABLE games DROP CONSTRAINT IF EXISTS games_status_check;
ALTER TABLE games
ADD CONSTRAINT games_status_check CHECK (status IN ('waiting', 'playing', 'ended'));
-- +goose StatementEnd
