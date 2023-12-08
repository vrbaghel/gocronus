CREATE TABLE `notification` (
  id INT NOT NULL UNIQUE AUTO_INCREMENT,
  n_action VARCHAR(100) NOT NULL,
  n_timezone ENUM('IST', 'GMT') NOT NULL,
  n_timestamp TIMESTAMP NOT NULL,
  n_device ENUM('ios', 'android') NOT NULL,
  n_status ENUM(
    'scheduled',
    'running',
    'completed',
    'terminated'
  ) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);