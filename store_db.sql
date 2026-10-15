-- MariaDB dump 10.19  Distrib 10.4.32-MariaDB, for Win64 (AMD64)
--
-- Host: localhost    Database: store_db
-- ------------------------------------------------------
-- Server version	10.4.32-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `store_db`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `store_db` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

USE `store_db`;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `categories` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `description` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_categories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (1,'2026-10-15 20:18:48.345','2026-10-15 20:18:48.345',NULL,'Laptops & Computers','โน้ตบุ๊ก คอมพิวเตอร์ตั้งโต๊ะ และอุปกรณ์ประมวลผล'),(2,'2026-10-15 20:18:48.345','2026-10-15 20:18:48.345',NULL,'Smartphones & Tablets','สมาร์ตโฟน แท็บเล็ต และอุปกรณ์พกพาอัจฉริยะ'),(3,'2026-10-15 20:18:48.345','2026-10-15 20:18:48.345',NULL,'Peripherals & Accessories','เมาส์ คีย์บอร์ด หูฟัง และอุปกรณ์ต่อพ่วง'),(4,'2026-10-15 20:18:48.345','2026-10-15 20:18:48.345',NULL,'Gadgets & Cables','สายเชื่อมต่อ ฮับ อะแดปเตอร์ และอุปกรณ์เสริมไอที'),(5,'2026-10-15 20:18:48.372','2026-10-15 20:18:48.372',NULL,'Audio & Sound','หูฟัง ลำโพง และอุปกรณ์เสียง'),(6,'2026-10-15 20:30:22.073','2026-10-15 20:30:22.073',NULL,'Audio & Sound','หูฟัง ลำโพง และอุปกรณ์เสียง'),(7,'2026-10-15 20:38:52.009','2026-10-15 20:38:52.009',NULL,'Audio & Accessories','หูฟัง ลำโพง และอุปกรณ์เสียง'),(8,'2026-10-15 20:40:20.711','2026-10-15 20:40:20.711',NULL,'Audio & Sound','หูฟัง ลำโพง และอุปกรณ์เสียง'),(9,'2026-10-15 20:42:56.857','2026-10-15 20:42:56.857',NULL,'Audio & Sound','หูฟัง ลำโพง และอุปกรณ์เสียง'),(10,'2026-10-15 20:50:51.147','2026-10-15 20:50:51.147',NULL,'Audio & Sound','หูฟัง ลำโพง และอุปกรณ์เสียง');
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `customers`
--

