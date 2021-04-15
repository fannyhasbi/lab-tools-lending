CREATE TABLE IF NOT EXISTS tools (
  id BIGSERIAL NOT NULL,
  name VARCHAR(300) NOT NULL,
  brand VARCHAR(100),
  product_type VARCHAR(100),
  weight FLOAT,
  stock INT NOT NULL DEFAULT 0,
  additional_info TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS tools_id_idx ON tools ("id");
CREATE INDEX IF NOT EXISTS tools_stock_idx ON tools ("stock");