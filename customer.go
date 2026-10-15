package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ฟังก์ชันเพิ่มลูกค้าใหม่ (Create)
func createCustomerDB(db *gorm.DB, name string, email string, phone string) (*Customer, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	phone = strings.TrimSpace(phone)

	if name == "" {
		err := errors.New("ชื่อลูกค้าต้องไม่ว่างเปล่า")
		fmt.Printf("[Error] createCustomerDB: %v\n", err)
		return nil, err
	}
	if email == "" {
		err := errors.New("อีเมลต้องไม่ว่างเปล่า")
		fmt.Printf("[Error] createCustomerDB: %v\n", err)
		return nil, err
	}
	if phone == "" {
		err := errors.New("เบอร์โทรศัพท์ต้องไม่ว่างเปล่า")
		fmt.Printf("[Error] createCustomerDB: %v\n", err)
		return nil, err
	}

	customer := &Customer{
		Name:  name,
		Email: email,
		Phone: phone,
	}

	// คำสั่งบันทึกข้อมูลลง table
	result := db.Create(customer)
	if result.Error != nil {
		fmt.Printf("[Error] createCustomerDB: %v\n", result.Error)
		return nil, fmt.Errorf("ไม่สามารถเพิ่มลูกค้าได้: %w", result.Error)
	}

	fmt.Printf("[Create Customer สำเร็จ] ID: %d | ชื่อ: '%s' | อีเมล: %s | โทร: %s\n",
		customer.ID, customer.Name, customer.Email, customer.Phone)
	return customer, nil
}

// ฟังก์ชันค้นหาลูกค้าด้วย Email (Read)
func getCustomerByEmailDB(db *gorm.DB, email string) (*Customer, error) {
	var customer Customer
	cleanEmail := strings.TrimSpace(email)
	result := db.Where("email = ?", cleanEmail).First(&customer)
	if result.Error != nil {
		fmt.Printf("[Error] getCustomerByEmailDB: ไม่พบลูกค้าที่มีอีเมล '%s'\n", cleanEmail)
		return nil, fmt.Errorf("ไม่พบลูกค้าที่มีอีเมล '%s': %w", cleanEmail, result.Error)
	}

	fmt.Printf("\n--- ข้อมูลลูกค้า (ค้นหาจาก Email: %s) ---\n", cleanEmail)
	fmt.Printf("ID: %d | ชื่อ: %s | อีเมล: %s | เบอร์โทรศัพท์: %s\n",
		customer.ID, customer.Name, customer.Email, customer.Phone)
	return &customer, nil
}

// ฟังก์ชันแก้ไขข้อมูลลูกค้าตาม ID (Update)
func updateCustomerByIDDB(db *gorm.DB, id uint, name string, email string, phone string) (*Customer, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	phone = strings.TrimSpace(phone)

	if name == "" || email == "" || phone == "" {
		err := errors.New("ข้อมูลชื่อ อีเมล และเบอร์โทรต้องไม่ว่างเปล่า")
		fmt.Printf("[Error] updateCustomerByIDDB: %v\n", err)
		return nil, err
	}

	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		fmt.Printf("[Error] updateCustomerByIDDB: ไม่พบลูกค้า ID %d\n", id)
		return nil, fmt.Errorf("ไม่พบลูกค้า ID %d", id)
	}

	customer.Name = name
	customer.Email = email
	customer.Phone = phone

	// คำสั่งบันทึกการอัปเดตลง table
	result := db.Save(&customer)
	if result.Error != nil {
		fmt.Printf("[Error] updateCustomerByIDDB: %v\n", result.Error)
		return nil, fmt.Errorf("ไม่สามารถอัปเดตข้อมูลลูกค้าได้: %w", result.Error)
	}

	fmt.Printf("[Update Customer สำเร็จ] แก้ไขลูกค้า ID %d: ชื่อ='%s' | อีเมล='%s' | โทร='%s'\n",
		customer.ID, customer.Name, customer.Email, customer.Phone)
	return &customer, nil
}

// ฟังก์ชันลบลูกค้าตาม ID (Delete)
// กฎ: ห้ามลบลูกค้าที่มีประวัติคำสั่งซื้อ
func deleteCustomerByIDDB(db *gorm.DB, id uint) error {
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		fmt.Printf("[Error] deleteCustomerByIDDB: ไม่พบลูกค้า ID %d\n", id)
		return fmt.Errorf("ไม่พบลูกค้า ID %d", id)
	}

	// ตรวจสอบว่ามีประวัติคำสั่งซื้อหรือไม่
	var orderCount int64
	db.Model(&Order{}).Where("customer_id = ?", id).Count(&orderCount)
	if orderCount > 0 {
		err := fmt.Errorf("ไม่สามารถลบลูกค้า ID %d ได้ เนื่องจากมีประวัติคำสั่งซื้อ %d รายการ", id, orderCount)
		fmt.Printf("[ป้องกันการลบ] %v\n", err)
		return err
	}

	// คำสั่งลบลูกค้า
	result := db.Delete(&customer)
	if result.Error != nil {
		fmt.Printf("[Error] deleteCustomerByIDDB: %v\n", result.Error)
		return fmt.Errorf("ไม่สามารถลบลูกค้าได้: %w", result.Error)
	}

	fmt.Printf("[Delete Customer สำเร็จ] ลบลูกค้า ID %d ('%s') เรียบร้อยแล้ว\n", id, customer.Name)
	return nil
}

// ฟังก์ชันดึงข้อมูลลูกค้าทั้งหมด (Read)
func getAllCustomersDB(db *gorm.DB) ([]Customer, error) {
	var customers []Customer
	result := db.Find(&customers)
	if result.Error != nil {
		fmt.Printf("[Error] getAllCustomersDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- รายชื่อลูกค้าทั้งหมด (%d รายการ) ---\n", len(customers))
	fmt.Printf("%-4s %-25s %-30s %-15s\n", "ID", "ชื่อลูกค้า", "อีเมล", "เบอร์โทร")
	fmt.Println(strings.Repeat("-", 78))
	for _, c := range customers {
		fmt.Printf("%-4d %-25s %-30s %-15s\n", c.ID, c.Name, c.Email, c.Phone)
	}
	return customers, nil
}

// ฟังก์ชันดึงข้อมูลลูกค้าตาม ID (Read)
func getCustomerByIDDB(db *gorm.DB, id uint) (*Customer, error) {
	var customer Customer
	result := db.First(&customer, id)
	if result.Error != nil {
		fmt.Printf("[Error] getCustomerByIDDB: ไม่พบลูกค้า ID %d\n", id)
		return nil, result.Error
	}

	fmt.Printf("\n--- ข้อมูลลูกค้า ID %d ---\n", id)
	fmt.Printf("ID: %d | ชื่อ: %s | อีเมล: %s | เบอร์โทร: %s\n",
		customer.ID, customer.Name, customer.Email, customer.Phone)
	return &customer, nil
}
