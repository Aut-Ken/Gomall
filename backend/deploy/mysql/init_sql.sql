-- Add product_image column to orders table
-- ALTER TABLE `orders`
-- ADD COLUMN `product_image` VARCHAR(500) DEFAULT NULL COMMENT '商品图片快照' AFTER `product_name`;

-- Note: The column `product_image` might have been created by GORM AutoMigrate.
-- If it doesn't exist, uncomment the line above.

-- Update existing product images for new items
-- iPad Air
UPDATE `products`
SET
    `image_url` = 'http://localhost:8080/photos/product_ipad.jpg'
WHERE
    `name` LIKE '%iPad Air%';

-- MacBook Pro 14
UPDATE `products`
SET
    `image_url` = 'http://localhost:8080/photos/product_macbook.jpg'
WHERE
    `name` LIKE '%MacBook Pro 14%';