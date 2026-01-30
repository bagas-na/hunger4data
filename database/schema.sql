CREATE TABLE payments (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    country_id UUID NOT NULL,

    transaction_type VARCHAR(32) NOT NULL, -- payment, refund

    amount BIGINT NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'IDR',

    provider VARCHAR(32) NOT NULL,
    provider_object_id VARCHAR(255),

    status VARCHAR(32) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE payment_events (
    id UUID PRIMARY KEY,
    payment_id UUID NOT NULL REFERENCES payments(id),

    event_type VARCHAR(32) NOT NULL, -- created, pending, paid, failed, expired
    source VARCHAR(16) NOT NULL, -- system, webhook

    provider_event_id VARCHAR(255) UNIQUE,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);
