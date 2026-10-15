package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// กำหนดค่าการเชื่อมต่อ MySQL
const (
	host     = "localhost"
	port     = "3306"
	user     = "root"
	password = ""
	dbname   = "store_db"
)

// ==========================================
// ส่วนที่ 1: MODELS (โครงสร้างทั้ง 4 ตาราง)
// ==========================================

// โครงสร้าง Category (ประเภทสินค้า) - ความสัมพันธ์ 1:N กับ Product
type Category struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// โครงสร้าง Product (สินค้า) - Foreign Key -> Category, ความสัมพันธ์ 1:N กับ Order
type Product struct {
	gorm.Model
	Name       string   `gorm:"type:varchar(150);not null" json:"name"`
	Price      float64  `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	CategoryID uint     `gorm:"not null;index" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category,omitempty"`
	Orders     []Order  `gorm:"foreignKey:ProductID" json:"orders,omitempty"`
}

// โครงสร้าง Customer (ลูกค้า) - ความสัมพันธ์ 1:N กับ Order
type Customer struct {
	gorm.Model
	Name   string  `gorm:"type:varchar(100);not null" json:"name"`
	Email  string  `gorm:"type:varchar(100);not null;uniqueIndex" json:"email"`
	Phone  string  `gorm:"type:varchar(20);not null" json:"phone"`
	Orders []Order `gorm:"foreignKey:CustomerID" json:"orders,omitempty"`
}

// โครงสร้าง Order (คำสั่งซื้อ) - Foreign Key -> Customer และ Product
type Order struct {
	gorm.Model
	CustomerID uint     `gorm:"not null;index" json:"customer_id"`
	Customer   Customer `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"customer,omitempty"`
	ProductID  uint     `gorm:"not null;index" json:"product_id"`
	Product    Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"product,omitempty"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	TotalPrice float64  `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     string   `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"` // Pending, Completed, Cancelled
}

// ==========================================
// ส่วนที่ 1: DATABASE CONNECTION & MIGRATION
// ==========================================

// ฟังก์ชันเชื่อมต่อฐานข้อมูล MySQL และทำ AutoMigrate
func initDB() *gorm.DB {
	// 1. สร้างฐานข้อมูล store_db หากยังไม่มี
	createDatabaseIfNotExists()

	// 2. สร้าง connection DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)

	// 3. ปิดข้อความ Log รบกวน
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      false,
		},
	)

	// 4. เชื่อมต่อฐานข้อมูลด้วย GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	fmt.Println("Database connection success: เชื่อมต่อ MySQL สำเร็จ")

	// 5. ทำ AutoMigrate ทั้ง 4 ตาราง
	err = db.AutoMigrate(&Category{}, &Product{}, &Customer{}, &Order{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
	fmt.Println("Migration success: สร้างตาราง categories, products, customers, orders สำเร็จ")

	return db
}

// ฟังก์ชันสร้าง Database store_db บน MySQL หากยังไม่มี
func createDatabaseIfNotExists() {
	dsnRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port)
	sqlDB, err := sql.Open("mysql", dsnRoot)
	if err != nil {
		panic("Failed to connect to MySQL server: " + err.Error())
	}
	defer sqlDB.Close()

	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbname)
	_, err = sqlDB.Exec(createSQL)
	if err != nil {
		panic("Failed to create database: " + err.Error())
	}
}

// ==========================================
// ข้อมูลทดสอบขั้นต่ำ (Seed Data)
// ==========================================

