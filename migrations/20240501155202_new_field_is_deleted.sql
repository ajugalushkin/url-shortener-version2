-- +goose Up
-- +goose StatementBegin
ALTER TABLE shorten_urls
ADD is_deleted BOOLEAN;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE shorten_urls
DROP COLUMN is_deleted;
-- +goose StatementEnd
