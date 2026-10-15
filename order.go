package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// สถานะคำสั่งซื้อที่อนุญาต
var validOrderStatuses = map[string]bool{
	"Pending":   true,
	"Completed": true,
	"Cancelled": true,
}

// ==========================================
// ส่วนที่ 2: CREATE WITH TRANSACTION
// ==========================================

// ฟังก์ชันเพิ่มคำสั่งซื้อโดยใช้ Database Transaction (Create)
// ขั้นตอน:
// 1. ตรวจสอบว่ามีลูกค้าและสินค้าในระบบ
// 2. ตรวจสอบสต็อกว่าเพียงพอกับจำนวนที่สั่งหรือไม่
// 3. คำนวณ TotalPrice = Price * Quantity
// 4. ลดสต็อกสินค้า
// 5. บันทึกคำสั่งซื้อ (หากมีขั้นตอนใดล้มเหลว Transaction จะ Rollback กลับทั้งหมด)
func createOrderDB(db *gorm.DB, customerID uint, productID uint, quantity int) (*Order, error) {
	if quantity <= 0 {
		err := errors.New("จำนวนสินค้าที่สั่งซื้อต้องมากกว่า 0")
		fmt.Printf("[Error] createOrderDB: %v\n", err)
		return nil, err
	}

	var createdOrder Order

	// เริ่มใช้งาน Transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. ตรวจสอบลูกค้า
		var customer Customer
		if err := tx.First(&customer, customerID).Error; err != nil {
			return fmt.Errorf("ไม่พบข้อมูลลูกค้า ID %d", customerID)
		}

		// 2. ตรวจสอบสินค้า
		var product Product
		if err := tx.First(&product, productID).Error; err != nil {
			return fmt.Errorf("ไม่พบข้อมูลสินค้า ID %d", productID)
		}

		// 3. ตรวจสอบสต็อกคงเหลือ
		if product.Quantity < quantity {
			return fmt.Errorf("สต็อกสินค้า '%s' ไม่เพียงพอ (สั่ง: %d, คงเหลือ: %d)",
				product.Name, quantity, product.Quantity)
		}

		// 4. คำนวณยอดรวมราคา
		totalPrice := product.Price * float64(quantity)

		// 5. ตัดสต็อกสินค้า
		product.Quantity -= quantity
		if err := tx.Save(&product).Error; err != nil {
			return fmt.Errorf("เกิดข้อผิดพลาดในการตัดสต็อกสินค้า: %w", err)
		}

		// 6. บันทึกคำสั่งซื้อลง table
		createdOrder = Order{
			CustomerID: customerID,
			ProductID:  productID,
			Quantity:   quantity,
			TotalPrice: totalPrice,
			Status:     "Pending",
		}

		if err := tx.Create(&createdOrder).Error; err != nil {
			return fmt.Errorf("เกิดข้อผิดพลาดในการบันทึกคำสั่งซื้อ: %w", err)
		}

		createdOrder.Customer = customer
		createdOrder.Product = product
		return nil
	})

	if err != nil {
		fmt.Printf("[Transaction ล้มเหลว/Rollback] %v\n", err)
		return nil, err
	}

	fmt.Printf("[Create Order สำเร็จ] Order #%d | ลูกค้า: %s | สินค้า: %s x %d ชิ้น | ยอดรวม: %.2f บาท | สถานะ: %s\n",
		createdOrder.ID, createdOrder.Customer.Name, createdOrder.Product.Name,
		createdOrder.Quantity, createdOrder.TotalPrice, createdOrder.Status)
	return &createdOrder, nil
}

// ==========================================
// ส่วนที่ 3: READ
// ==========================================

