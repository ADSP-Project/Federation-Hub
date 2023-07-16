CREATE TABLE shops (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) UNIQUE,
  description VARCHAR(255),
  webhookURL VARCHAR(255),
  publicKey VARCHAR(1024)
);

INSERT INTO shops(id, name, description, webhookurl, publickey) VALUES
(1, 'Tech Mart', 'Your one stop for all tech gadgets', 'http://localhost:8001/webhook', 'PublicKey1'),
(2, 'Garden Central', 'Everything you need for your garden', 'http://localhost:8002/webhook', 'PublicKey2'),
(3, 'Sports Gear Galore', 'Sports equipment for all ages', 'http://localhost:8003/webhook', 'PublicKey3'),
(4, 'Fashion Boutique', 'Latest fashion trends for you', 'http://localhost:8004/webhook', 'PublicKey4'),
(5, 'Pet Paradise', 'Pet food and accessories', 'http://localhost:8005/webhook', 'PublicKey5'),
(6, 'Home Decor Hub', 'Make your home beautiful', 'http://localhost:8006/webhook', 'PublicKey6'),
(7, 'Beauty Bliss', 'Beauty products for everyone', 'http://localhost:8007/webhook', 'PublicKey7'),
(8, 'Fitness Fanatics', 'Gym equipment and sportswear', 'http://localhost:8008/webhook', 'PublicKey8'),
(9, 'Kids Kingdom', 'Toys and clothes for kids', 'http://localhost:8009/webhook', 'PublicKey9'),
(10, 'Auto Accessories', 'Accessories for your vehicle', 'http://localhost:8010/webhook', 'PublicKey10'),
(11, 'Healthy Harvest', 'Organic produce for your home', 'http://localhost:8011/webhook', 'PublicKey11'),
(12, 'Book Barn', 'Books for all genres', 'http://localhost:8012/webhook', 'PublicKey12'),
(13, 'Music Mania', 'Instruments and music gear', 'http://localhost:8013/webhook', 'PublicKey13'),
(14, 'Travel Treasures', 'Travel gear for your adventures', 'http://localhost:8014/webhook', 'PublicKey14'),
(15, 'Artistic Alley', 'Art supplies and crafts', 'http://localhost:8015/webhook', 'PublicKey15'),
(16, 'Outdoor Outfitters', 'Gear for your outdoor activities', 'http://localhost:8016/webhook', 'PublicKey16'),
(17, 'Gourmet Grocers', 'Fine foods and ingredients', 'http://localhost:8017/webhook', 'PublicKey17'),
(18, 'Stationery Stop', 'Office supplies and stationery', 'http://localhost:8018/webhook', 'PublicKey18'),
(19, 'Tool Town', 'Tools for your DIY projects', 'http://localhost:8019/webhook', 'PublicKey19'),
(20, 'Luxury Linens', 'High-end linens for your home', 'http://localhost:8020/webhook', 'PublicKey20');
