SET NAMES utf8mb4;
CREATE DATABASE IF NOT EXISTS `music_box` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `music_box`;

-- ==========================================
-- 0. 卸载旧表 (按外键依赖的逆序删除)
-- ==========================================
DROP TABLE IF EXISTS `playlist_songs`;
DROP TABLE IF EXISTS `play_histories`;
DROP TABLE IF EXISTS `song_likes`;
DROP TABLE IF EXISTS `playlists`;
DROP TABLE IF EXISTS `songs`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `system_settings`;

-- ==========================================
-- 1. 创建表结构
-- ==========================================

-- 1. 用户表 (users)
CREATE TABLE `users` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
  `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
  `password_hash` VARCHAR(256) NOT NULL COMMENT '密码哈希',
  `avatar_url` VARCHAR(200) COMMENT '头像URL',
  `role` VARCHAR(20) NOT NULL DEFAULT 'User' COMMENT '角色',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-禁用, 1-正常',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. 歌曲表 (songs)
CREATE TABLE `songs` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '歌曲ID',
  `title` VARCHAR(200) NOT NULL,
  `artist` VARCHAR(100) NOT NULL,
  `album` VARCHAR(100),
  `duration` INT COMMENT '时长(秒)',
  `file_url` VARCHAR(300) NOT NULL COMMENT '文件URL',
  `cover_url` VARCHAR(300) COMMENT '封面URL',
  `lyric_url` VARCHAR(300) COMMENT '歌词URL',
  `upload_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. 点赞表 (song_likes)
CREATE TABLE `song_likes` (
  `song_id` INT NOT NULL COMMENT '歌曲外键',
  `user_id` INT NOT NULL COMMENT '用户外键',
  `liked_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`song_id`, `user_id`),
  CONSTRAINT `fk_likes_song` FOREIGN KEY (`song_id`) REFERENCES `songs` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. 歌单表 (playlists)
CREATE TABLE `playlists` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '歌单ID',
  `user_id` INT NOT NULL COMMENT '创建者ID',
  `title` VARCHAR(100) NOT NULL COMMENT '歌单标题',
  `description` VARCHAR(500) COMMENT '歌单描述',
  `cover_url` VARCHAR(300) COMMENT '歌单封面',
  `is_public` BOOLEAN DEFAULT TRUE COMMENT '是否公开',
  `play_count` INT DEFAULT 0 COMMENT '播放次数/热度',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-待审核, 1-通过, 2-驳回',
  `reject_reason` VARCHAR(255) DEFAULT '' COMMENT '驳回理由',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT `fk_playlist_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. 播放历史表 (play_histories)
CREATE TABLE `play_histories` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '播放历史ID',
  `user_id` INT NOT NULL COMMENT '用户ID',
  `song_id` INT NOT NULL COMMENT '歌曲ID',
  `played_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '播放时间',
  `duration` INT DEFAULT 0 COMMENT '本次播放时长(秒)',
  `source` VARCHAR(100) DEFAULT '' COMMENT '播放来源，如每日推荐/歌单/我喜欢的音乐',
  CONSTRAINT `fk_play_history_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_play_history_song` FOREIGN KEY (`song_id`) REFERENCES `songs` (`id`) ON DELETE CASCADE,
  KEY `idx_play_histories_user_played_at` (`user_id`, `played_at`),
  KEY `idx_play_histories_song_id` (`song_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户播放历史表';

-- 6. 歌单歌曲关联表 (playlist_songs)
CREATE TABLE `playlist_songs` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `playlist_id` INT NOT NULL,
  `song_id` INT NOT NULL,
  `sort_order` INT DEFAULT 0 COMMENT '排序权重',
  `added_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT `fk_ps_playlist` FOREIGN KEY (`playlist_id`) REFERENCES `playlists` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_ps_song` FOREIGN KEY (`song_id`) REFERENCES `songs` (`id`) ON DELETE CASCADE,
  UNIQUE KEY `uk_playlist_song` (`playlist_id`, `song_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 7. 系统设置表 (system_settings)
CREATE TABLE `system_settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `setting_key` varchar(255) NOT NULL COMMENT '设置键',
  `setting_value` text NOT NULL COMMENT '设置值',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_setting_key` (`setting_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ==========================================
-- 2. 初始化最小管理员与系统设置
-- ==========================================

INSERT INTO `users` (`id`, `username`, `password_hash`, `avatar_url`, `role`, `status`) VALUES
(1, 'admin', '$2a$10$emZD53.StcoycaaBcIXGg.CcPKQMUl0d5FCHCVznV8xMq.ikHPNJe', 'uploads/avatars/1.jpg', 'Admin', 1)
ON DUPLICATE KEY UPDATE
  `avatar_url` = VALUES(`avatar_url`),
  `role` = VALUES(`role`),
  `status` = VALUES(`status`);

INSERT INTO `system_settings` (`setting_key`, `setting_value`, `description`) VALUES
('allow_register', 'true', '是否开放用户注册')
ON DUPLICATE KEY UPDATE
  `setting_value` = VALUES(`setting_value`),
  `description` = VALUES(`description`);