// ฟังก์ชันดึงประวัติคำสั่งซื้อของลูกค้า (Read)
func getOrdersByCustomerDB(db *gorm.DB, customerID uint) ([]Order, error) {
	var orders []Order
	result := db.Where("customer_id = ?", customerID).
		Preload("Customer").
		Preload("Product").
		Preload("Product.Category").
		Find(&orders)
	if result.Error != nil {
		fmt.Printf("[Error] getOrdersByCustomerDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- ประวัติคำสั่งซื้อของลูกค้า ID %d (พบ %d รายการ) ---\n", customerID, len(orders))
	for _, o := range orders {
		fmt.Printf("- Order #%d | สินค้า: %-25s x %d | ยอดรวม: %10.2f บาท | สถานะ: %s\n",
			o.ID, o.Product.Name, o.Quantity, o.TotalPrice, o.Status)
	}
	return orders, nil
}

// ฟังก์ชันดึงคำสั่งซื้อตามสถานะ เช่น "Pending", "Completed" (Read)
func getOrdersByStatusDB(db *gorm.DB, status string) ([]Order, error) {
	var orders []Order
	status = strings.TrimSpace(status)
	result := db.Where("status = ?", status).
		Preload("Customer").
		Preload("Product").
		Preload("Product.Category").
		Find(&orders)
	if result.Error != nil {
		fmt.Printf("[Error] getOrdersByStatusDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- รายการคำสั่งซื้อสถานะ '%s' (พบ %d รายการ) ---\n", status, len(orders))
	for _, o := range orders {
		fmt.Printf("- Order #%d | ลูกค้า: %-15s | สินค้า: %-25s x %d | ยอดรวม: %10.2f บาท | สถานะ: %s\n",
			o.ID, o.Customer.Name, o.Product.Name, o.Quantity, o.TotalPrice, o.Status)
	}
	return orders, nil
}

// ฟังก์ชันดึงคำสั่งซื้อตามช่วงวันที่ (Read)
func getOrdersByDateRangeDB(db *gorm.DB, startDate time.Time, endDate time.Time) ([]Order, error) {
	var orders []Order
	result := db.Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Preload("Customer").
		Preload("Product").
		Preload("Product.Category").
		Find(&orders)
	if result.Error != nil {
		fmt.Printf("[Error] getOrdersByDateRangeDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- คำสั่งซื้อช่วงวันที่ %s ถึง %s (พบ %d รายการ) ---\n",
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), len(orders))
	for _, o := range orders {
		fmt.Printf("- Order #%d | วันที่: %s | ลูกค้า: %-15s | ยอดรวม: %10.2f บาท | สถานะ: %s\n",
			o.ID, o.CreatedAt.Format("2006-01-02 15:04"), o.Customer.Name, o.TotalPrice, o.Status)
	}
	return orders, nil
}

// ฟังก์ชันดึงคำสั่งซื้อพร้อมข้อมูลลูกค้าและสินค้าทั้งหมด (Read with Details)
func getOrdersWithDetailsDB(db *gorm.DB) ([]Order, error) {
	var orders []Order
	result := db.Preload("Customer").
		Preload("Product").
		Preload("Product.Category").
		Find(&orders)
	if result.Error != nil {
		fmt.Printf("[Error] getOrdersWithDetailsDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- รายการคำสั่งซื้อทั้งหมดพร้อมรายละเอียด (%d รายการ) ---\n", len(orders))
	fmt.Printf("%-5s %-20s %-28s %-6s %-12s %-10s\n", "ID", "ลูกค้า", "สินค้า (หมวดหมู่)", "จำนวน", "ยอดรวม (บาท)", "สถานะ")
	fmt.Println(strings.Repeat("-", 90))
	for _, o := range orders {
		catName := ""
		if o.Product.Category.Name != "" {
			catName = " (" + o.Product.Category.Name + ")"
		}
		fmt.Printf("%-5d %-20s %-28s %-6d %-12.2f %-10s\n",
			o.ID, o.Customer.Name, o.Product.Name+catName, o.Quantity, o.TotalPrice, o.Status)
	}
	return orders, nil
}

// ==========================================
// ส่วนที่ 4: UPDATE
// ==========================================

// ฟังก์ชันเปลี่ยนสถานะคำสั่งซื้อ (Update)
// กฎ: สถานะต้องเป็น Pending, Completed หรือ Cancelled เท่านั้น
func updateOrderStatusDB(db *gorm.DB, id uint, newStatus string) (*Order, error) {
	newStatus = strings.TrimSpace(newStatus)
	if !validOrderStatuses[newStatus] {
		err := fmt.Errorf("สถานะไม่ถูกต้อง '%s' (อนุญาตเฉพาะ: Pending, Completed, Cancelled)", newStatus)
		fmt.Printf("[Validation Error] updateOrderStatusDB: %v\n", err)
		return nil, err
	}

	var order Order
	if err := db.First(&order, id).Error; err != nil {
		fmt.Printf("[Error] updateOrderStatusDB: ไม่พบคำสั่งซื้อ ID %d\n", id)
		return nil, fmt.Errorf("ไม่พบคำสั่งซื้อ ID %d", id)
	}

	oldStatus := order.Status
	order.Status = newStatus
	result := db.Save(&order)
	if result.Error != nil {
		fmt.Printf("[Error] updateOrderStatusDB: %v\n", result.Error)
		return nil, result.Error
	}

	db.Preload("Customer").Preload("Product").First(&order, id)
	fmt.Printf("[Update Status สำเร็จ] Order #%d: เปลี่ยนสถานะจาก '%s' -> '%s'\n",
		order.ID, oldStatus, order.Status)
	return &order, nil
}

// ฟังก์ชันดึงคำสั่งซื้อตาม ID (Read)
func getOrderByIDDB(db *gorm.DB, id uint) (*Order, error) {
	var order Order
	result := db.Preload("Customer").Preload("Product").Preload("Product.Category").First(&order, id)
	if result.Error != nil {
		fmt.Printf("[Error] getOrderByIDDB: ไม่พบคำสั่งซื้อ ID %d\n", id)
		return nil, fmt.Errorf("ไม่พบคำสั่งซื้อ ID %d", id)
	}

	fmt.Printf("\n--- รายละเอียดคำสั่งซื้อ ID %d ---\n", id)
	fmt.Printf("Order #%d | ลูกค้า: %s | สินค้า: %s x %d | ยอดรวม: %.2f บาท | สถานะ: %s\n",
		order.ID, order.Customer.Name, order.Product.Name, order.Quantity, order.TotalPrice, order.Status)
	return &order, nil
}
