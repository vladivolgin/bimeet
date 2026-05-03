ALTER TABLE events
    DROP COLUMN IF EXISTS reminder_3d_sent,
    DROP COLUMN IF EXISTS reminder_1d_sent;
