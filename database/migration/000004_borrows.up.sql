CREATE TABLE IF NOT EXISTS borrows (
  id BIGSERIAL NOT NULL,
  amount INT NOT NULL DEFAULT 1,
  return_date DATE,
  status VARCHAR(50) NOT NULL,
  user_id BIGINT NOT NULL,
  tool_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (tool_id) REFERENCES tools(id) 
);

CREATE INDEX IF NOT EXISTS borrows_id_idx ON borrows ("id");
CREATE INDEX IF NOT EXISTS borrows_userid_idx ON borrows ("user_id");
CREATE INDEX IF NOT EXISTS borrows_status_idx ON borrows ("status");