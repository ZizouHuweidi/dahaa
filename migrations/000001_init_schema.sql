-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(6) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL,
    players JSONB NOT NULL DEFAULT '[]',
    rounds JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text TEXT NOT NULL,
    answer TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    filler_answers TEXT [] NOT NULL DEFAULT '{}',
    image_path TEXT,
    image_alt TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT category_length CHECK (char_length(category) >= 3)
);

CREATE INDEX idx_games_code ON games(code);
CREATE INDEX idx_questions_category ON questions(category);

ALTER TABLE games
ADD CONSTRAINT games_status_check CHECK (status IN ('waiting', 'playing', 'finished'));
ALTER TABLE questions
ADD CONSTRAINT questions_category_length_check CHECK (
        char_length(category) >= 2
        AND char_length(category) <= 50
    );

CREATE TRIGGER update_games_updated_at BEFORE
UPDATE ON games FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_questions_updated_at BEFORE
UPDATE ON questions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE games IS 'Stores game sessions and their state';
COMMENT ON TABLE questions IS 'Stores trivia questions with optional images';
COMMENT ON COLUMN games.players IS 'JSON array of player objects';
COMMENT ON COLUMN games.rounds IS 'JSON array of round objects';
COMMENT ON COLUMN questions.image_path IS 'Path to the question image file';
COMMENT ON COLUMN questions.image_alt IS 'Alt text for the question image';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
COMMENT ON TABLE games IS NULL;
COMMENT ON TABLE questions IS NULL;
COMMENT ON COLUMN games.players IS NULL;
COMMENT ON COLUMN games.rounds IS NULL;
COMMENT ON COLUMN questions.image_path IS NULL;
COMMENT ON COLUMN questions.image_alt IS NULL;

ALTER TABLE games DROP CONSTRAINT IF EXISTS games_status_check;
ALTER TABLE questions DROP CONSTRAINT IF EXISTS questions_category_length_check;
ALTER TABLE questions DROP CONSTRAINT IF EXISTS category_length;

DROP INDEX IF EXISTS idx_games_code;
DROP INDEX IF EXISTS idx_questions_category;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS questions;
DROP FUNCTION IF EXISTS update_updated_at_column();
-- +goose StatementEnd
