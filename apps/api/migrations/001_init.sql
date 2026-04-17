CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  display_name VARCHAR(64) NOT NULL,
  role VARCHAR(32) NOT NULL DEFAULT 'admin',
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL
);

CREATE TABLE IF NOT EXISTS departments (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL UNIQUE,
  manager VARCHAR(64) NULL,
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL
);

INSERT INTO users (username, password_hash, display_name, role, status, created_at, updated_at)
VALUES ('admin', '$2a$10$QgI8ql6f8sXtVGzJfQ0f6.ULJLSm4Qo8v7NZaY1f3T5iAG2wE8VPa', '超级管理员', 'admin', 1, NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE updated_at = NOW(3);

