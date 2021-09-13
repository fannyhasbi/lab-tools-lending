CREATE TABLE IF NOT EXISTS tool_photos (
  id BIGSERIAL NOT NULL,
  tool_id BIGINT NOT NULL,
  file_id TEXT NOT NULL,
  file_unique_id TEXT NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (tool_id) REFERENCES tools(id)
);

CREATE INDEX IF NOT EXISTS tool_photos_id_idx ON tool_photos ("id");
CREATE INDEX IF NOT EXISTS tool_photos_toolid_idx ON tool_photos ("tool_id");