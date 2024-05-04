-- +goose Up
-- +goose StatementBegin
ALTER TABLE shorten_urls
ADD is_deleted BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE shorten_urls
DROP COLUMN is_deleted;
-- +goose StatementEnd
