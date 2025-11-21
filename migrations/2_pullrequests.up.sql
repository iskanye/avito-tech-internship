CREATE TABLE IF NOT EXISTS pull_requests_id
(
    id SERIAL PRIMARY KEY,
    pull_request_id TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS pull_requests
(
    id SERIAL PRIMARY KEY,
    pull_request_id INTEGER REFERENCES pull_requests_id (id),
    pull_request_name TEXT
    author_id INTEGER REFERENCES users (id),
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    merged_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reviewers
(
    pull_request_id INTEGER REFERENCES pull_requests (id),
    user_id INTEGER REFERENCES users (id)
);
