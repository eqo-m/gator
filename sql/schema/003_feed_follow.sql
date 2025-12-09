-- +goose Up
CREATE TABLE feed_follow (
    id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users
    ON DELETE CASCADE,
    FOREIGN KEY (feed_id) REFERENCES feeds
    ON DELETE CASCADE,
    PRIMARY KEY(id),
    UNIQUE(user_id,feed_id)
);

-- +goose Down
DROP TABLE feed_follow;
