CREATE USER ozon_root WITH password 'qwerty123';


DROP DATABASE IF EXISTS ozon_auth_db;
CREATE DATABASE ozon_auth_db
    WITH OWNER ozon_root
    ENCODING 'utf8';
GRANT ALL PRIVILEGES ON database ozon_auth_db TO ozon_root;
\connect ozon_auth_db;

DROP TABLE IF EXISTS auth_users CASCADE;
CREATE TABLE auth_users (
    id SERIAL NOT NULL PRIMARY KEY,
    email TEXT NOT NULL,
    password BYTEA NOT NULL
);

GRANT ALL PRIVILEGES ON TABLE auth_users TO ozon_root;
