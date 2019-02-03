-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE notes(
    id VARCHAR(36) NOT NULL,
    account_id VARCHAR(36) NOT NULL,
    note text DEFAULT "",
    created timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    KEY notes_created_idx (created),
    KEY notes_updated_idx (updated)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE notes;
