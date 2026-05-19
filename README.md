# BoxdWatcher

A Go service that polls Letterboxd RSS feeds on a schedule, detects new activity, and sends notifications to a downstream service (e.g. a Twitch bot).

## How it works

On each poll cycle, BoxdWatcher:

1. Fetches active subscriptions from PostgreSQL
2. HTTP GETs each user's Letterboxd RSS feed and parses it
3. Compares against the last seen GUID to find only new entries
4. Formats a notification message (single film or batch) and POSTs it to a notification endpoint
5. Updates the last seen GUID in the database

A failed feed does not stop the cycle — errors are logged and skipped.

## Configuration

All config is via environment variables:

| Variable | Required | Default | Description |
|---|---|---|---|
| `DATABASE_URL` | Yes | — | PostgreSQL DSN (`postgres://user:pass@host/db?sslmode=disable`) |
| `NOTIFICATION_ENDPOINT` | Yes | — | URL of the downstream notification service |
| `POLL_INTERVAL_MINUTES` | No | `20` | How often to poll feeds |

## Running with Docker

```sh
docker compose up -d
```

## Database

PostgreSQL schema is in [db/init.sql](db/init.sql). Subscriptions are stored in `botsuite.letterboxd_feed_subscriptions` with one row per tracked user, including their Letterboxd feed URL and last seen GUID.

## Notification format

- **Single film:** `🎬 @username watched FilmTitle (Year) on date | ★★½ | ❤️`
- **Multiple films:** lists up to 480 characters, then appends a profile link
