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
    firstName TEXT,
    lastName TEXT,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    avatar TEXT
);


DROP TABLE IF EXISTS category CASCADE;
CREATE TABLE category (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);


DROP TABLE IF EXISTS products CASCADE;
CREATE TABLE products (
    id SERIAL NOT NULL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    rating NUMERIC(4, 2) NOT NULL,
    description TEXT,
    baseCost INTEGER NOT NULL,
    discount INTEGER NOT NULL,
    images TEXT[] NOT NULL,
    idCategory INTEGER NOT NULL,

    FOREIGN KEY (idCategory) REFERENCES category(id),

    CONSTRAINT ratingValue CHECK (rating >= 0 AND rating <= 10),
    CONSTRAINT discountValue CHECK (rating >= 0 AND rating <= 100)
);


DROP TABLE IF EXISTS subsetCategory CASCADE;
CREATE TABLE subsetCategory (
    idCategory INTEGER NOT NULL,
    idSubSet INTEGER NOT NULL,
    level INTEGER NOT NULL,

    FOREIGN KEY (idCategory) REFERENCES category(id),
    FOREIGN KEY (idSubSet) REFERENCES category(id),

    CONSTRAINT levelValue CHECK (level > 0)
);


-- Data for testing
INSERT INTO category(name) VALUES ('Home');
INSERT INTO category(name) VALUES ('Kitchen');
INSERT INTO category(name) VALUES ('Dishes');
INSERT INTO category(name) VALUES ('Electronics');
INSERT INTO category(name) VALUES ('Food');

INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (1, 1, 1);

INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (2, 1, 1);
INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (2, 2, 2);

INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (3, 1, 1);
INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (3, 2, 2);
INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (3, 3, 3);

INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (4, 1, 1);
INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (4, 2, 2);
INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (4, 4, 3);

INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (5, 1, 1);
INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (5, 2, 2);
INSERT INTO subsetCategory(idCategory, idSubSet, level) VALUES (5, 5, 3);

INSERT INTO products(title, rating, description, baseCost, discount, images, idCategory)
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
            4
       );

INSERT INTO products(title, rating, description, baseCost, discount, images, idCategory)
VALUES (
            'Chupa Chups assorted caramel',
            3,
            'Chupa Chups Mini is Chupa Chups'' favorite candy on a stick ' ||
            'in mini format. In the showbox there are 100 Chupa. ' ||
            'Chups with the best flavors: strawberry, cola, orange, apple.',
            6.25,
            0,
            '{"/product/6024670802.jpg", "/product/6024670803.jpg",
            "/product/6024670804.jpg", "/product/6024670805.jpg"}',
            5
       );

INSERT INTO products(title, rating, description, baseCost, discount, images, idCategory)
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

INSERT INTO products(title, rating, description, baseCost, discount, images, idCategory)
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


CREATE OR REPLACE FUNCTION getPathOfCategory(INT)
RETURNS TEXT[] AS $$
select array(
            SELECT c.name FROM  subsetCategory s
            LEFT JOIN category c ON c.id = s.idSubset
            WHERE s.idCategory = $1
            ORDER BY s.level
        )
$$
LANGUAGE SQL;
