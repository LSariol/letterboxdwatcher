-- Create Tables
CREATE SEQUENCE IF NOT EXISTS botsuite.feeds_id_seq;

CREATE TABLE botsuite.letterboxd_feed_subscriptions (
    id                  BIGINT PRIMARY KEY DEFAULT nextval('botsuite.feeds_id_seq'),
    username            TEXT NOT NULL,
    user_id             TEXT NOT NULL,
    letterboxd_username VARCHAR(255) NOT NULL,
    feed_url            TEXT NOT NULL,
    last_seen_guid      TEXT,
    last_alerted_at     TIMESTAMPTZ,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    alert_channels      TEXT[] NOT NULL DEFAULT '{}',
    CONSTRAINT feeds_feed_url_key UNIQUE (feed_url),
    CONSTRAINT feeds_user_id_letterboxd_username_key UNIQUE (user_id, letterboxd_username)
);
