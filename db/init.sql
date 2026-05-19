-- Create Tables
CREATE TABLE botsuite.letterboxd_feed_subscriptions (
    id                  BIGSERIAL PRIMARY KEY,
    username            TEXT NOT NULL,
    user_id             TEXT NOT NULL,

    letterboxd_username VARCHAR(255) NOT NULL,        -- e.g. "John's Letterboxd"
    feed_url            TEXT NOT NULL UNIQUE,

    last_seen_guid      TEXT,
    last_alerted_at     TIMESTAMPTZ,

    created_at          TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (user_id, letterboxd_username)
);

CREATE TABLE botsuite.feed_entries (
    id              SERIAL PRIMARY KEY,
    feed_id         INTEGER NOT NULL REFERENCES letterboxd.feeds(id) ON DELETE CASCADE,
    guid            TEXT NOT NULL,
    title           TEXT,                     -- "Land of Bad, 2024 - ★★★"
    entry_url       TEXT,
    film_title      TEXT,                     -- "Land of Bad"
    film_year       SMALLINT,                 -- 2024
    movie_id        INTEGER,                  -- 969492
    member_rating   NUMERIC(2,1),             -- 3.0
    is_rewatch      BOOLEAN DEFAULT FALSE,
    watched_date    DATE,                     -- 2026-04-18 (more precise than pubDate)
    published_at    TIMESTAMPTZ,              -- full pubDate timestamp
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (feed_id, guid)
);