// ฟังก์ชันบันทึกข้อมูลเริ่มต้น (Category: 4, Product: 15, Customer: 8, Order: 12)
func seedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&Category{}).Count(&count)
	if count > 0 {
		return // มีข้อมูลอยู่แล้ว ไม่ต้องเพิ่มซ้ำ
	}

	fmt.Println(">> กำลังเตรียมข้อมูลทดสอบเริ่มต้น (4 Categories, 15 Products, 8 Customers, 12 Orders)...")

	// 1. หมวดหมู่สินค้า (4 Categories)
	cats := []Category{
		{Name: "Laptops & Computers", Description: "โน้ตบุ๊ก คอมพิวเตอร์ตั้งโต๊ะ และอุปกรณ์ประมวลผล"},
		{Name: "Smartphones & Tablets", Description: "สมาร์ตโฟน แท็บเล็ต และอุปกรณ์พกพาอัจฉริยะ"},
		{Name: "Peripherals & Accessories", Description: "เมาส์ คีย์บอร์ด หูฟัง และอุปกรณ์ต่อพ่วง"},
		{Name: "Gadgets & Cables", Description: "สายเชื่อมต่อ ฮับ อะแดปเตอร์ และอุปกรณ์เสริมไอที"},
	}
	db.Create(&cats)

	// 2. สินค้า (15 Products หลากหลายราคาและสต็อก)
	prods := []Product{
		{Name: "MacBook Air M3", Price: 39900.00, Quantity: 10, CategoryID: cats[0].ID},
		{Name: "ASUS Vivobook 15", Price: 21900.00, Quantity: 15, CategoryID: cats[0].ID},
		{Name: "Lenovo Legion 5 Pro", Price: 44900.00, Quantity: 8, CategoryID: cats[0].ID},
		{Name: "Dell XPS 13 Plus", Price: 52900.00, Quantity: 5, CategoryID: cats[0].ID},

		{Name: "iPhone 16 Pro Max", Price: 48900.00, Quantity: 12, CategoryID: cats[1].ID},
		{Name: "iPhone 16", Price: 29900.00, Quantity: 20, CategoryID: cats[1].ID},
		{Name: "Samsung Galaxy S24 Ultra", Price: 43900.00, Quantity: 10, CategoryID: cats[1].ID},
		{Name: "Samsung Galaxy A55 5G", Price: 13900.00, Quantity: 25, CategoryID: cats[1].ID},
		{Name: "iPad Air M2 11-inch", Price: 23900.00, Quantity: 14, CategoryID: cats[1].ID},

		{Name: "Logitech MX Master 3S Mouse", Price: 3590.00, Quantity: 30, CategoryID: cats[2].ID},
		{Name: "Logitech G Pro X Superlight 2", Price: 5290.00, Quantity: 4, CategoryID: cats[2].ID}, // สต็อกต่ำ (<5)
		{Name: "Keychron K2 Mechanical Keyboard", Price: 3890.00, Quantity: 3, CategoryID: cats[2].ID}, // สต็อกต่ำ (<5)
		{Name: "Razer BlackWidow V4 Pro", Price: 7990.00, Quantity: 7, CategoryID: cats[2].ID},

		{Name: "Anker 7-in-1 USB-C Hub", Price: 1890.00, Quantity: 35, CategoryID: cats[3].ID},
		{Name: "Belkin Connect Thunderbolt 4 Hub", Price: 8990.00, Quantity: 2, CategoryID: cats[3].ID}, // สต็อกต่ำ (<5)
	}
	db.Create(&prods)

	// 3. ลูกค้า (8 Customers)
	custs := []Customer{
		{Name: "สมชาย สายไอที", Email: "somchai@gmail.com", Phone: "081-111-2222"},
		{Name: "วิภาดา รักเรียน", Email: "wiphada@hotmail.com", Phone: "082-333-4444"},
		{Name: "กิตติศักดิ์ พัฒนา", Email: "kittisak@yahoo.com", Phone: "083-555-6666"},
		{Name: "ณัฐธิดา สดใส", Email: "natthida@gmail.com", Phone: "084-777-8888"},
		{Name: "ธนภัทร มั่งมี", Email: "thanapat@outlook.com", Phone: "085-999-0000"},
		{Name: "อรทัย ใจดี", Email: "orathai@gmail.com", Phone: "086-123-4567"},
		{Name: "ปกรณ์ กิจการดี", Email: "pakorn@gmail.com", Phone: "087-234-5678"},
		{Name: "สิรินทรา นาวา", Email: "sirintra@outlook.com", Phone: "088-345-6789"},
	}
	db.Create(&custs)

	// 4. คำสั่งซื้อ (12 Orders)
	orders := []struct {
		cIdx int
		pIdx int
		qty  int
		st   string
	}{
		{0, 0, 1, "Completed"}, // สมชาย -> MacBook Air
		{0, 9, 1, "Completed"}, // สมชาย -> Logitech Mouse
		{1, 5, 1, "Completed"}, // วิภาดา -> iPhone 16
		{1, 13, 2, "Completed"}, // วิภาดา -> Anker USB-C Hub
		{2, 1, 1, "Pending"},   // กิตติศักดิ์ -> ASUS Vivobook
		{2, 11, 1, "Pending"},  // กิตติศักดิ์ -> Keychron K2
		{3, 4, 1, "Completed"}, // ณัฐธิดา -> iPhone 16 Pro Max
		{4, 6, 1, "Pending"},   // ธนภัทร -> Samsung S24 Ultra
		{4, 10, 1, "Completed"}, // ธนภัทร -> Logitech G Pro
		{5, 13, 1, "Cancelled"}, // อรทัย -> Anker USB-C Hub
		{6, 8, 1, "Completed"}, // ปกรณ์ -> iPad Air M2
		{7, 2, 1, "Pending"},   // สิรินทรา -> Lenovo Legion 5
	}

	for _, o := range orders {
		c := custs[o.cIdx]
		p := prods[o.pIdx]
		db.Create(&Order{
			CustomerID: c.ID,
			ProductID:  p.ID,
			Quantity:   o.qty,
			TotalPrice: p.Price * float64(o.qty),
			Status:     o.st,
		})
	}

	fmt.Println(">> เตรียมข้อมูลเริ่มต้นสำเร็จเรียบร้อยแล้ว!")
}