DROP TABLE IF EXISTS `customers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `customers` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `phone` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_customers_email` (`email`),
  KEY `idx_customers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `customers`
--

LOCK TABLES `customers` WRITE;
/*!40000 ALTER TABLE `customers` DISABLE KEYS */;
INSERT INTO `customers` VALUES (1,'2026-10-15 20:18:48.352','2026-10-15 20:18:48.352',NULL,'สมชาย สายไอที','somchai@gmail.com','081-111-2222'),(2,'2026-10-15 20:18:48.352','2026-10-15 20:18:48.352',NULL,'วิภาดา รักเรียน','wiphada@hotmail.com','082-333-4444'),(3,'2026-10-15 20:18:48.352','2026-10-15 20:18:48.352',NULL,'กิตติศักดิ์ พัฒนา','kittisak@yahoo.com','083-555-6666'),(4,'2026-10-15 20:18:48.352','2026-10-15 20:18:48.352',NULL,'ณัฐธิดา สดใส','natthida@gmail.com','084-777-8888'),(5,'2026-10-15 20:18:48.352','2026-10-15 20:18:48.352',NULL,'ธนภัทร มั่งมี','thanapat@outlook.com','085-999-0000'),(6,'2026-10-15 20:18:48.352','2026-10-15 20:18:48.352',NULL,'อรทัย ใจดี','orathai@gmail.com','086-123-4567'),(7,'2026-10-15 20:18:48.352','2026-10-15 20:18:48.352',NULL,'ปกรณ์ กิจการดี','pakorn@gmail.com','087-234-5678'),(8,'2026-10-15 20:18:48.352','2026-10-15 20:18:48.352',NULL,'สิรินทรา นาวา','sirintra@outlook.com','088-345-6789'),(9,'2026-10-15 20:18:48.375','2026-10-15 20:18:48.404',NULL,'วีรชัย ผู้เจริญ','veerachai.new@gmail.com','089-999-8877'),(10,'2026-10-15 20:30:22.081','2026-10-15 20:30:22.081',NULL,'วีรชัย ใจกล้า','veerachai@gmail.com','089-112-3344'),(12,'2026-10-15 20:40:20.715','2026-10-15 20:40:20.745',NULL,'วีรชัย ผู้เจริญ','test_customer_79500@gmail.com','089-999-8877'),(13,'2026-10-15 20:42:56.862','2026-10-15 20:42:56.920',NULL,'วีรชัย ผู้เจริญ','test_customer_12900@gmail.com','089-999-8877'),(14,'2026-10-15 20:50:51.150','2026-10-15 20:50:51.209',NULL,'วีรชัย ผู้เจริญ','test_customer_20200@gmail.com','089-999-8877');
/*!40000 ALTER TABLE `customers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `orders` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `customer_id` bigint(20) unsigned NOT NULL,
  `product_id` bigint(20) unsigned NOT NULL,
  `quantity` bigint(20) NOT NULL,
  `total_price` decimal(10,2) NOT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'Pending',
  PRIMARY KEY (`id`),
  KEY `idx_orders_deleted_at` (`deleted_at`),
  KEY `idx_orders_customer_id` (`customer_id`),
  KEY `idx_orders_product_id` (`product_id`),
  CONSTRAINT `fk_customers_orders` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`),
  CONSTRAINT `fk_products_orders` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
INSERT INTO `orders` VALUES (1,'2026-10-15 20:18:48.355','2026-10-15 20:18:48.355',NULL,1,1,1,39900.00,'Completed'),(2,'2026-10-15 20:18:48.358','2026-10-15 20:18:48.358',NULL,1,10,1,3590.00,'Completed'),(3,'2026-10-15 20:18:48.359','2026-10-15 20:18:48.359',NULL,2,6,1,29900.00,'Completed'),(4,'2026-10-15 20:18:48.361','2026-10-15 20:18:48.361',NULL,2,14,2,3780.00,'Completed'),(5,'2026-10-15 20:18:48.363','2026-10-15 20:18:48.363',NULL,3,2,1,21900.00,'Pending'),(6,'2026-10-15 20:18:48.364','2026-10-15 20:18:48.364',NULL,3,12,1,3890.00,'Pending'),(7,'2026-10-15 20:18:48.365','2026-10-15 20:18:48.365',NULL,4,5,1,48900.00,'Completed'),(8,'2026-10-15 20:18:48.366','2026-10-15 20:18:48.366',NULL,5,7,1,43900.00,'Pending'),(9,'2026-10-15 20:18:48.368','2026-10-15 20:18:48.368',NULL,5,11,1,5290.00,'Completed'),(10,'2026-10-15 20:18:48.369','2026-10-15 20:18:48.369',NULL,6,14,1,1890.00,'Cancelled'),(11,'2026-10-15 20:18:48.370','2026-10-15 20:18:48.370',NULL,7,9,1,23900.00,'Completed'),(12,'2026-10-15 20:18:48.371','2026-10-15 20:18:48.371',NULL,8,3,1,44900.00,'Pending'),(13,'2026-10-15 20:18:48.378','2026-10-15 20:18:48.406',NULL,9,16,2,27980.00,'Completed'),(14,'2026-10-15 20:30:22.087','2026-10-15 20:30:22.150',NULL,10,18,2,27980.00,'Completed'),(15,'2026-10-15 20:40:20.719','2026-10-15 20:40:20.746',NULL,12,21,2,27980.00,'Completed'),(16,'2026-10-15 20:42:56.869','2026-10-15 20:42:56.922',NULL,13,23,2,27980.00,'Completed'),(17,'2026-10-15 20:50:51.158','2026-10-15 20:50:51.211',NULL,14,25,2,27980.00,'Completed');
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `products` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(150) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `quantity` bigint(20) NOT NULL,
  `category_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_deleted_at` (`deleted_at`),
  KEY `idx_products_category_id` (`category_id`),
  CONSTRAINT `fk_categories_products` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
INSERT INTO `products` VALUES (1,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'MacBook Air M3',39900.00,10,1),(2,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'ASUS Vivobook 15',21900.00,15,1),(3,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'Lenovo Legion 5 Pro',44900.00,8,1),(4,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'Dell XPS 13 Plus',52900.00,5,1),(5,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'iPhone 16 Pro Max',48900.00,12,2),(6,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'iPhone 16',29900.00,20,2),(7,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'Samsung Galaxy S24 Ultra',43900.00,10,2),(8,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'Samsung Galaxy A55 5G',13900.00,25,2),(9,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'iPad Air M2 11-inch',23900.00,14,2),(10,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'Logitech MX Master 3S Mouse',3590.00,30,3),(11,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'Logitech G Pro X Superlight 2',5290.00,4,3),(12,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'Keychron K2 Mechanical Keyboard',3890.00,3,3),(13,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'Razer BlackWidow V4 Pro',7990.00,7,3),(14,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'Anker 7-in-1 USB-C Hub',1890.00,35,4),(15,'2026-10-15 20:18:48.348','2026-10-15 20:18:48.348',NULL,'Belkin Connect Thunderbolt 4 Hub',8990.00,2,4),(16,'2026-10-15 20:18:48.374','2026-10-15 20:18:48.402',NULL,'Sony WH-1000XM5 (Black)',12500.00,13,5),(18,'2026-10-15 20:30:22.078','2026-10-15 20:30:22.143',NULL,'Sony WH-1000XM5 (Black)',12500.00,13,6),(20,'2026-10-15 20:38:52.011','2026-10-15 20:38:52.011',NULL,'Sony WH-1000XM5',13990.00,10,7),(21,'2026-10-15 20:40:20.714','2026-10-15 20:40:20.742',NULL,'Sony WH-1000XM5 (Black Edition)',12500.00,13,8),(23,'2026-10-15 20:42:56.860','2026-10-15 20:42:56.918',NULL,'Sony WH-1000XM5 (Black Edition)',12500.00,13,9),(25,'2026-10-15 20:50:51.149','2026-10-15 20:50:51.206',NULL,'Sony WH-1000XM5 (Black Edition)',12500.00,13,10);
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-10-15 20:56:20
