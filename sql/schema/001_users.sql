-- +goose Up
CREATE TABLE users (
    id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name varchar(255) NOT NULL,
    PRIMARY KEY(ID),
    UNIQUE(name) 
);

-- +goose Down
DROP TABLE users;