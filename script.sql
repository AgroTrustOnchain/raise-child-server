CREATE TABLE public.background_children_withdraw_requests (
  id character varying NOT NULL,
  profile_id character varying NOT NULL,
  actor_address character varying NOT NULL,
  region character varying NOT NULL,
  is_executed boolean NOT NULL DEFAULT false,
  raw_proposed_date character varying NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT background_children_withdraw_requests_pkey PRIMARY KEY (id),
  CONSTRAINT fk_background_children_withdraw_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.bank_profiles (
  id character varying NOT NULL,
  profile_id character varying NOT NULL UNIQUE,
  owner character varying NOT NULL,
  bank_org character varying NOT NULL,
  bank_code character varying NOT NULL,
  owner_name character varying NOT NULL,
  payos_client_id text,
  payos_api_key text,
  payos_check_sum_key text,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT bank_profiles_pkey PRIMARY KEY (id),
  CONSTRAINT fk_bank_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.center_requests (
  id character varying NOT NULL,
  profile_id character varying NOT NULL,
  region character varying NOT NULL,
  address character varying NOT NULL,
  phone_number character varying NOT NULL,
  image_blob_id character varying NOT NULL,
  approvers ARRAY,
  refusers ARRAY,
  refuse_reasons ARRAY,
  status character varying NOT NULL,
  is_confirm_register boolean NOT NULL DEFAULT false,
  created_by character varying NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  closed_at timestamp with time zone NOT NULL,
  CONSTRAINT center_requests_pkey PRIMARY KEY (id),
  CONSTRAINT fk_center_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.child_task_details (
  id character varying NOT NULL,
  child_id character varying NOT NULL,
  purpose character varying NOT NULL,
  target character varying NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT child_task_details_pkey PRIMARY KEY (id)
);

CREATE TABLE public.donations (
  id character varying NOT NULL,
  purpose character varying NOT NULL,
  target character varying NOT NULL,
  meal_duration_id character varying UNIQUE,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT donations_pkey PRIMARY KEY (id),
  CONSTRAINT fk_donation_duration FOREIGN KEY (meal_duration_id) REFERENCES public.meal_support_durations(id)
);

CREATE TABLE public.leader_notis (
  id character varying NOT NULL,
  need_id character varying NOT NULL UNIQUE,
  need_type character varying NOT NULL,
  child_id character varying NOT NULL UNIQUE,
  region character varying NOT NULL,
  assigned_leaders ARRAY NOT NULL DEFAULT '{}'::text[],
  expected_withdraw_periods ARRAY NOT NULL DEFAULT '{}'::text[],
  general_content text NOT NULL,
  contents ARRAY NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT leader_notis_pkey PRIMARY KEY (id)
);

CREATE TABLE public.meal_support_durations (
  id character varying NOT NULL,
  start_period character varying NOT NULL,
  end_period character varying NOT NULL,
  CONSTRAINT meal_support_durations_pkey PRIMARY KEY (id)
);

CREATE TABLE public.numeric_configs (
  id character varying NOT NULL,
  key character varying NOT NULL UNIQUE,
  value bigint NOT NULL,
  description text NOT NULL,
  actor_profile_id character varying,
  actor_address character varying,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT numeric_configs_pkey PRIMARY KEY (id),
  CONSTRAINT fk_num_config_profile FOREIGN KEY (actor_profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.payments (
  id character varying NOT NULL,
  actor character varying NOT NULL,
  profile_id character varying NOT NULL,
  proposal_id character varying,
  donation_id character varying,
  is_donate_tx boolean NOT NULL,
  transaction_id character varying NOT NULL UNIQUE,
  amount bigint NOT NULL,
  currency character varying NOT NULL,
  status character varying NOT NULL,
  method character varying NOT NULL,
  cancel_reason character varying,
  message text NOT NULL,
  expired_at timestamp with time zone NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  proof_blob_id character varying,
  reviewed_by character varying,
  review_status character varying NOT NULL DEFAULT 'Pending'::character varying,
  is_transferred boolean NOT NULL DEFAULT false,
  transferred_at timestamp with time zone,
  CONSTRAINT payments_pkey PRIMARY KEY (id),
  CONSTRAINT fk_payment_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id),
  CONSTRAINT fk_payment_donation FOREIGN KEY (donation_id) REFERENCES public.donations(id)
);

CREATE TABLE public.pending_campaigns (
  id character varying NOT NULL,
  actor_profile_id character varying NOT NULL,
  actor_address character varying NOT NULL,
  target bigint NOT NULL,
  description text NOT NULL,
  proof_blob_id character varying NOT NULL,
  ai_evaluation character varying NOT NULL,
  review_status character varying NOT NULL DEFAULT 'Pending'::character varying,
  reviewed_by character varying,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT pending_campaigns_pkey PRIMARY KEY (id),
  CONSTRAINT fk_campaign_profile FOREIGN KEY (actor_profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.pending_child_special_need_proposals (
  id character varying NOT NULL,
  child_id character varying NOT NULL,
  region character varying NOT NULL,
  actor_profile_id character varying NOT NULL,
  actor_address character varying NOT NULL,
  target bigint NOT NULL,
  description text NOT NULL,
  proof_blob_id character varying,
  ai_evaluation character varying NOT NULL,
  review_status character varying NOT NULL DEFAULT 'Pending'::character varying,
  reviewed_by character varying,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT pending_child_special_need_proposals_pkey PRIMARY KEY (id),
  CONSTRAINT fk_special_need_profile FOREIGN KEY (actor_profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.pending_withdraw_proposals (
  id character varying NOT NULL,
  profile_id character varying NOT NULL,
  creator character varying NOT NULL,
  pool_id character varying NOT NULL,
  pool_name character varying NOT NULL,
  purpose character varying NOT NULL,
  target character varying NOT NULL,
  withdraw_amount bigint NOT NULL,
  proof_blob_id character varying,
  description text NOT NULL,
  status character varying NOT NULL DEFAULT 'Pending'::character varying,
  ai_evaluation character varying NOT NULL DEFAULT 'Pending'::character varying,
  reviewed_by character varying,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT pending_withdraw_proposals_pkey PRIMARY KEY (id),
  CONSTRAINT fk_proposal_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.profiles (
  id character varying NOT NULL,
  salt character varying NOT NULL UNIQUE,
  status character varying NOT NULL DEFAULT 'Active'::character varying,
  identity_code character varying UNIQUE,
  first_name character varying,
  last_name character varying,
  gender character varying,
  date_of_birth character varying,
  phone_number character varying UNIQUE,
  email character varying UNIQUE,
  token text,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  wallet_address character varying,
  CONSTRAINT profiles_pkey PRIMARY KEY (id)
);

CREATE TABLE public.registration_requests (
  id character varying NOT NULL,
  profile_id character varying NOT NULL,
  register_role character varying NOT NULL,
  identity_code character varying NOT NULL,
  identity_card_blob_id character varying NOT NULL,
  avatar_blob_id character varying NOT NULL,
  region character varying,
  first_name character varying NOT NULL,
  last_name character varying NOT NULL,
  gender character varying NOT NULL,
  date_of_birth character varying NOT NULL,
  phone_number character varying NOT NULL,
  email character varying NOT NULL,
  approvers ARRAY,
  refusers ARRAY,
  refuse_reasons ARRAY,
  status character varying NOT NULL DEFAULT 'Pending'::character varying,
  is_confirm_register boolean NOT NULL DEFAULT false,
  created_by character varying NOT NULL,
  on_chain_id character varying UNIQUE,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  closed_at timestamp with time zone NOT NULL,
  CONSTRAINT registration_requests_pkey PRIMARY KEY (id),
  CONSTRAINT fk_registration_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.supported_region_suggestions (
  id character varying NOT NULL,
  profile_id character varying NOT NULL,
  region character varying NOT NULL UNIQUE,
  content text NOT NULL,
  status character varying NOT NULL DEFAULT 'Pending'::character varying,
  refuse_reason text,
  created_by character varying NOT NULL,
  reviewed_by character varying,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT supported_region_suggestions_pkey PRIMARY KEY (id),
  CONSTRAINT fk_region_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.task_proofs (
  id character varying NOT NULL,
  task_id character varying NOT NULL,
  description text NOT NULL DEFAULT ''::text,
  actor_profile_id character varying NOT NULL,
  actor_address character varying NOT NULL,
  image_walrus_blob_id character varying NOT NULL,
  image_cloudinary_blob_id character varying NOT NULL,
  reviewed_by character varying,
  ai_evaluation character varying NOT NULL,
  ai_reason text NOT NULL,
  review_status character varying NOT NULL DEFAULT 'Pending'::character varying,
  raw_submit_date character varying NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT task_proofs_pkey PRIMARY KEY (id),
  CONSTRAINT fk_proof_profile FOREIGN KEY (actor_profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.tasks (
  id character varying NOT NULL,
  is_child_task boolean NOT NULL,
  child_task_id character varying,
  created_by character varying NOT NULL,
  assigned_profile_id character varying,
  assgined_staff character varying,
  region character varying NOT NULL,
  description text NOT NULL,
  start_period timestamp with time zone NOT NULL,
  end_period timestamp with time zone NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT tasks_pkey PRIMARY KEY (id),
  CONSTRAINT fk_task_profile FOREIGN KEY (assigned_profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.upload_child_requests (
  id character varying NOT NULL,
  profile_id character varying NOT NULL,
  identity_code character varying NOT NULL,
  avatar_blob_id character varying NOT NULL,
  home_blob_id character varying NOT NULL,
  region character varying NOT NULL,
  first_name character varying NOT NULL,
  last_name character varying NOT NULL,
  gender character varying NOT NULL,
  date_of_birth character varying NOT NULL,
  home_address character varying NOT NULL,
  first_guardian_name character varying NOT NULL,
  first_guardian_phone character varying NOT NULL,
  first_guardian_relation character varying NOT NULL,
  first_guardian_identity_card_blob_id character varying NOT NULL,
  second_guardian_name character varying,
  second_guardian_phone character varying,
  second_guardian_relation character varying,
  second_guardian_identity_card_blob_id character varying,
  status character varying NOT NULL DEFAULT 'Pending'::character varying,
  is_confirm_upload boolean NOT NULL DEFAULT false,
  created_by character varying NOT NULL,
  reviewed_by character varying,
  on_chain_id character varying UNIQUE,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  birth_certificate_blob_id character varying NOT NULL,
  CONSTRAINT upload_child_requests_pkey PRIMARY KEY (id),
  CONSTRAINT fk_upload_child_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id)
);

CREATE TABLE public.withdraw_proposals (
  id character varying NOT NULL,
  purpose character varying NOT NULL,
  proposal_id character varying UNIQUE,
  target character varying NOT NULL,
  local_pool_id character varying NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT withdraw_proposals_pkey PRIMARY KEY (id)
);