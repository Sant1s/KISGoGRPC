CREATE TABLE IF NOT EXISTS posts (
    id BIGINT PRIMARY KEY,
    author VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    create_time TIMESTAMP NOT NULL,
    comment_count BIGINT DEFAULT 0,
    likes_count BIGINT DEFAULT 0,
    liked BOOLEAN DEFAULT FALSE,
    parent_id BIGINT,
    CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES posts(id)
);