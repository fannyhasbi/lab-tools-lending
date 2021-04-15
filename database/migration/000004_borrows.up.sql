CREATE TABLE IF NOT EXISTS borrows (
  id BIGSERIAL NOT NULL,
  amount INT NOT NULL,
  return_date DATE NOT NULL,
  user_id BIGINT NOT NULL,
  tool_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (tool_id) REFERENCES tools(id) 
);

CREATE INDEX IF NOT EXISTS users_id_idx ON users ("id");
CREATE INDEX IF NOT EXISTS users_userid_idx ON users ("user_id");