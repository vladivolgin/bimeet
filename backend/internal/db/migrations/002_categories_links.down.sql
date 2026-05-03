DROP TABLE IF EXISTS event_links;

ALTER TABLE events
    DROP COLUMN IF EXISTS dress_code,
    DROP COLUMN IF EXISTS category;
