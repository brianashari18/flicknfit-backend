-- Migration: Add OAuth fields to users table
-- Date: 2025-11-20
-- Description: Add support for Google/Facebook OAuth authentication

-- Add new columns for OAuth support
ALTER TABLE `users` 
ADD COLUMN `auth_provider` ENUM('local', 'google', 'facebook') DEFAULT 'local' COMMENT 'Authentication provider' AFTER `role`,
ADD COLUMN `auth_provider_id` VARCHAR(255) NULL COMMENT 'OAuth provider user ID (Firebase UID)' AFTER `auth_provider`,
ADD COLUMN `profile_picture_url` VARCHAR(512) NULL COMMENT 'User profile picture URL from OAuth provider' AFTER `auth_provider_id`,
ADD COLUMN `is_email_verified` BOOLEAN DEFAULT FALSE COMMENT 'Email verification status' AFTER `profile_picture_url`,
ADD COLUMN `last_login` TIMESTAMP NULL COMMENT 'Last login timestamp' AFTER `is_email_verified`;

-- Modify password column to be nullable (for OAuth users without password)
ALTER TABLE `users` 
MODIFY COLUMN `password` VARCHAR(255) NULL DEFAULT NULL COMMENT 'Hashed password (NULL for OAuth users)';

-- Modify phone_number to be nullable (optional field)
ALTER TABLE `users`
MODIFY COLUMN `phone_number` VARCHAR(20) NULL DEFAULT NULL COMMENT 'User phone number (optional)';

-- Add index for faster OAuth lookups
CREATE INDEX `idx_auth_provider_id` ON `users`(`auth_provider_id`);
CREATE INDEX `idx_auth_provider` ON `users`(`auth_provider`);

-- Add unique constraint for OAuth provider ID (one account per provider)
ALTER TABLE `users`
ADD UNIQUE KEY `unique_auth_provider_user` (`auth_provider`, `auth_provider_id`);
