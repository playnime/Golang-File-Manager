ALTER TABLE users
    DROP COLUMN totp_secret,
    DROP COLUMN totp_enabled;

ALTER TABLE sessions DROP COLUMN confirmed;