ALTER TABLE comments
ADD CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id);

ALTER TABLE comments
ADD CONSTRAINT fk_shop FOREIGN KEY(shop_id) REFERENCES coffee_shops(id);
