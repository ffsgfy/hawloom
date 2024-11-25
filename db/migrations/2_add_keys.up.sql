BEGIN;
    CREATE TABLE key (
        id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        data BYTEA NOT NULL
    );
COMMIT;