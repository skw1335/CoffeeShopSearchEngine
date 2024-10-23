CREATE TABLE IF NOT EXISTS roles(
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  level int NOT NULL DEFAULT 1,
  description TEXT
);

INSERT INTO
  roles (name, description, level)
VALUES (
  'user',
  'A user can create ratings and comments',
  3
);

INSERT INTO
  roles (name, description, level)
VALUES (
  'moderator',
  'A moderator can update others ratings and comments',
  2
);

INSERT INTO
  roles (name, description, level)
VALUES (
  'admin',
  'An admin has full control over all posts',
  3
);




