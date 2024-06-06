CREATE DATABASE greenlight;

\c greenlight;

CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';

CREATE EXTENSION IF NOT EXISTS citext;

DATABASE_URL="postgres://greenlight:pa55word@localhost:5432/greenlight?sslmode=disable"

GRANT ALL ON DATABASE mydb TO admin;

ALTER DATABASE mydb OWNER TO admin;

\dt

SELECT * FROM schema_migrations;

\d movies

\d users

-- Set the activated field for alice@example.com to true.
UPDATE users SET activated = true WHERE email = 'alice@example.com';

-- Give all users the 'movies:read' permission
INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE code = 'movies:read') FROM users;

-- Give faith@example.com the 'movies:write' permission
INSERT INTO users_permissions
    VALUES (
    (SELECT id FROM users WHERE email = 'faith@example.com'),
    (SELECT id FROM permissions WHERE code = 'movies:write')
);

-- List all activated users and their permissions.
SELECT email, array_agg(permissions.code) as permissions
FROM permissions
INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
INNER JOIN users ON users_permissions.user_id = users.id
WHERE users.activated = true
GROUP BY email;

SELECT email, code FROM users
INNER JOIN users_permissions ON users.id = users_permissions.user_id
INNER JOIN permissions ON users_permissions.permission_id = permissions.id
WHERE users.email = 'grace@example.com';