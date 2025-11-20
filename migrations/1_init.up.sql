CREATE TABLE IF NOT EXISTS users_id
(
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS teams
(
    id SERIAL PRIMARY KEY,
    team_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users_id (id),
    team_id INTEGER REFERENCES teams (id),
    is_active BOOLEAN NOT NULL
);
