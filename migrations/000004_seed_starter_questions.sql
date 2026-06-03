-- +goose Up
-- +goose StatementBegin
INSERT INTO questions (text, answer, category, filler_answers)
VALUES
  ('Which city is known as the City of a Thousand Minarets?', 'Cairo', 'geography', ARRAY['Istanbul', 'Damascus', 'Marrakesh', 'Rabat', 'Baghdad']),
  ('What is the largest desert in the world?', 'Antarctic Desert', 'geography', ARRAY['Sahara Desert', 'Arabian Desert', 'Gobi Desert', 'Kalahari Desert', 'Patagonian Desert']),
  ('Which river runs through Baghdad?', 'Tigris', 'geography', ARRAY['Euphrates', 'Nile', 'Jordan', 'Orontes', 'Indus']),
  ('Who wrote One Hundred Years of Solitude?', 'Gabriel Garcia Marquez', 'literature', ARRAY['Pablo Neruda', 'Jorge Luis Borges', 'Isabel Allende', 'Mario Vargas Llosa', 'Julio Cortazar']),
  ('In which language was The Divine Comedy originally written?', 'Italian', 'literature', ARRAY['Latin', 'Greek', 'French', 'Spanish', 'German']),
  ('Who is the author of Things Fall Apart?', 'Chinua Achebe', 'literature', ARRAY['Wole Soyinka', 'Ngugi wa Thiongo', 'Ben Okri', 'Nadine Gordimer', 'Ayi Kwei Armah']),
  ('In which year did the Berlin Wall fall?', '1989', 'history', ARRAY['1987', '1988', '1990', '1991', '1992']),
  ('Which empire built Machu Picchu?', 'Inca Empire', 'history', ARRAY['Aztec Empire', 'Maya Civilization', 'Spanish Empire', 'Moche Civilization', 'Wari Empire']),
  ('Who was the first woman to win a Nobel Prize?', 'Marie Curie', 'history', ARRAY['Rosalind Franklin', 'Ada Lovelace', 'Dorothy Hodgkin', 'Irène Joliot-Curie', 'Lise Meitner']),
  ('What particle has a negative electric charge?', 'Electron', 'science', ARRAY['Proton', 'Neutron', 'Photon', 'Muon', 'Neutrino']),
  ('What is the chemical symbol for gold?', 'Au', 'science', ARRAY['Ag', 'Gd', 'Go', 'Pb', 'Pt']),
  ('Which planet has the most prominent ring system?', 'Saturn', 'science', ARRAY['Jupiter', 'Uranus', 'Neptune', 'Mars', 'Venus']);
-- +goose StatementEnd

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
