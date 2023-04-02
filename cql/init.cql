CREATE KEYSPACE app WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 2};

CREATE TABLE IF NOT EXISTS app.users
(
    id         UUID PRIMARY KEY,
    name       TEXT,
    age        INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS app.tasks
(
    id         UUID,
    user_id    UUID,
    title      TEXT,
    checked    BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY ((user_id), id )
)