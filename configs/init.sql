CREATE USER ozon_root WITH password 'qwerty123';


DROP DATABASE IF EXISTS ozon_db;
CREATE DATABASE ozon_db
    WITH OWNER ozon_root
    ENCODING 'utf8';
GRANT ALL PRIVILEGES ON database ozon_db TO ozon_root;
\connect ozon_db;


DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
                       id SERIAL NOT NULL PRIMARY KEY,
                       first_name TEXT,
                       last_name TEXT,
                       email TEXT NOT NULL,
                       password BYTEA NOT NULL,
                       avatar TEXT
);


DROP TABLE IF EXISTS categories CASCADE;
CREATE TABLE categories (
                            id SERIAL NOT NULL PRIMARY KEY,
                            name TEXT NOT NULL
);



DROP TABLE IF EXISTS products CASCADE;
CREATE TABLE products (
                          id SERIAL NOT NULL PRIMARY KEY,
                          title TEXT NOT NULL UNIQUE,
                          rating NUMERIC(4, 2) NOT NULL,
                          description TEXT,
                          base_cost INTEGER NOT NULL,
                          discount INTEGER NOT NULL,
                          images TEXT[] NOT NULL,
                          id_category INTEGER NOT NULL,

                          FOREIGN KEY (id_category) REFERENCES categories(id),

                          CONSTRAINT rating_value CHECK (rating >= 0 AND rating <= 10),
                          CONSTRAINT discount_value CHECK (rating >= 0 AND rating <= 100)
);



DROP TABLE IF EXISTS subsets_category CASCADE;
CREATE TABLE subsets_category (
                                  id_category INTEGER NOT NULL,
                                  id_subset INTEGER NOT NULL,
                                  level INTEGER NOT NULL,

                                  FOREIGN KEY (id_category) REFERENCES categories (id),
                                  FOREIGN KEY (id_subset) REFERENCES categories (id),

                                  CONSTRAINT level_value CHECK (level > 0)
);

DROP TABLE IF EXISTS user_orders CASCADE;
CREATE TABLE user_orders (
                             id SERIAL NOT NULL PRIMARY KEY,
                             user_id INTEGER NOT NULL,
                             first_name TEXT NOT NULL,
                             last_name TEXT NOT NULL,
                             email TEXT NOT NULL,
                             address TEXT NOT NULL,
                             base_cost INTEGER NOT NULL,
                             total_cost INTEGER NOT NULL,
                             discount INTEGER NOT NULL,

                             FOREIGN KEY (user_id) REFERENCES users(id)
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


GRANT ALL PRIVILEGES ON TABLE users TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE ordered_products TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE categories TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE products TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE subsets_category TO ozon_root;
GRANT ALL PRIVILEGES ON TABLE user_orders TO ozon_root;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO ozon_root;


INSERT INTO categories(id, name) VALUES (1, 'Base');
INSERT INTO categories(id, name) VALUES (2, 'Электроника');
INSERT INTO categories(id, name) VALUES (3, 'Мобильные телефоны');
INSERT INTO categories(id, name) VALUES (4, 'Чехлы');
INSERT INTO categories(id, name) VALUES (7, 'Наушники');
INSERT INTO categories(id, name) VALUES (5, 'Для дома');
INSERT INTO categories(id, name) VALUES (6, 'Для уюта');


INSERT INTO subsets_category(id_category, id_subset, level) VALUES (1, 1, 1);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (2, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (2, 2, 2);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (3, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (3, 2, 2);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (3, 3, 3);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (7, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (7, 2, 2);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (7, 7, 3);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (4, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (4, 2, 2);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (4, 3, 3);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (4, 4, 4);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (5, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (5, 5, 2);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (6, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (6, 5, 2);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (6, 6, 3);

INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Чехол противоударный Armor Case для Samsung Galaxy A31, черный', 5, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 570, 21, '{"/product/6023600636.jpg", "/product/6023600623.jpg", "/product/6023600635.jpg", "/product/6023600630.jpg", "/product/6023600633.jpg", "/product/6023600625.jpg"}', 4);
INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Защитное стекло TORUS для Huawei Honor 10X Lite, закруглённные края, полное покрытие', 0, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 799, 75, '{"/product/6045097510.jpg", "/product/6045097505.jpg", "/product/6036662865.jpg", "/product/6036669672.jpg"}', 4);
INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Силиконовый чехол-накладка для iPhone 11/ Apple Silicone Case светло-голубой', 5, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 849, 50, '{"/product/6032890130.jpg", "/product/6032860336.jpg"}', 4);
INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Чехол-книжка MyPads для Meizu M5 Note прошитый по контуру с необычным геометрическим швом синий', 5, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 1600, 56, '{"/product/1037359526.jpg", "/product/1037359514.jpg", "/product/1037359516.jpg", "/product/1037359518.jpg", "/product/1037359520.jpg", "/product/1037359523.jpg"}', 4);
INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Аккумулятор для Apple iPhone 7 Plus', 0, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 1150, 36, '{"/product/6048563870.jpg"}', 4);
INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Защитное стекло для OPPO A5s', 3, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 280, 35, '{"/product/6036067021.jpg"}', 4);

INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Наушники для IPhone ligtning в футляре проводные для 7 8 X Xr Xs 11 (seria 81) , белые', 0, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 1990, 70, '{"/product/6046394320.jpg", "/product/6046394303.jpg", "/product/6046394306.jpg"}', 7);
INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Наушники беспроводные внутриканальные Defender OutFit B725, красный', 4, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 289, 42, '{"/product/1034408244.jpg", "/product/1034408253.jpg", "/product/6014331709.jpg"}', 7);
INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Наушники Sony MDR-EX15LP, белый', 5, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 699, 32, '{"/product/6011283025.jpg", "/product/6011283027.jpg", "/product/6011283026.jpg"}', 7);
INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Наушники Sony MDR-XD150, белый', 4, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 1236, 0, '{"/product/1007465889.jpg", "/product/1007465888.jpg", "/product/1007465887.jpg", "/product/1007465886.jpg"}', 7);

INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Фоторамка Veld Co "10*15", 1 фото', 4, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 210, 23, '{"/product/6026860610.jpg", "/product/6026825797.jpg"}', 6);
INSERT INTO products(title, rating, description, base_cost, discount, images, id_category) VALUES ('Фотоальбом Fotografia, 200 фото, 10 x 15 см (4 x 6")', 5, 'Насос предназначен для использования на гибридных велосипедах. Также он подходит для подкачивания колес городских, горных, BMX и детских велосипедов. Встроенный шлаг, который точно не потеряется! Ручка : 100.0% Полиамид 6.6', 639, 0, '{"/product/1036393602.jpg"}', 6);