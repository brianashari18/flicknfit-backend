-- Rename shopping_cart_id to saved_items_id in shopping_cart_items table
ALTER TABLE `shopping_cart_items` 
CHANGE COLUMN `shopping_cart_id` `saved_items_id` BIGINT UNSIGNED NOT NULL;

-- Update foreign key constraint name (optional, for consistency)
ALTER TABLE `shopping_cart_items` 
DROP FOREIGN KEY `fk_shopping_carts_shopping_cart_items`;

ALTER TABLE `shopping_cart_items` 
ADD CONSTRAINT `fk_shopping_carts_saved_items_list` 
FOREIGN KEY (`saved_items_id`) REFERENCES `shopping_carts`(`id`) 
ON DELETE CASCADE ON UPDATE CASCADE;
