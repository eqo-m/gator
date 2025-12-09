-- +goose Up
CREATE TABLE posts (
    id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title varchar(255),
    url varchar(2048) NOT NULL,
    description varchar(2048),
    published TIMESTAMP,
    feed_id UUID NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY (feed_id) REFERENCES feeds(id)
    ON DELETE CASCADE 
);

-- +goose Down
DROP TABLE posts;
