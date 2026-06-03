-- name: GetRandomQuestion :one
SELECT id, text, answer, category, filler_answers, created_at, updated_at
FROM questions
WHERE category = $1
ORDER BY RANDOM()
LIMIT 1;

-- name: GetQuestionCategories :many
SELECT DISTINCT category
FROM questions
ORDER BY category;

-- name: GetQuestionByID :one
SELECT id, text, answer, category, filler_answers, created_at, updated_at
FROM questions
WHERE id = $1
LIMIT 1;

-- name: CreateQuestion :one
INSERT INTO questions (text, answer, category, filler_answers)
VALUES ($1, $2, $3, $4)
RETURNING id, text, answer, category, filler_answers, created_at, updated_at;

-- name: UpdateQuestion :one
UPDATE questions
SET text = $1,
    answer = $2,
    category = $3,
    filler_answers = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $5
RETURNING id, text, answer, category, filler_answers, created_at, updated_at;

-- name: DeleteQuestion :execrows
DELETE FROM questions
WHERE id = $1;
