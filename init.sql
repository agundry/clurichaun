CREATE TABLE `records` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `monitor_id` bigint(20) unsigned NOT NULL,
  `record_epoch` datetime(3) NOT NULL,
  `value` datetime(3) NOT NULL,
  `updated` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_monitor_id_record_epoch` (`monitor_id`,`record_epoch`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPRESSED;
