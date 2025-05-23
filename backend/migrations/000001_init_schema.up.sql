-- Create games table
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(6) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL,
    players JSONB NOT NULL DEFAULT '[]',
    rounds JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create questions table
CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text TEXT NOT NULL,
    answer TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    filler_answers TEXT[] NOT NULL DEFAULT '{}', -- Array of pre-defined filler answers
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT category_length CHECK (char_length(category) >= 3)
);

-- Create indexes
CREATE INDEX idx_games_code ON games(code);
CREATE INDEX idx_questions_category ON questions(category);

-- Add constraints
ALTER TABLE games
    ADD CONSTRAINT games_status_check 
    CHECK (status IN ('waiting', 'playing', 'finished'));

ALTER TABLE questions
    ADD CONSTRAINT questions_category_length_check 
    CHECK (char_length(category) >= 2 AND char_length(category) <= 50);

-- Add comments
COMMENT ON TABLE games IS 'Stores game sessions and their state';
COMMENT ON TABLE questions IS 'Stores trivia questions with optional images';
COMMENT ON COLUMN games.players IS 'JSON array of player objects';
COMMENT ON COLUMN games.rounds IS 'JSON array of round objects';
COMMENT ON COLUMN questions.image_path IS 'Path to the question image file';
COMMENT ON COLUMN questions.image_alt IS 'Alt text for the question image'; 