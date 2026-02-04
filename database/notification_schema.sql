CREATE TABLE email_logs (
    id         BIGSERIAL PRIMARY KEY,
    from_name  TEXT NOT NULL,
    from_email TEXT NOT NULL,
    to_email   TEXT NOT NULL,
    subject    TEXT NOT NULL,
    status     TEXT NOT NULL, -- sent | failed
    error      TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
