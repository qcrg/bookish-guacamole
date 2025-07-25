CREATE TABLE users (
  id UUID PRIMARY KEY
);
INSERT INTO users (id) VALUES ('01983781-0bab-7ad3-8c1d-dd713d2eadd8');
INSERT INTO users (id) VALUES ('01983784-e5ff-786a-bdd0-ab93f98362d2');
INSERT INTO users (id) VALUES ('01983784-f1f9-79b7-b58a-7895bb24d23b');
INSERT INTO users (id) VALUES ('01983787-435a-7688-be24-45bfd4883a9e');

CREATE TABLE tokens (
  id UUID PRIMARY KEY,
  refresh CHAR(60) NOT NULL
);

CREATE TABLE sessions (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) NOT NULL,
  token_id UUID REFERENCES tokens(id) NOT NULL,
  user_agent CHAR(88) NOT NULL,
  init_ip CHAR(88) NOT NULL,
  expires_at DATE
);
