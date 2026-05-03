ALTER TABLE collection_contributions
  ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'not_paid'
    CHECK (status IN ('not_paid', 'pending', 'paid')),
  ADD COLUMN IF NOT EXISTS receipt_url TEXT;

UPDATE collection_contributions SET status = 'paid' WHERE paid = TRUE;
