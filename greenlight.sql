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