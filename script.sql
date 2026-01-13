CREATE TABLE registration_requests (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    register_role         VARCHAR(50) NOT NULL,
    identity_code         VARCHAR(100) NOT NULL,
    identity_card_blob_id VARCHAR(255) NOT NULL,
    avatar_blob_id        VARCHAR(255) NOT NULL,
    region                VARCHAR(100) NOT NULL,
    first_name            VARCHAR(100) NOT NULL,
    last_name             VARCHAR(100) NOT NULL,
    gender                VARCHAR(20) NOT NULL,
    date_of_birth         VARCHAR(20) NOT NULL,
    phone_number          VARCHAR(20) NOT NULL,
    email                 VARCHAR(255) NOT NULL,
    approvers             TEXT[] DEFAULT '{}', -- PostgreSQL Array type
    refusers              TEXT[] DEFAULT '{}',
    refuse_reasons        TEXT[] DEFAULT '{}',
    status                VARCHAR(50) DEFAULT 'Pending',
    is_confirm_register   BOOLEAN DEFAULT FALSE,
    created_by            VARCHAR(255) NOT NULL,
    created_at            TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    closed_at             TIMESTAMPTZ
);

-- Index for performance on common lookups
CREATE INDEX idx_registration_status ON registration_requests(status);
CREATE INDEX idx_registration_email ON registration_requests(email);