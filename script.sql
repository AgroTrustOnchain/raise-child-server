CREATE TABLE profiles (
    id character varying(100) PRIMARY KEY,
    salt character varying(30) NOT NULL UNIQUE,
    status character varying(10) NOT NULL DEFAULT 'Active',
    identity_code character varying(20) UNIQUE,
    first_name character varying(20),
    last_name character varying(50),
    gender character varying(10),
    date_of_birth character varying(10),
    phone_number character varying(15) UNIQUE,
    email character varying(30) UNIQUE,
    token TEXT,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE bank_profiles (
    id character varying(100) PRIMARY KEY,
    profile_id character varying(100) NOT NULL UNIQUE,
    owner character varying(100) NOT NULL,
    bank_org character varying(20) NOT NULL,
    bank_code character varying(20) NOT NULL,
    owner_name character varying(50) NOT NULL,
    payos_client_id TEXT,
    payos_api_key TEXT,
    payos_check_sum_key TEXT,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_bank_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE center_requests (
    id character varying(100) PRIMARY KEY,
    profile_id character varying(100) NOT NULL,
    region character varying(30) NOT NULL,
    address character varying(80) NOT NULL,
    phone_number character varying(20) NOT NULL,
    image_blob_id character varying(100) NOT NULL,
    approvers TEXT[] DEFAULT '{}' NOT NULL,
    refusers TEXT[] DEFAULT '{}' NOT NULL,
    refuse_reasons TEXT[] DEFAULT '{}' NOT NULL,
    status character varying(10) NOT NULL,
    is_available_to_confirm BOOLEAN NOT NULL DEFAULT FALSE,
    is_confirm_register BOOLEAN NOT NULL DEFAULT FALSE,
    created_by character varying(100) NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    closed_at timestamptz NOT NULL,
    CONSTRAINT fk_center_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE meal_support_durations (
    id character varying(100) PRIMARY KEY,
    start_period character varying(10) NOT NULL,
    end_period character varying(10) NOT NULL
);

CREATE TABLE donations (
    id character varying(100) PRIMARY KEY,
    purpose character varying(20) NOT NULL,
    target character varying(100) NOT NULL,
    meal_duration_id character varying(100) UNIQUE,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_donation_duration FOREIGN KEY (meal_duration_id) REFERENCES meal_support_durations(id) ON DELETE CASCADE
);

CREATE TABLE withdraw_proposals (
    id character varying(100) PRIMARY KEY,
    purpose character varying(20) NOT NULL,
    proposal_id character varying(100) UNIQUE,
    target character varying(100) NOT NULL,
    local_pool_id character varying(100) NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE payments (
    id character varying(100) PRIMARY KEY,
    actor character varying(100) NOT NULL,
    profile_id character varying(100) NOT NULL,
    proposal_id character varying(100),
    donation_id character varying(100),
    is_donate_tx BOOLEAN NOT NULL,
    transaction_id character varying(50) UNIQUE NOT NULL,
    amount BIGINT NOT NULL,
    currency character varying(10) NOT NULL,
    status character varying(10) NOT NULL,
    method character varying(10) NOT NULL,
    cancel_reason character varying(50),
    message character varying(50) NOT NULL,
    expired_at timestamptz NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_payment_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
    CONSTRAINT fk_payment_proposal FOREIGN KEY (proposal_id) REFERENCES withdraw_proposals(id) ON DELETE CASCADE,
    CONSTRAINT fk_payment_donation FOREIGN KEY (donation_id) REFERENCES donations(id) ON DELETE CASCADE,
    CONSTRAINT check_single_detail CHECK (
        (proposal_id IS NOT NULL AND donation_id IS NULL) OR 
        (proposal_id IS NULL AND donation_id IS NOT NULL)
    )
);

CREATE TABLE registration_requests (
    id                    character varying(100) PRIMARY KEY,
    profile_id            character varying(100) NOT NULL,
    register_role         character varying(15) NOT NULL,
    identity_code         character varying(20) NOT NULL,
    identity_card_blob_id character varying(100) NOT NULL,
    avatar_blob_id        character varying(100) NOT NULL,
    region                character varying(30),
    first_name            character varying(10) NOT NULL,
    last_name             character varying(50) NOT NULL,
    gender                character varying(10) NOT NULL,
    date_of_birth         character varying(10) NOT NULL,
    phone_number          character varying(20) NOT NULL,
    email                 character varying(25) NOT NULL,
    approvers             TEXT[] DEFAULT '{}' NOT NULL,
    refusers              TEXT[] DEFAULT '{}' NOT NULL,
    refuse_reasons        TEXT[] DEFAULT '{}' NOT NULL,
    status                character varying(10) NOT NULL DEFAULT 'Pending',
    is_available_to_confirm BOOLEAN NOT NULL DEFAULT FALSE,
    is_confirm_register   BOOLEAN NOT NULL DEFAULT FALSE,
    created_by            character varying(100) NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    closed_at timestamptz NOT NULL,
    CONSTRAINT fk_registration_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE upload_child_requests (
    id                    character varying(100) PRIMARY KEY,
    profile_id            character varying(100) NOT NULL,
    identity_code         character varying(20) NOT NULL,
    avatar_blob_id        character varying(100) NOT NULL,
    home_blob_id          character varying(100) NOT NULL,
    region                character varying(30) NOT NULL,
    first_name            character varying(10) NOT NULL,
    last_name             character varying(50) NOT NULL,
    gender                character varying(10) NOT NULL,
    date_of_birth         character varying(10) NOT NULL,
    home_address          character varying(80) NOT NULL,
    first_guardian_name   character varying(50) NOT NULL,
    first_guardian_phone  character varying(15) NOT NULL,
    first_guardian_relation character varying(10) NOT NULL,
    first_guardian_identity_card_blob_id character varying(100) NOT NULL,
    second_guardian_name   character varying(50),
    second_guardian_phone  character varying(15),
    second_guardian_relation character varying(10),
    second_guardian_identity_card_blob_id character varying(100),
    approvers             TEXT[] DEFAULT '{}' NOT NULL,
    refusers              TEXT[] DEFAULT '{}' NOT NULL,
    refuse_reasons        TEXT[] DEFAULT '{}' NOT NULL,
    ai_evaluation         character varying(50) NOT NULL,
    status                character varying(10) NOT NULL DEFAULT 'Pending',
    review_status         character varying(10) NOT NULL DEFAULT 'Pending',
    is_confirm_upload     BOOLEAN NOT NULL DEFAULT FALSE,
    created_by            character varying(100) NOT NULL,
    reviewed_by           character varying(100),
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    closed_at timestamptz,
    CONSTRAINT fk_upload_child_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE supported_region_suggestions (
    id character varying(100) PRIMARY KEY,
    profile_id character varying(100) NOT NULL,
    region character varying(30) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    status character varying(10) NOT NULL DEFAULT 'Pending',
    created_by character varying(100) NOT NULL,
    reviewed_by character varying(100),
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_region_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE pending_withdraw_proposals ( 
    id character varying(100) PRIMARY KEY,
    profile_id character varying(100) NOT NULL,
    creator character varying(100) NOT NULL,
    pool_id character varying(100) NOT NULL,
    pool_name character varying(30) NOT NULL,
    purpose character varying(30) NOT NULL,
    target character varying(100) NOT NULL,
    withdraw_amount BIGINT NOT NULL,
    proof_blob_id character varying(100),
    description TEXT NOT NULL,
    status character varying(10) NOT NULL DEFAULT 'Pending',
    ai_evaluation character varying(20) NOT NULL DEFAULT 'Pending',
    reviewed_by character varying(100),
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_proposal_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE pending_child_special_need_proposals(
    id character varying(100) PRIMARY KEY,
    child_id character varying(100) NOT NULL,
    region character varying(30) NOT NULL,
    actor_profile_id character varying(100) NOT NULL,
    actor_address character varying(100) NOT NULL,
    target BIGINT NOT NULL,
    description TEXT[] NOT NULL,
    proof_blob_id character varying(100),
    ai_evaluation character varying(50) NOT NULL,
    review_status character varying(10) NOT NULL DEFAULT 'Pending',
    reviewed_by character varying(100),
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_special_need_profile FOREIGN KEY (actor_profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE volunteer_tasks (
    id character varying(100) PRIMARY KEY,
    assigned_profile_id character varying(100),
    assgined_volunteer character varying(100),
    child_id character varying(100) NOT NULL,
    region character varying(30) NOT NULL,
    content TEXT NOT NULL,
    start_period character varying(10) NOT NULL,
    end_period character varying(10) NOT NULL,
    is_end BOOLEAN NOT NULL DEFAULT FALSE,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_volunteer_task_profile FOREIGN KEY (assigned_profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE leader_notis (
    id character varying(100) PRIMARY KEY,
    need_id character varying(100) NOT NULL UNIQUE,
    need_type character varying(20) NOT NULL,
    child_id character varying(100) NOT NULL UNIQUE,
    region character varying(30) NOT NULL,
    assigned_leaders TEXT[] DEFAULT '{}' NOT NULL,
    expected_withdraw_periods TEXT[] DEFAULT '{}' NOT NULL,
    contents TEXT[] NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE child_task_details (
    id character varying(100) PRIMARY KEY,
    child_id character varying(100) NOT NULL,
    purpose character varying(20) NOT NULL,
    target character varying(100) NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE tasks (
    id character varying(100) PRIMARY KEY,
    is_child_task BOOLEAN NOT NULL,
    child_task_id character varying(100),
    created_by character varying(100) NOT NULL,
    assigned_profile_id character varying(100),
    assgined_staff character varying(100),
    review_profile_status character varying(10) NOT NULL DEFAULT 'Pending',
    reviewed_by character varying(100),
    region character varying(30) NOT NULL,
    description TEXT[] NOT NULL,
    start_period timestamptz NOT NULL,
    end_period timestamptz NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_task_profile FOREIGN KEY (assigned_profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE task_proofs (
    id character varying(100) PRIMARY KEY,
    task_id character varying(100) NOT NULL,
    description TEXT[] NOT NULL,
    actor_profile_id character varying(100) NOT NULL,
    actor_address character varying(100) NOT NULL,
    image_blob_id character varying(100) NOT NULL,
    reviewed_by character varying(100),
    review_status character varying(10) NOT NULL DEFAULT 'Pending',
    raw_submit_date character varying(10) NOT NULL,
     created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_proof_profile FOREIGN KEY (actor_profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

-- Index for performance on common lookups
CREATE INDEX idx_registration_status ON registration_requests(status);
CREATE INDEX idx_registration_email ON registration_requests(email);
CREATE INDEX idx_center_status ON center_requests(status);



