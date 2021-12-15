CREATE DATABASE `mfano` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for mfano_group
-- ----------------------------
DROP TABLE IF EXISTS `mfano_group`;
CREATE TABLE `mfano_group` (
                               `id` int(11) NOT NULL AUTO_INCREMENT,
                               `name` varchar(64) NOT NULL,
                               `created` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
                               `creator` varchar(64) NOT NULL,
                               `creator_id` int(11) NOT NULL,
                               PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for mfano_project
-- ----------------------------
DROP TABLE IF EXISTS `mfano_project`;
CREATE TABLE `mfano_project` (
                                 `id` int(11) NOT NULL AUTO_INCREMENT,
                                 `name` varchar(64) NOT NULL,
                                 `created` datetime(6) NOT NULL,
                                 `url` varchar(255) NOT NULL,
                                 `creator_id` int(11) NOT NULL,
                                 `creator` varchar(64) NOT NULL,
                                 `group` varchar(64) NOT NULL,
                                 `group_id` int(11) NOT NULL,
                                 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for mfano_user
-- ----------------------------
DROP TABLE IF EXISTS `mfano_user`;
CREATE TABLE `mfano_user` (
                              `id` int(11) NOT NULL,
                              `name` varchar(64) NOT NULL,
                              `nickname` varchar(64) NOT NULL,
                              `password` varchar(255) NOT NULL,
                              `email` varchar(255) NOT NULL,
                              PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
