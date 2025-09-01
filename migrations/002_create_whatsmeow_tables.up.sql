-- Drop old tables if they exist
DROP TABLE IF EXISTS whatsmeow_chat_settings CASCADE;
DROP TABLE IF EXISTS whatsmeow_contacts CASCADE;
DROP TABLE IF EXISTS whatsmeow_app_state_version CASCADE;
DROP TABLE IF EXISTS whatsmeow_app_state_sync_keys CASCADE;
DROP TABLE IF EXISTS whatsmeow_sender_keys CASCADE;
DROP TABLE IF EXISTS whatsmeow_sessions CASCADE;
DROP TABLE IF EXISTS whatsmeow_pre_keys CASCADE;
DROP TABLE IF EXISTS whatsmeow_identity_keys CASCADE;
DROP TABLE IF EXISTS whatsmeow_device CASCADE;

-- Create whatsmeow compatible tables
CREATE TABLE whatsmeow_device (
    jid         TEXT PRIMARY KEY,
    registration_id BIGINT NOT NULL CHECK (registration_id >= 0 AND registration_id < 4294967296),
    noise_key       BYTEA NOT NULL CHECK (length(noise_key) = 32),
    identity_key    BYTEA NOT NULL CHECK (length(identity_key) = 32),
    signed_pre_key  BYTEA NOT NULL,
    signed_pre_key_id INTEGER NOT NULL CHECK (signed_pre_key_id >= 0 AND signed_pre_key_id < 16777216),
    signed_pre_key_sig BYTEA NOT NULL CHECK (length(signed_pre_key_sig) = 64),
    adv_secret_key  BYTEA NOT NULL CHECK (length(adv_secret_key) = 32),
    next_pre_key_id INTEGER NOT NULL DEFAULT 1,
    first_unuploaded_pre_key_id INTEGER NOT NULL DEFAULT 1,
    account_sync_timestamp BIGINT NOT NULL DEFAULT 0,
    myapp_state_key_id BYTEA CHECK (myapp_state_key_id IS NULL OR length(myapp_state_key_id) = 32),
    platform TEXT NOT NULL DEFAULT '',
    business_name TEXT NOT NULL DEFAULT '',
    push_name TEXT NOT NULL DEFAULT ''
);

CREATE TABLE whatsmeow_identity_keys (
    our_jid   TEXT,
    their_id  TEXT,
    identity  BYTEA NOT NULL CHECK (length(identity) = 32),
    PRIMARY KEY (our_jid, their_id),
    FOREIGN KEY (our_jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE
);

CREATE TABLE whatsmeow_pre_keys (
    jid     TEXT,
    key_id  INTEGER,
    key     BYTEA NOT NULL CHECK (length(key) = 32),
    uploaded BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY (jid, key_id),
    FOREIGN KEY (jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE
);

CREATE TABLE whatsmeow_sessions (
    our_jid   TEXT,
    their_jid TEXT,
    session   BYTEA NOT NULL,
    PRIMARY KEY (our_jid, their_jid),
    FOREIGN KEY (our_jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE
);

CREATE TABLE whatsmeow_sender_keys (
    our_jid    TEXT,
    chat_jid   TEXT,
    sender_jid TEXT,
    sender_key BYTEA NOT NULL,
    PRIMARY KEY (our_jid, chat_jid, sender_jid),
    FOREIGN KEY (our_jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE
);

CREATE TABLE whatsmeow_app_state_sync_keys (
    jid        TEXT,
    key_id     BYTEA,
    key_data   BYTEA NOT NULL,
    timestamp  BIGINT NOT NULL,
    fingerprint BYTEA NOT NULL,
    PRIMARY KEY (jid, key_id),
    FOREIGN KEY (jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE
);

CREATE TABLE whatsmeow_app_state_version (
    jid      TEXT,
    name     TEXT,
    version  BIGINT NOT NULL DEFAULT 0,
    hash     BYTEA NOT NULL CHECK (length(hash) = 128),
    PRIMARY KEY (jid, name),
    FOREIGN KEY (jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE
);

CREATE TABLE whatsmeow_contacts (
    our_jid      TEXT,
    their_jid    TEXT,
    first_name   TEXT,
    full_name    TEXT,
    push_name    TEXT,
    business_name TEXT,
    PRIMARY KEY (our_jid, their_jid),
    FOREIGN KEY (our_jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE
);

CREATE TABLE whatsmeow_chat_settings (
    our_jid           TEXT,
    chat_jid          TEXT,
    muted_until       BIGINT NOT NULL DEFAULT 0,
    pinned            BOOLEAN NOT NULL DEFAULT false,
    archived          BOOLEAN NOT NULL DEFAULT false,
    disappearing_timer BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (our_jid, chat_jid),
    FOREIGN KEY (our_jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE
);

-- Additional tables that might be needed
CREATE TABLE whatsmeow_message_secrets (
    our_jid TEXT,
    chat_jid TEXT,
    sender_jid TEXT,
    message_id TEXT,
    secret BYTEA NOT NULL,
    PRIMARY KEY (our_jid, chat_jid, sender_jid, message_id),
    FOREIGN KEY (our_jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE
);

CREATE TABLE whatsmeow_privacy_tokens (
    our_jid TEXT,
    their_jid TEXT,
    token BYTEA NOT NULL,
    timestamp BIGINT NOT NULL,
    PRIMARY KEY (our_jid, their_jid),
    FOREIGN KEY (our_jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX whatsmeow_identity_keys_our_jid ON whatsmeow_identity_keys(our_jid);
CREATE INDEX whatsmeow_pre_keys_jid ON whatsmeow_pre_keys(jid);
CREATE INDEX whatsmeow_sessions_our_jid ON whatsmeow_sessions(our_jid);
CREATE INDEX whatsmeow_sender_keys_our_jid ON whatsmeow_sender_keys(our_jid, chat_jid);
CREATE INDEX whatsmeow_app_state_sync_keys_jid ON whatsmeow_app_state_sync_keys(jid);
CREATE INDEX whatsmeow_contacts_our_jid ON whatsmeow_contacts(our_jid);
CREATE INDEX whatsmeow_chat_settings_our_jid ON whatsmeow_chat_settings(our_jid);

