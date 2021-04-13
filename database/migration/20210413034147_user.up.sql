CREATE TABLE IF NOT EXISTS users (
  id SERIAL NOT NULL,
  chat_id INT NOT NULL,
  name VARCHAR(100) NOT NULL,
  nim VARCHAR(20),
  batch SMALLINT,
  address VARCHAR(500),
  created_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS users_id_idx ON users ("id");
CREATE INDEX IF NOT EXISTS users_chat_id_idx ON users ("chat_id");
