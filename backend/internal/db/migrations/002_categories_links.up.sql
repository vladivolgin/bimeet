ALTER TABLE events
    ADD COLUMN IF NOT EXISTS category  TEXT NOT NULL DEFAULT 'ordinary'
        CHECK (category IN ('ordinary', 'business')),
    ADD COLUMN IF NOT EXISTS dress_code TEXT;

CREATE TABLE IF NOT EXISTS event_links (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    url         TEXT NOT NULL,
    created_by  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_event_links_event ON event_links(event_id);
