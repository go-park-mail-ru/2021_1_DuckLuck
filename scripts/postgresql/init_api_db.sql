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

    CONSTRAINT discount_value CHECK (discount >= 0 AND discount <= 100),
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
    status TEXT NOT NULL DEFAULT 'оформлен',

    FOREIGN KEY (user_id) REFERENCES data_users(id),

    CONSTRAINT discount_value CHECK (discount >= 0),
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

    CONSTRAINT rating_value CHECK (rating >= 1 AND rating <= 5)
);

DROP TABLE IF EXISTS favorites CASCADE;
CREATE TABLE favorites (
    id SERIAL NOT NULL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    date_added TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (user_id) REFERENCES data_users(id),

    CONSTRAINT uniq_like UNIQUE (product_id, user_id)
);

GRANT ALL PRIVILEGES ON TABLE data_users TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE promo_codes TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE ordered_products TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE categories TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE products TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE user_orders TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE reviews TO ozon_root;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO ozon_root;


INSERT INTO promo_codes(id, code, sale) VALUES (1, 'test1', 10);
INSERT INTO promo_codes(id, code, sale) VALUES (2, 'test2', 5);
INSERT INTO promo_codes(id, code, sale) VALUES (3, 'test3', 100);

INSERT INTO categories(id, name, left_node, right_node, level) VALUES (1, 'Base', 0, 13, 0);
INSERT INTO categories(id, name, left_node, right_node, level) VALUES (2, 'Электроника', 1, 8, 1);
INSERT INTO categories(id, name, left_node, right_node, level) VALUES (3, 'Мобильные телефоны', 2, 5, 2);
INSERT INTO categories(id, name, left_node, right_node, level) VALUES (4, 'Чехлы', 3, 4, 3);
INSERT INTO categories(id, name, left_node, right_node, level) VALUES (7, 'Наушники', 6, 7, 2);
INSERT INTO categories(id, name, left_node, right_node, level) VALUES (5, 'Для дома', 9, 12, 1);
INSERT INTO categories(id, name, left_node, right_node, level) VALUES (6, 'Для уюта', 10, 11, 2);

INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties) VALUES ('Чехол противоударный Armor Case для Samsung Galaxy A31, черный', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 570, 451, 21, '{"/product/6023600636.webp", "/product/6023600623.webp", "/product/6023600635.webp", "/product/6023600630.webp", "/product/6023600633.webp", "/product/6023600625.webp"}', 4, '{ "вес" : "10 кг", "цвет" : "чёрный" }');
INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties, sale_group) VALUES ('Защитное стекло TORUS для Huawei Honor 10X Lite, закруглённные края, полное покрытие', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 799, 200, 75, '{"/product/6045097510.webp", "/product/6045097505.webp", "/product/6036662865.webp", "/product/6036669672.webp"}', 4, '{ "вес" : "10 кг", "цвет" : "чёрный" }', '{1, 2}');
INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties) VALUES ('Силиконовый чехол-накладка для iPhone 11/ Apple Silicone Case светло-голубой', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 849, 425, 50, '{"/product/6032890130.webp", "/product/6032860336.webp"}', 4, '{ "вес" : "10 кг", "цвет" : "чёрный" }');
INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties, sale_group) VALUES ('Чехол-книжка MyPads для Meizu M5 Note прошитый по контуру с необычным геометрическим швом синий', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 1600, 704, 56, '{"/product/1037359526.webp", "/product/1037359514.webp", "/product/1037359516.webp", "/product/1037359518.webp", "/product/1037359520.webp", "/product/1037359523.webp"}', 4, '{ "вес" : "10 кг", "цвет" : "чёрный" }', '{2}');
INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties) VALUES ('Аккумулятор для Apple iPhone 7 Plus', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 1150, 736, 36, '{"/product/6048563870.webp"}', 4, '{ "вес" : "10 кг", "цвет" : "чёрный" }');
INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties, sale_group) VALUES ('Защитное стекло для OPPO A5s', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 280, 182, 35, '{"/product/6036067021.webp"}', 4, '{ "вес" : "10 кг", "цвет" : "чёрный" }', '{1}');

INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties, sale_group) VALUES ('Наушники для IPhone ligtning в футляре проводные для 7 8 X Xr Xs 11 (seria 81) , белые', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 1990, 597, 70, '{"/product/6046394320.webp", "/product/6046394303.webp", "/product/6046394306.webp"}', 7, '{ "вес" : "10 кг", "цвет" : "чёрный" }', '{1, 2}');
INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties) VALUES ('Наушники беспроводные внутриканальные Defender OutFit B725, красный', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 289, 168, 42, '{"/product/1034408244.webp", "/product/1034408253.webp", "/product/6014331709.webp"}', 7, '{ "вес" : "10 кг", "цвет" : "чёрный" }');
INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties, sale_group) VALUES ('Наушники Sony MDR-EX15LP, белый', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 699, 476, 32, '{"/product/6011283025.webp", "/product/6011283027.webp", "/product/6011283026.webp"}', 7, '{ "вес" : "10 кг", "цвет" : "чёрный" }', '{3}');
INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties) VALUES ('Наушники Sony MDR-XD150, белый', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 1236, 1236, 0, '{"/product/1007465889.webp", "/product/1007465888.webp", "/product/1007465887.webp", "/product/1007465886.webp"}', 7, '{ "вес" : "10 кг", "цвет" : "чёрный" }');

INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties) VALUES ('Фоторамка Veld Co "10*15", 1 фото', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 210, 162, 23, '{"/product/6026860610.webp", "/product/6026825797.webp"}', 6, '{ "вес" : "10 кг", "цвет" : "чёрный" }');
INSERT INTO products(title, description, base_cost, total_cost, discount, images, id_category, properties, sale_group) VALUES ('Фотоальбом Fotografia, 200 фото, 10 x 15 см (4 x 6")', 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 639, 639, 0, '{"/product/1036393602.webp"}', 6, '{ "вес" : "10 кг", "цвет" : "чёрный" }', '{2}');