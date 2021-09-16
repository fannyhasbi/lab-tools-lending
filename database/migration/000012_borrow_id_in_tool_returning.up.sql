BEGIN;

ALTER TABLE tool_returning ADD COLUMN IF NOT EXISTS borrow_id BIGINT;
-- populate
UPDATE tool_returning SET borrow_id = borrows.id
FROM borrows
WHERE tool_returning.user_id = borrows.user_id;
ALTER TABLE tool_returning ALTER COLUMN borrow_id SET NOT NULL;

ALTER TABLE tool_returning ADD FOREIGN KEY (borrow_id) REFERENCES borrows(id);
CREATE INDEX IF NOT EXISTS tool_returning_borrowid_idx ON tool_returning ("borrow_id");


ALTER TABLE tool_returning DROP COLUMN IF EXISTS user_id;
ALTER TABLE tool_returning DROP COLUMN IF EXISTS tool_id;

COMMIT;