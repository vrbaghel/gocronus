CREATE TABLE `notification_pack` (
  id INT NOT NULL UNIQUE,
  np_id VARCHAR(16),
  np_order_id VARCHAR(16),
  np_filter_id VARCHAR(16),
  np_tool_id VARCHAR(16),
  np_name VARCHAR(50),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`id`) REFERENCES `notification_data` (`id`)
);