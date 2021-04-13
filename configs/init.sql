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


-- Data for testing
INSERT INTO categories(name)VALUES ('Base');
INSERT INTO categories(name)VALUES ('Home');
INSERT INTO categories(name)VALUES ('Kitchen');
INSERT INTO categories(name)VALUES ('Dishes');
INSERT INTO categories(name)VALUES ('Electronics');
INSERT INTO categories(name)VALUES ('Mixer');

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (1, 1, 1);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (2, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (2, 2, 2);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (3, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (3, 2, 2);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (3, 3, 3);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (4, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (4, 2, 2);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (4, 4, 3);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (5, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (5, 2, 2);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (5, 5, 3);

INSERT INTO subsets_category(id_category, id_subset, level) VALUES (6, 1, 1);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (6, 2, 2);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (6, 5, 3);
INSERT INTO subsets_category(id_category, id_subset, level) VALUES (6, 6, 4);

INSERT INTO products(title, rating, description, base_cost, discount, images, id_category)
VALUES (
           'Hair dryer brush Rowenta',
           4,
           'The rotating Brush Activ airstyler provides ' ||
           'unsurpassed drying results. Power of 1000 ' ||
           'W guarantees fast drying effortlessly, two ' ||
           'rotating brushes with a diameter of 50 or 40 mm provide ' ||
           'professional styling. Ion generator and ' ||
           'ceramic coating smoothes hair, leaving it soft ' ||
           'and more brilliant.',
           20,
           20,
           '{"/product/1021166584.jpg", "/product/1021166585.jpg",
           "/product/1021166586.jpg", "/product/6043447767.jpg"}',
           6
       );

INSERT INTO products(title, rating, description, base_cost, discount, images, id_category)
VALUES (
           'Chupa Chups assorted caramel',
           3,
           'Chupa Chups Mini is Chupa Chups'' favorite candy on a stick ' ||
           'in mini format. In the showbox there are 100 Chupa. ' ||
           'Chups with the best flavors: strawberry, cola, orange, apple.',
           6,
           0,
           '{"/product/6024670802.jpg", "/product/6024670803.jpg",
           "/product/6024670804.jpg", "/product/6024670805.jpg"}',
           5
       );

INSERT INTO products(title, rating, description, base_cost, discount, images, id_category)
VALUES (
           'Electric Toothbrush Oral-B PRO 6000',
           4,
           'Oral-B is the # 1 brand of toothbrushes recommended ' ||
           'by most dentists in the world! Discover the Oral-B PRO 6000. ' ||
           'Smart Series Triumph! The Oral-B PRO 6000 Smart Series Triumph Toothbrush ' ||
           'features Bluetooth 4.0 to sync with the free Oral-B App. ' ||
           'Take your brushing to the next level in 2 minutes as ' ||
           'recommended by your dentist for superior cleansing and gum health.',
           50,
           5,
           '{"/product/6023124975.jpg", "/product/6023125065.jpg",
           "/product/6023125066.jpg"}',
           4
       );

INSERT INTO products(title, rating, description, base_cost, discount, images, id_category)
VALUES (
           'Bosch VitaPower Serie 4 jug blender',
           5,
           'Pro-Performance System: Optimal texture for smoothies, ' ||
           'even with frozen fruits and hard ingredients. Reliable Bosch motor: ' ||
           '1200 W power with engine speeds up to 30,000 rpm. Easy assembly and ' ||
           'cleaning: knife, bowl, lid are dishwasher safe, easy to install bowl ' ||
           'and store cable. High-quality assembly in our own Bosch factory: ' ||
           'durability, reliable construction and high-quality materials. ' ||
           'ProEdge Blades: Brilliant mixing results thanks to efficient ' ||
           'and durable stainless steel blades. Made in Germany.",',
           30,
           0,
           '{"/product/6026466446.jpg", "/product/6043224204.jpg",
           "/product/6043224631.jpg"}',
           4
       );
