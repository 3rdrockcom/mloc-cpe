-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `api_keys` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) NOT NULL,
  `key` varchar(80) NOT NULL,
  `date_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `api_keys` (`id`, `customer_id`, `key`, `date_created`) VALUES
(1, 0, 'LOGIN', NOW()),
(2, 0, 'bd95200a60e47be9736970fd665f6195', NOW());

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `api_keys`;