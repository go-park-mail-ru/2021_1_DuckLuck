CREATE USER ozon_root WITH password 'qwerty123';


DROP DATABASE IF EXISTS ozon_api_db;
CREATE DATABASE ozon_api_db
    WITH OWNER ozon_root
    ENCODING 'utf8';
GRANT ALL PRIVILEGES ON database ozon_api_db TO ozon_root;
\connect ozon_api_db;


CREATE TEXT SEARCH DICTIONARY russian_ispell (
    TEMPLATE = ispell,
    DictFile = russian,
    AffFile = russian,
    StopWords = russian
);

CREATE TEXT SEARCH CONFIGURATION ru (COPY=russian);

ALTER TEXT SEARCH CONFIGURATION ru
    ALTER MAPPING FOR hword, hword_part, word
    WITH russian_ispell, russian_stem;

SET default_text_search_config = 'ru';

DROP TABLE IF EXISTS data_users CASCADE;
CREATE TABLE data_users (
    id SERIAL NOT NULL PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    avatar TEXT,
    email TEXT NOT NULL
);


DROP TABLE IF EXISTS categories CASCADE;
CREATE TABLE categories (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    left_node INTEGER NOT NULL,
    right_node INTEGER NOT NULL,
    level INTEGER NOT NULL,

    CONSTRAINT left_value CHECK (left_node >= 0 AND left_node < right_node),
    CONSTRAINT right_value CHECK (right_node > left_node),
    CONSTRAINT level_value CHECK (level >= 0)
);

DROP TABLE IF EXISTS promo_codes CASCADE;
CREATE TABLE promo_codes (
    id SERIAL NOT NULL PRIMARY KEY,
    code TEXT NOT NULL,
    sale INTEGER NOT NULL
);

DROP TABLE IF EXISTS products CASCADE;
CREATE TABLE products (
    id SERIAL NOT NULL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    base_cost INTEGER NOT NULL,
    total_cost INTEGER NOT NULL,
    discount INTEGER NOT NULL,
    images TEXT[] NOT NULL,
    id_category INTEGER NOT NULL,
    date_added TIMESTAMP NOT NULL DEFAULT NOW(),
    properties JSONB NOT NULL,
    sale_group INTEGER[] NOT NULL DEFAULT '{}',
    fts TSVECTOR,

    FOREIGN KEY (id_category) REFERENCES categories(id),

    CONSTRAINT discount_value CHECK (discount >= 0 AND discount <= base_cost),
    CONSTRAINT base_cost_value CHECK (base_cost >= 0),
    CONSTRAINT total_cost_value CHECK (total_cost >= 0 AND total_cost <= base_cost)
);

DROP SEQUENCE IF EXISTS order_num CASCADE;
CREATE SEQUENCE order_num
    INCREMENT 1
    CACHE 20;

DROP SEQUENCE IF EXISTS order_serial CASCADE;
CREATE SEQUENCE order_serial
    START WITH 17
    INCREMENT 3
    CACHE 20;
    
DROP INDEX IF EXISTS products_fts CASCADE;
CREATE INDEX products_fts ON products USING GIN (fts);

DROP FUNCTION IF EXISTS create_fts CASCADE;
CREATE OR REPLACE FUNCTION create_fts() RETURNS TRIGGER AS $$
BEGIN
    NEW.fts = setweight(to_tsvector('ru', NEW.title), 'A')
        || setweight(to_tsvector('ru', NEW.description), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS t_create_fts on products CASCADE;
CREATE TRIGGER t_create_fts
    BEFORE INSERT ON products
    FOR EACH ROW
    EXECUTE PROCEDURE create_fts();



DROP TABLE IF EXISTS user_orders CASCADE;
CREATE TABLE user_orders (
    id SERIAL NOT NULL PRIMARY KEY,
    order_num TEXT UNIQUE NOT NULL
        DEFAULT lpad(text(nextval('order_num')), 8, '0')
        || '-' || lpad(text(nextval('order_serial')), 4, '0'),
    user_id INTEGER NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    address TEXT NOT NULL,
    base_cost INTEGER NOT NULL,
    total_cost INTEGER NOT NULL,
    discount INTEGER NOT NULL,
    date_added TIMESTAMP NOT NULL DEFAULT NOW(),
    date_delivery TIMESTAMP NOT NULL DEFAULT NOW(),
    status_pay TEXT NOT NULL DEFAULT 'оплачено',
    status_delivery TEXT NOT NULL DEFAULT 'получено',

    FOREIGN KEY (user_id) REFERENCES data_users(id),

    CONSTRAINT discount_value CHECK (discount >= 0 AND discount <= base_cost),
    CONSTRAINT base_cost_value CHECK (base_cost >= 0),
    CONSTRAINT total_cost_value CHECK (total_cost >= 0 AND total_cost <= base_cost)
);

DROP TABLE IF EXISTS ordered_products CASCADE;
CREATE TABLE ordered_products (
    id SERIAL NOT NULL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    order_id INTEGER NOT NULL,
    num INTEGER NOT NULL,
    base_cost INTEGER NOT NULL,
    discount INTEGER NOT NULL,

    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (order_id) REFERENCES user_orders(id),

    CONSTRAINT num_value CHECK (num >= 0)
);

DROP TABLE IF EXISTS reviews CASCADE;
CREATE TABLE reviews (
    id SERIAL NOT NULL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    rating INTEGER NOT NULL,
    advantages TEXT,
    disadvantages TEXT,
    comment TEXT,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    date_added TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (user_id) REFERENCES data_users(id),

    CONSTRAINT rating_value CHECK (rating >= 0 AND rating <= 5)
);

GRANT ALL PRIVILEGES ON TABLE data_users TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE promo_codes TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE ordered_products TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE categories TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE products TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE user_orders TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE reviews TO ozon_root;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO ozon_root;

COPY categories(id,name,right_node,left_node,level)
    FROM '/categories.csv'
    DELIMITER ','
    CSV HEADER;

COPY products(title,description,base_cost,total_cost,discount,images,id_category,properties)
    FROM '/products.csv'
    DELIMITER ','
    CSV HEADER;
