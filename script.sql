CREATE TABLE profiles (
    id character varying(100) PRIMARY KEY,
    salt character varying(100) NOT NULL,
    identity_code character varying(20),
    first_name character varying(20),
    last_name character varying(50),
    gender character varying(10),
    date_of_birth character varying(10),
    phone_number character varying(20),
    email character varying(25),
    token character varying(100),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bank_profiles (
    id character varying(100) PRIMARY KEY,
    owner character varying(100) NOT NULL,
    bank_org character varying(20) NOT NULL,
    bank_code character varying(20) NOT NULL,
    owner_name character varying(50) NOT NULL,
    payos_client_id character varying(100),
    payos_api_key character varying(100),
    payos_check_sum_key character varying(100),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE center_requests (
    id character varying(100) PRIMARY KEY,
    region character varying(30) NOT NULL,
    address character varying(80) NOT NULL,
    phone_number character varying(20) NOT NULL,
    image_blob_id character varying(100) NOT NULL,
    approvers []TEXT,
    refusers []TEXT,
    refuse_reasons []TEXT,
    status character varying(10) NOT NULL,
    is_available_to_confirm BOOLEAN NOT NULL DEFAULT FALSE,
    is_confirm_register BOOLEAN NOT NULL DEFAULT FALSE,
    created_by character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    closed_at timestamp without time zone
);

CREATE TABLE payments (
    id character varying(100) PRIMARY KEY,
    actor character varying(100) NOT NULL,
    sub character varying(20) NOT NULL,
    proposal_id character varying(100),
    donation_id character varying(100),
    is_donate_tx BOOLEAN NOT NULL,
    transaction_id character varying(50) NOT NULL,
    amount BIGINT NOT NULL,
    currency character varying(10) NOT NULL,
    status character varying(10) NOT NULL,
    method character varying(10) NOT NULL,
    cancel_reason character varying(50),
    message character varying(50) NOT NULL,
    expired_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_payment_profile FOREIGN KEY (sub) REFERENCES profiles(id) ON DELETE CASCADE,
    CONSTRAINT fk_payment_proposal FOREIGN KEY (proposal_id) REFERENCES withdraw_proposals(id) ON DELETE CASCADE,
    CONSTRAINT fk_payment_donation FOREIGN KEY (donation_id) REFERENCES donations(id) ON DELETE CASCADE,
    CONSTRAINT check_single_detail CHECK (
        (proposal_id IS NOT NULL AND donation_id IS NULL) OR 
        (proposal_id IS NULL AND donation_id IS NOT NULL)
    )
);

CREATE TABLE withdraw_proposals (
    id character varying(100) PRIMARY KEY,
    purpose character varying(20) NOT NULL,
    proposal_id character varying(100),
    target character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE donations (
    id character varying(100) PRIMARY KEY,
    purpose character varying(20) NOT NULL,
    target character varying(100) NOT NULL,
    start_period character varying(20),
    end_period character varying(20),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE registration_requests (
    id                    character varying(100) PRIMARY KEY,
    register_role         character varying(10) NOT NULL,
    identity_code         character varying(20) NOT NULL,
    identity_card_blob_id character varying(100) NOT NULL,
    avatar_blob_id        character varying(100) NOT NULL,
    region                character varying(30),
    first_name            character varying(10) NOT NULL,
    last_name             character varying(50) NOT NULL,
    gender                character varying(10) NOT NULL,
    date_of_birth         character varying(10) NOT NULL,
    phone_number          character varying(20) NOT NULL,
    email                 character varying(30) NOT NULL,
    approvers             TEXT[],
    refusers              TEXT[],
    refuse_reasons        TEXT[],
    status                character varying(50) NOT NULL DEFAULT 'Pending',
    is_available_to_confirm BOOLEAN NOT NULL DEFAULT FALSE,
    is_confirm_register   BOOLEAN NOT NULL DEFAULT FALSE,
    created_by            character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    closed_at timestamp without time zone
);

CREATE TABLE registration_requests (
    id                    character varying(100) PRIMARY KEY,
    identity_code         character varying(20) NOT NULL,
    avatar_blob_id        character varying(100) NOT NULL,
    region                character varying(30),
    first_name            character varying(10) NOT NULL,
    last_name             character varying(50) NOT NULL,
    gender                character varying(10) NOT NULL,
    date_of_birth         character varying(10) NOT NULL,
    approvers             TEXT[],
    refusers              TEXT[],
    refuse_reasons        TEXT[],
    status                character varying(50) NOT NULL DEFAULT 'Pending',
    is_confirm_upload BOOLEAN NOT NULL DEFAULT FALSE,
    created_by            character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    closed_at timestamp without time zone
);

-- Index for performance on common lookups
CREATE INDEX idx_registration_status ON registration_requests(status);
CREATE INDEX idx_registration_email ON registration_requests(email);