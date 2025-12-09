-- +goose Up
CREATE TABLE feeds (
    id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name varchar(255) NOT NULL,
    url varchar(2048) NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY(id),
    UNIQUE(name),
    UNIQUE(url),
    FOREIGN KEY (user_id) REFERENCES users
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
