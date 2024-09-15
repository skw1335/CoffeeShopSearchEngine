create TABLE if not exists ratings (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `UserID` INT NOT NULL,
  `ShopID` INT NOT NULL,
  `Ambiance` VARCHAR(255),
  `Coffee`   INT,
  `Overall`  VARCHAR(255),
  `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (ShopID) REFERENCES coffee_shops(ID),
  FOREIGN KEY (UserID) REFERENCES users(id),
);
