ALTER TABLE collection_contributions
  DROP COLUMN IF EXISTS status,
  DROP COLUMN IF EXISTS receipt_url;
