/*
 Navicat MySQL Data Transfer

 Source Server         : conntest
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : localhost:3306
 Source Schema         : user_system

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 20/04/2022 00:30:21
*/

SET NAMES utf8mb4;
DROP TABLE IF EXISTS `users_0`;
CREATE TABLE `users_0` (
     `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
     `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
     `password` varchar(32) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL COMMENT '密码',
     `salt` varchar(14) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL COMMENT '密码盐值',
     `nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '昵称',
     `pic_path` varchar(255) CHARACTER SET ascii COLLATE ascii_general_ci DEFAULT NULL COMMENT '图片路径',
     `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
     `version` int unsigned NOT NULL DEFAULT '1' COMMENT '记录版本号',
     PRIMARY KEY (`id`),
     UNIQUE KEY `uniq_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `users_1`;
DROP TABLE IF EXISTS `users_2`;
DROP TABLE IF EXISTS `users_3`;
DROP TABLE IF EXISTS `users_4`;
DROP TABLE IF EXISTS `users_5`;
DROP TABLE IF EXISTS `users_6`;
DROP TABLE IF EXISTS `users_7`;
DROP TABLE IF EXISTS `users_8`;
DROP TABLE IF EXISTS `users_9`;
CREATE TABLE users_1 LIKE users_0;
CREATE TABLE users_2 LIKE users_0;
CREATE TABLE users_3 LIKE users_0;
CREATE TABLE users_4 LIKE users_0;
CREATE TABLE users_5 LIKE users_0;
CREATE TABLE users_6 LIKE users_0;
CREATE TABLE users_7 LIKE users_0;
CREATE TABLE users_8 LIKE users_0;
CREATE TABLE users_9 LIKE users_0;
