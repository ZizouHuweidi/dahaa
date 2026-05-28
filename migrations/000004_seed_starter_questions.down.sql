-- +goose Down
-- +goose StatementBegin
DELETE FROM questions
WHERE (category, text) IN (
  ('geography', 'Which city is known as the City of a Thousand Minarets?'),
  ('geography', 'What is the largest desert in the world?'),
  ('geography', 'Which river runs through Baghdad?'),
  ('literature', 'Who wrote One Hundred Years of Solitude?'),
  ('literature', 'In which language was The Divine Comedy originally written?'),
  ('literature', 'Who is the author of Things Fall Apart?'),
  ('history', 'In which year did the Berlin Wall fall?'),
  ('history', 'Which empire built Machu Picchu?'),
  ('history', 'Who was the first woman to win a Nobel Prize?'),
  ('science', 'What particle has a negative electric charge?'),
  ('science', 'What is the chemical symbol for gold?'),
  ('science', 'Which planet has the most prominent ring system?')
);
-- +goose StatementEnd
