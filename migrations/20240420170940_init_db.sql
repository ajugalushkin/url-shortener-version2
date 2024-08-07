-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shorten_urls
(
    short_url
    VARCHAR
(
    20
) NOT NULL PRIMARY KEY,
    correlation_id VARCHAR
(
    250
) NOT NULL DEFAULT '',
    original_url VARCHAR
(
    2048
) NOT NULL
    );
CREATE UNIQUE INDEX original_url ON shorten_urls (original_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shorten_urls;
-- +goose StatementEnd
