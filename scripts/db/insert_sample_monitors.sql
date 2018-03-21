INSERT INTO `monitors` (`id`, `name`, `symbol`, `monitor_category`, `created_epoch`, `enabled`)
VALUES
	(1, 'Bitcoin', 'BTC', 1, 1519610400, 1),
	(2, 'Ethereum', 'ETH', 1, 1521512140, 1),
	(3, 'Litecoin', 'LTC', 1, 1521512140, 1);

INSERT INTO `monitor_categories` (`id`, `name`, `type`)
VALUES
	(1, 'cryptocurrency', 1);
