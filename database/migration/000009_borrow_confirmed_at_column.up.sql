ALTER TABLE borrows ADD COLUMN IF NOT EXISTS duration INT NOT NULL DEFAULT 7;
ALTER TABLE borrows ADD COLUMN IF NOT EXISTS confirmed_at TIMESTAMP;
ALTER TABLE borrows DROP COLUMN IF EXISTS return_date;