START TRANSACTION;

CREATE TABLE `customers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(50) NOT NULL,
  `last_name` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `mobile_number` VARCHAR(20) NOT NULL,
  `program_id` int(11) DEFAULT '1',
  `program_customer_id` int(11) NOT NULL,
  `program_customer_mobile` varchar(50) NOT NULL,
  `cust_unique_id` varchar(100) NOT NULL,
  `last_transaction_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cust_unique_id` (`cust_unique_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

COMMIT;