CREATE TABLE if not exists comments (
  `ID` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `UserID` INT UNSIGNED NOT NULL,
  `ShopID` INT UNSIGNED NOT NULL,
  `Content` TEXT NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `FOREIGN KEY` (ShopID) REFERENCES coffee_shops(ID),
  `FOREIGN KEY` (UserID) REFERENCES users(ID),
);
