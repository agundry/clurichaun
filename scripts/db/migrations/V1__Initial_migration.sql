CREATE TABLE `monitor_categories` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  `type` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPRESSED;

CREATE TABLE `monitors` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  `symbol` varchar(30) NOT NULL,
  `monitor_category` int(11) NOT NULL,
  `created_epoch` int(11) NOT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPRESSED;

CREATE TABLE `records` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `monitor_id` bigint(20) unsigned NOT NULL,
  `record_epoch` int(11) NOT NULL,
  `value` decimal(10,2) NOT NULL,
  `updated` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_monitor_id_record_epoch` (`monitor_id`,`record_epoch`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPRESSED;
