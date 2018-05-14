-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `customers` (
  `id` int(11) NOT NULL,
  `first_name` varchar(80) NOT NULL,
  `last_name` varchar(80) NOT NULL,
  `email` varchar(255) NOT NULL,
  `program_id` int(11) DEFAULT '1',
  `program_customer_id` int(11) DEFAULT NULL,
  `program_customer_mobile` varchar(50) DEFAULT NULL,
  `cust_unique_id` varchar(100) DEFAULT NULL,
  `last_transaction_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `customers`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `customers`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `customers`;