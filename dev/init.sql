CREATE TABLE IF NOT EXISTS package_dep (
  id SERIAL PRIMARY KEY,
  packageName TEXT NOT NULL UNIQUE,
  version TEXT NOT NULL,
  dependencies TEXT,
  updatedAt TEXT
);