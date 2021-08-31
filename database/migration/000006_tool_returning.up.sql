CREATE TABLE IF NOT EXISTS tool_returning (
  id BIGSERIAL NOT NULL,
  user_id BIGINT NOT NULL,
  tool_id BIGINT NOT NULL,
  status VARCHAR(50) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  additional_info TEXT,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (tool_id) REFERENCES tools(id)
);

CREATE INDEX IF NOT EXISTS tool_returning_id_idx ON tool_returning ("id");
CREATE INDEX IF NOT EXISTS tool_returning_userid_idx ON tool_returning ("user_id");