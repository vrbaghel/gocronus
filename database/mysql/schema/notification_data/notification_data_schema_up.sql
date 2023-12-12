CREATE TABLE `notification_data` (
  id INT NOT NULL UNIQUE,
  nd_uuid INT NOT NULL,
  nd_title TEXT,
  nd_body TEXT,
  nd_source INT NOT NULL DEFAULT 0,
  nd_category ENUM('text', 'carousel', 'thumbnail_image', 'gif') NOT NULL,
  nd_navtype ENUM(
    'ai_tool',
    'ai_filter',
    'ai_photo',
    'profile',
    'pack_detail'
  ) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`id`) REFERENCES `notification` (`id`)
);