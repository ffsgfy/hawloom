BEGIN;
    CREATE TABLE vote (
        ver UUID NOT NULL REFERENCES ver (id) ON DELETE CASCADE,
        doc UUID NOT NULL REFERENCES doc (id) ON DELETE CASCADE,
        vord_num INT NOT NULL,
        account INT NOT NULL REFERENCES account (id) ON DELETE CASCADE,

        PRIMARY KEY (ver, account),
        FOREIGN KEY (doc, vord_num) REFERENCES vord (doc, num) ON DELETE CASCADE
    );

    CREATE INDEX ON vote (doc, vord_num, account);
COMMIT;
