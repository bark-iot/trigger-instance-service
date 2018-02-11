CREATE TABLE trigger_instances (
  id SERIAL PRIMARY KEY,
  trigger_id INTEGER,
  input_data jsonb,
  created_at TIMESTAMP
);