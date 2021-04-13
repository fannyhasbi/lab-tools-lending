CREATE TABLE IF NOT EXISTS chat_sessions (
  id SERIAL NOT NULL,
  status VARCHAR(50),
  user_id INT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS chat_session_details (
  id SERIAL NOT NULL,
  topic VARCHAR(50) NOT NULL,
  chat_session_id INT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (chat_session_id) REFERENCES chat_sessions(id)
);

CREATE INDEX IF NOT EXISTS chat_sessions_id_idx ON chat_sessions ("id");
CREATE INDEX IF NOT EXISTS chat_sessions_user_id_idx ON chat_sessions ("user_id");

CREATE INDEX IF NOT EXISTS chat_session_details_chat_session_id_idx ON chat_session_details ("chat_session_id");