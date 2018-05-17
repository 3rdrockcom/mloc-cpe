-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `transactions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) NOT NULL,
  `description` varchar(255) NOT NULL,
  `credit` decimal(10,2) NOT NULL,
  `debit` decimal(10,2) NOT NULL,
  `running_balance` decimal(10,2) NOT NULL,
  `balance` decimal(10,2) NOT NULL,
  `datetime` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `customer_id` (`customer_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `transactions`;