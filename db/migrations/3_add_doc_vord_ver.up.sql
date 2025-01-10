-- NOTE: vord = VOting RounD

BEGIN;
    CREATE TABLE doc (
        id UUID NOT NULL PRIMARY KEY,
        title VARCHAR(256) NOT NULL,
        flags INT NOT NULL,
        created_by INT NOT NULL REFERENCES account (id) ON DELETE CASCADE,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        vord_duration INT NOT NULL CHECK (vord_duration > 0),
        current_ver UUID NOT NULL
    );

    CREATE TABLE vord (
        doc UUID NOT NULL REFERENCES doc (id) ON DELETE CASCADE,
        num INT NOT NULL CHECK (num >= -1), -- -1 is the currently active vord
        flags INT NOT NULL,
        start_at TIMESTAMP WITH TIME ZONE NOT NULL,
        finish_at TIMESTAMP WITH TIME ZONE NOT NULL,

        PRIMARY KEY (doc, num)
    );

    CREATE TABLE ver (
        id UUID NOT NULL PRIMARY KEY,
        doc UUID NOT NULL REFERENCES doc (id) ON DELETE CASCADE,
        vord_num INT NOT NULL,
        created_by INT REFERENCES account (id) ON DELETE SET NULL,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        summary TEXT NOT NULL,
        content TEXT NOT NULL,

        FOREIGN KEY (doc, vord_num) REFERENCES vord (doc, num) ON DELETE CASCADE
    );

    ALTER TABLE doc
    ADD CONSTRAINT doc_current_ver_fkey
        FOREIGN KEY (current_ver) REFERENCES ver (id) ON DELETE RESTRICT
        DEFERRABLE INITIALLY DEFERRED;

    CREATE INDEX ON vord (finish_at)
    WHERE num = -1;

    CREATE INDEX ON ver (doc, vord_num);
COMMIT;
