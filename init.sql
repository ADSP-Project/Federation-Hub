CREATE TABLE shops (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) UNIQUE,
  description VARCHAR(255),
  webhookURL VARCHAR(255),
  publicKey VARCHAR(1024)
);