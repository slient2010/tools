-- MySQL dump 10.13  Distrib 5.1.73, for redhat-linux-gnu (x86_64)
--
-- Host: localhost    Database: ljops
-- ------------------------------------------------------
-- Server version	5.1.73

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `netdeviceinfo`
--

DROP TABLE IF EXISTS `netdeviceinfo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `netdeviceinfo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `device_name` varchar(240) DEFAULT NULL COMMENT '设备名称',
  `device_ip` varchar(240) DEFAULT NULL COMMENT '设备IP',
  `device_type` varchar(240) DEFAULT NULL COMMENT '设备类型',
  `vendor` varchar(240) DEFAULT NULL COMMENT '设备供应商',
  `device_no` varchar(240) DEFAULT NULL COMMENT '设备型号',
  `username` varchar(240) DEFAULT NULL COMMENT '设备用户名',
  `password` varchar(240) DEFAULT NULL COMMENT '设备密码',
  `device_detail` text COMMENT '详细配置',
  `device_status` varchar(240) DEFAULT NULL COMMENT '设备使用状态',
  `device_location` varchar(240) DEFAULT NULL COMMENT '设备使用位置',
  `device_position` varchar(240) DEFAULT NULL COMMENT '设备位置',
  `created_time` datetime NOT NULL COMMENT '设备创建(购买)时间',
  `update_time` datetime NOT NULL COMMENT '设备信息修改时间',
  `notes` text COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `device_ip` (`device_ip`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `netdeviceinfo`
--

LOCK TABLES `netdeviceinfo` WRITE;
/*!40000 ALTER TABLE `netdeviceinfo` DISABLE KEYS */;
/*!40000 ALTER TABLE `netdeviceinfo` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `serverinfo`
--

DROP TABLE IF EXISTS `serverinfo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `serverinfo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `server_name` varchar(240) DEFAULT NULL COMMENT '服务器名称',
  `IP` varchar(240) DEFAULT NULL COMMENT '服务器IP',
  `manage_ip` varchar(240) DEFAULT NULL COMMENT '服务器bmc管理口IP',
  `vendor` varchar(240) DEFAULT NULL COMMENT '服务器供应商',
  `server_type` varchar(240) DEFAULT NULL COMMENT '服务器型号',
  `memory` int(11) DEFAULT NULL COMMENT '服务器内存G',
  `cpu` varchar(240) DEFAULT NULL COMMENT '服务器CPU',
  `interface` text COMMENT '服务器网口',
  `disk` varchar(240) DEFAULT NULL COMMENT '服务器硬盘',
  `server_detail` text COMMENT '详细配置',
  `server_system` varchar(240) DEFAULT NULL COMMENT '服务器系统',
  `os_type` int(11) NOT NULL COMMENT '服务器类型虚拟机or物理机,0-物理机, 1-虚拟机',
  `os_type_id` int(11) NOT NULL COMMENT '服务器类型虚拟机or物理机id从属关系',
  `status` varchar(240) DEFAULT NULL COMMENT '服务器使用状态',
  `position` varchar(240) DEFAULT NULL COMMENT '服务器位置',
  `created_time` datetime NOT NULL COMMENT '服务器创建(购买)时间',
  `update_time` datetime NOT NULL COMMENT '服务器信息修改时间',
  `notes` text COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IP` (`IP`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `serverinfo`
--

LOCK TABLES `serverinfo` WRITE;
/*!40000 ALTER TABLE `serverinfo` DISABLE KEYS */;
/*!40000 ALTER TABLE `serverinfo` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `servers_manage`
--

DROP TABLE IF EXISTS `servers_manage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `servers_manage` (
  `name` varchar(128) NOT NULL,
  `hosts` varchar(1024) DEFAULT NULL,
  `rsa_key` varchar(1024) DEFAULT NULL,
  `lastip` varchar(1024) NOT NULL DEFAULT '0.0.0.0',
  `lasttime` varchar(1024) NOT NULL DEFAULT '2016-12-17 00:00:00',
  `created` int(10) DEFAULT NULL,
  `time` int(10) DEFAULT NULL,
  PRIMARY KEY (`name`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `servers_manage`
--

LOCK TABLES `servers_manage` WRITE;
/*!40000 ALTER TABLE `servers_manage` DISABLE KEYS */;
/*!40000 ALTER TABLE `servers_manage` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(240) DEFAULT NULL COMMENT '用户邮箱，即登录用户名',
  `password` varchar(240) DEFAULT NULL COMMENT '密码',
  `realname` varchar(240) DEFAULT NULL COMMENT '用户真实姓名',
  `phone` varchar(240) DEFAULT NULL COMMENT '用户电话',
  `created_time` datetime NOT NULL COMMENT '账号创建时间',
  `lastloginip` varchar(240) NOT NULL COMMENT '上次登录IP',
  `lastlogin` datetime NOT NULL COMMENT '上次登录时间',
  `notes` text COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `versions`
--

DROP TABLE IF EXISTS `versions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;

CREATE TABLE `versions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `versions_name` varchar(240) DEFAULT NULL COMMENT '版本名称',
  `versions` varchar(240) DEFAULT NULL COMMENT '版本号',
  `enviroment` varchar(240) DEFAULT NULL COMMENT '版本环境',
  `project` varchar(1024) NOT NULL DEFAULT 'teaching',
  `project_name` varchar(1024) NOT NULL DEFAULT '教学平台',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  `notes` text COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `versions` (`versions`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `versions`
--

LOCK TABLES `versions` WRITE;
/*!40000 ALTER TABLE `versions` DISABLE KEYS */;
/*!40000 ALTER TABLE `versions` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-12-17 10:07:46