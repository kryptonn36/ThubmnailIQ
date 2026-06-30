-- +goose Up
ALTER TABLE analyses ADD COLUMN relevance_score INT;
ALTER TABLE analyses ADD COLUMN relevance_reasoning TEXT;

-- +goose Down
ALTER TABLE analyses DROP COLUMN relevance_reasoning;
ALTER TABLE analyses DROP COLUMN relevance_score;
