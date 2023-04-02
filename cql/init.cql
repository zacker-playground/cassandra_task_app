CREATE KEYSPACE app WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 2};

CREATE TABLE IF NOT EXISTS app.users
(
    id         UUID PRIMARY KEY,
    name       TEXT,
    age        INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);