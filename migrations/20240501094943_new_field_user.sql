-- +goose Up
-- +goose StatementBegin
ALTER TABLE shorten_urls
ADD user_id VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE shorten_urls
DROP COLUMN user_id;
-- +goose StatementEnd
