BEGIN;

ALTER TABLE tool_returning ADD COLUMN IF NOT EXISTS user_id BIGINT;
-- populate
UPDATE tool_returning SET user_id = borrows.user_id
FROM borrows
WHERE tool_returning.borrow_id = borrows.id; 
ALTER TABLE tool_returning ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE tool_returning ADD FOREIGN KEY (user_id) REFERENCES users(id);
CREATE INDEX IF NOT EXISTS tool_returning_userid_idx ON tool_returning ("user_id");


ALTER TABLE tool_returning ADD COLUMN IF NOT EXISTS tool_id BIGINT;
-- populate
UPDATE tool_returning SET tool_id = borrows.tool_id
FROM borrows
WHERE tool_returning.borrow_id = borrow.id;
ALTER TABLE tool_returning ALTER COLUMN tool_id SET NOT NULL;
ALTER TABLE tool_returning ADD FOREIGN KEY (tool_id) REFERENCES tools(id);


ALTER TABLE tool_returning DROP COLUMN IF EXISTS borrow_id;

COMMIT;