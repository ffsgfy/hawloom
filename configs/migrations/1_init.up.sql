BEGIN;
    CREATE TABLE account (
        id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        name VARCHAR(36) UNIQUE NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        password_hash BYTEA NOT NULL,
        password_salt BYTEA NOT NULL
    );
COMMIT;
