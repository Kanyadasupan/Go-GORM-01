package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ==========================================
// ส่วนที่ 2: CREATE
// ==========================================

// ฟังก์ชันเพิ่มสินค้าพร้อมระบุประเภท (Create)
func createProductDB(db *gorm.DB, name string, price float64, quantity int, categoryID uint) (*Product, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		err := errors.New("ชื่อสินค้าต้องไม่ว่างเปล่า")
		fmt.Printf("[Error] createProductDB: %v\n", err)
		return nil, err
	}
	if price < 0 {
		err := errors.New("ราคาสินค้าต้องไม่ติดลบ")
		fmt.Printf("[Error] createProductDB: %v\n", err)
		return nil, err
	}
	if quantity < 0 {
		err := errors.New("จำนวนสต็อกต้องไม่ติดลบ")
		fmt.Printf("[Error] createProductDB: %v\n", err)
		return nil, err
	}

	// ตรวจสอบว่า Category มีอยู่จริงหรือไม่
	var category Category
	if err := db.First(&category, categoryID).Error; err != nil {
		errNotFound := fmt.Errorf("ไม่พบ Category ID %d ในระบบ", categoryID)
		fmt.Printf("[Error] createProductDB: %v\n", errNotFound)
		return nil, errNotFound
	}

	product := &Product{
		Name:       name,
		Price:      price,
		Quantity:   quantity,
		CategoryID: categoryID,
	}

	// คำสั่งบันทึกข้อมูลลง table
	result := db.Create(product)
	if result.Error != nil {
		fmt.Printf("[Error] createProductDB: %v\n", result.Error)
		return nil, fmt.Errorf("ไม่สามารถบันทึกสินค้าได้: %w", result.Error)
	}

	product.Category = category
	fmt.Printf("[Create Product สำเร็จ] ID: %d | ชื่อ: '%s' | ราคา: %.2f บาท | สต็อก: %d ชิ้น | หมวดหมู่: %s\n",
		product.ID, product.Name, product.Price, product.Quantity, category.Name)
	return product, nil
}

// ==========================================
// ส่วนที่ 3: READ / CONDITIONAL QUERIES
// ==========================================

// ฟังก์ชันเรียกดูสินค้าทั้งหมด (Read)
func getAllProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	result := db.Preload("Category").Find(&products)
	if result.Error != nil {
		fmt.Printf("[Error] getAllProductsDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- รายการสินค้าทั้งหมดในระบบ (%d รายการ) ---\n", len(products))
	fmt.Printf("%-4s %-32s %-12s %-8s %-25s\n", "ID", "ชื่อสินค้า", "ราคา (บาท)", "สต็อก", "หมวดหมู่")
	fmt.Println(strings.Repeat("-", 85))
	for _, p := range products {
		fmt.Printf("%-4d %-32s %-12.2f %-8d %-25s\n", p.ID, p.Name, p.Price, p.Quantity, p.Category.Name)
	}
	return products, nil
}

// ฟังก์ชันค้นหาสินค้าจาก ID (Read)
func getProductByIDDB(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	result := db.Preload("Category").First(&product, id)
	if result.Error != nil {
		fmt.Printf("[Error] getProductByIDDB: ไม่พบสินค้า ID %d\n", id)
		return nil, fmt.Errorf("ไม่พบสินค้า ID %d: %w", id, result.Error)
	}

	fmt.Printf("\n--- ข้อมูลสินค้า ID %d ---\n", id)
	fmt.Printf("ID: %d | ชื่อ: %s | ราคา: %.2f บาท | สต็อกคงเหลือ: %d ชิ้น | หมวดหมู่: %s\n",
		product.ID, product.Name, product.Price, product.Quantity, product.Category.Name)
	return &product, nil
}

// ฟังก์ชันค้นหาชื่อสินค้าบางส่วนด้วย LIKE (Read)
func getProductByNameDB(db *gorm.DB, keyword string) ([]Product, error) {
	var products []Product
	result := db.Where("name LIKE ?", "%"+strings.TrimSpace(keyword)+"%").Preload("Category").Find(&products)
	if result.Error != nil {
		fmt.Printf("[Error] getProductByNameDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- ผลค้นหาสินค้าด้วยคำค้น '%s' (พบ %d รายการ) ---\n", keyword, len(products))
	for _, p := range products {
		fmt.Printf("- ID %-3d: %-30s | ราคา: %10.2f บาท | สต็อก: %3d ชิ้น | หมวด: %s\n",
			p.ID, p.Name, p.Price, p.Quantity, p.Category.Name)
	}
	return products, nil
}

// ฟังก์ชันค้นหาสินค้าตาม CategoryID (Read)
func getProductsByCategoryDB(db *gorm.DB, categoryID uint) ([]Product, error) {
	var products []Product
	result := db.Where("category_id = ?", categoryID).Preload("Category").Find(&products)
	if result.Error != nil {
		fmt.Printf("[Error] getProductsByCategoryDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- รายการสินค้าในหมวดหมู่ ID %d (พบ %d รายการ) ---\n", categoryID, len(products))
	for _, p := range products {
		fmt.Printf("- ID %-3d: %-30s | ราคา: %10.2f บาท | สต็อก: %3d ชิ้น\n",
			p.ID, p.Name, p.Price, p.Quantity)
	}
	return products, nil
}

// ฟังก์ชันค้นหาสินค้าตามช่วงราคา Min - Max (Read)
func getProductsByPriceRangeDB(db *gorm.DB, minPrice float64, maxPrice float64) ([]Product, error) {
	var products []Product
	result := db.Where("price BETWEEN ? AND ?", minPrice, maxPrice).Preload("Category").Find(&products)
	if result.Error != nil {
		fmt.Printf("[Error] getProductsByPriceRangeDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- สินค้าราคาช่วง %.2f - %.2f บาท (พบ %d รายการ) ---\n", minPrice, maxPrice, len(products))
	for _, p := range products {
		fmt.Printf("- ID %-3d: %-30s | ราคา: %10.2f บาท | หมวด: %s\n",
			p.ID, p.Name, p.Price, p.Category.Name)
	}
	return products, nil
}

// ฟังก์ชันค้นหาสินค้าสต็อกน้อยกว่าค่ากำหนด (Read)
func getProductsLowStockDB(db *gorm.DB, threshold int) ([]Product, error) {
	var products []Product
	result := db.Where("quantity < ?", threshold).Preload("Category").Find(&products)
	if result.Error != nil {
		fmt.Printf("[Error] getProductsLowStockDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- [แจ้งเตือน] สินค้าสต็อกต่ำกว่า %d ชิ้น (พบ %d รายการ) ---\n", threshold, len(products))
	for _, p := range products {
		fmt.Printf("- ID %-3d: %-30s | สต็อกคงเหลือ: %d ชิ้น (หมวด: %s)\n",
			p.ID, p.Name, p.Quantity, p.Category.Name)
	}
	return products, nil
}

// ฟังก์ชันค้นหาสินค้าตามประเภทและช่วงราคา (Read)
func getProductsByCategoryAndPriceDB(db *gorm.DB, categoryID uint, minPrice float64, maxPrice float64) ([]Product, error) {
	var products []Product
	result := db.Where("category_id = ? AND price BETWEEN ? AND ?", categoryID, minPrice, maxPrice).Preload("Category").Find(&products)
	if result.Error != nil {
		fmt.Printf("[Error] getProductsByCategoryAndPriceDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- สินค้าหมวดหมู่ ID %d ช่วงราคา %.2f - %.2f บาท (พบ %d รายการ) ---\n", categoryID, minPrice, maxPrice, len(products))
	for _, p := range products {
		fmt.Printf("- ID %-3d: %-30s | ราคา: %10.2f บาท | หมวด: %s\n",
			p.ID, p.Name, p.Price, p.Category.Name)
	}
	return products, nil
}

// ฟังก์ชันเรียงราคาสินค้า ASC หรือ DESC (Read)
func getProductsSortedByPriceDB(db *gorm.DB, order string) ([]Product, error) {
	order = strings.ToUpper(strings.TrimSpace(order))
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}
	var products []Product
	result := db.Order("price " + order).Preload("Category").Find(&products)
	if result.Error != nil {
		fmt.Printf("[Error] getProductsSortedByPriceDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- สินค้าเรียงลำดับราคา (%s) ---\n", order)
	for i, p := range products {
		fmt.Printf("%2d. ID %-3d: %-30s | ราคา: %10.2f บาท\n", i+1, p.ID, p.Name, p.Price)
	}
	return products, nil
}

// ฟังก์ชันดึงสินค้าที่แพงที่สุด N รายการแรก (Read)
func getTopExpensiveProductsDB(db *gorm.DB, limit int) ([]Product, error) {
	if limit <= 0 {
		limit = 5
	}
	var products []Product
	result := db.Order("price DESC").Limit(limit).Preload("Category").Find(&products)
	if result.Error != nil {
		fmt.Printf("[Error] getTopExpensiveProductsDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- สินค้าที่แพงที่สุด %d อันดับแรก ---\n", limit)
	for i, p := range products {
		fmt.Printf("%d. %-32s | ราคา: %10.2f บาท | หมวด: %s\n",
			i+1, p.Name, p.Price, p.Category.Name)
	}
	return products, nil
}

// ==========================================
// ส่วนที่ 4: UPDATE
// ==========================================

// ฟังก์ชันแก้ไขข้อมูลสินค้าตาม ID (Update)
func updateProductByIDDB(db *gorm.DB, id uint, name string, price float64, quantity int, categoryID uint) (*Product, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		err := errors.New("ชื่อสินค้าต้องไม่ว่างเปล่า")
		fmt.Printf("[Error] updateProductByIDDB: %v\n", err)
		return nil, err
	}
	if price < 0 {
		err := errors.New("ราคาสินค้าต้องไม่ติดลบ")
		fmt.Printf("[Error] updateProductByIDDB: %v\n", err)
		return nil, err
	}
	if quantity < 0 {
		err := errors.New("จำนวนสต็อกต้องไม่ติดลบ")
		fmt.Printf("[Error] updateProductByIDDB: %v\n", err)
		return nil, err
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		fmt.Printf("[Error] updateProductByIDDB: ไม่พบสินค้า ID %d\n", id)
		return nil, fmt.Errorf("ไม่พบสินค้า ID %d", id)
	}

	product.Name = name
	product.Price = price
	product.Quantity = quantity
	product.CategoryID = categoryID

	// คำสั่งบันทึกการอัปเดตลง table
	result := db.Save(&product)
	if result.Error != nil {
		fmt.Printf("[Error] updateProductByIDDB: %v\n", result.Error)
		return nil, fmt.Errorf("ไม่สามารถอัปเดตสินค้าได้: %w", result.Error)
	}

	fmt.Printf("[Update Product สำเร็จ] แก้ไขสินค้า ID %d เป็น '%s' | ราคา: %.2f บาท | สต็อก: %d ชิ้น\n",
		product.ID, product.Name, product.Price, product.Quantity)
	return &product, nil
}

// ฟังก์ชันแก้ไขราคาสินค้า (Update) - ราคาต้องไม่ติดลบ
func updateProductPriceDB(db *gorm.DB, id uint, newPrice float64) (*Product, error) {
	if newPrice < 0 {
		err := errors.New("ราคาสินค้าต้องไม่ติดลบ")
		fmt.Printf("[Validation Error] updateProductPriceDB: %v\n", err)
		return nil, err
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		fmt.Printf("[Error] updateProductPriceDB: ไม่พบสินค้า ID %d\n", id)
		return nil, fmt.Errorf("ไม่พบสินค้า ID %d", id)
	}

	oldPrice := product.Price
	product.Price = newPrice
	result := db.Save(&product)
	if result.Error != nil {
		fmt.Printf("[Error] updateProductPriceDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("[Update Price สำเร็จ] สินค้า ID %d ('%s'): ปรับราคาจาก %.2f -> %.2f บาท\n",
		product.ID, product.Name, oldPrice, product.Price)
	return &product, nil
}

// ฟังก์ชันเพิ่ม/ลดสต็อกสินค้า (Update) - สต็อกต้องไม่ติดลบ
func updateProductStockDB(db *gorm.DB, id uint, deltaQuantity int) (*Product, error) {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		fmt.Printf("[Error] updateProductStockDB: ไม่พบสินค้า ID %d\n", id)
		return nil, fmt.Errorf("ไม่พบสินค้า ID %d", id)
	}

	newStock := product.Quantity + deltaQuantity
	if newStock < 0 {
		err := fmt.Errorf("สต็อกสินค้าไม่สามารถติดลบได้ (สต็อกเดิม: %d, เปลี่ยนแปลง: %+d, ผลลัพธ์: %d)", product.Quantity, deltaQuantity, newStock)
		fmt.Printf("[Validation Error] updateProductStockDB: %v\n", err)
		return nil, err
	}

	oldStock := product.Quantity
	product.Quantity = newStock
	result := db.Save(&product)
	if result.Error != nil {
		fmt.Printf("[Error] updateProductStockDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("[Update Stock สำเร็จ] สินค้า ID %d ('%s'): ปรับสต็อก %+d ชิ้น (เดิม: %d -> ปัจจุบัน: %d ชิ้น)\n",
		product.ID, product.Name, deltaQuantity, oldStock, product.Quantity)
	return &product, nil
}

// ==========================================
// ส่วนที่ 5: DELETE
// ==========================================

// ฟังก์ชัน Soft Delete สินค้า (Delete)
func softDeleteProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		fmt.Printf("[Error] softDeleteProductDB: ไม่พบสินค้า ID %d\n", id)
		return fmt.Errorf("ไม่พบสินค้า ID %d", id)
	}

	result := db.Delete(&product)
	if result.Error != nil {
		fmt.Printf("[Error] softDeleteProductDB: %v\n", result.Error)
		return fmt.Errorf("ไม่สามารถ Soft Delete สินค้าได้: %w", result.Error)
	}

	fmt.Printf("[Soft Delete สำเร็จ] สินค้า ID %d ('%s') ถูก Soft Delete เรียบร้อยแล้ว\n", id, product.Name)
	return nil
}

// ฟังก์ชันลบสินค้าถาวร (Hard Delete)
// กฎ: ห้ามลบสินค้าที่มีประวัติอยู่ในคำสั่งซื้อ (Orders)
func hardDeleteProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.Unscoped().First(&product, id).Error; err != nil {
		fmt.Printf("[Error] hardDeleteProductDB: ไม่พบสินค้า ID %d\n", id)
		return fmt.Errorf("ไม่พบสินค้า ID %d", id)
	}

	// ตรวจสอบว่าสินค้ามีประวัติอยู่ในคำสั่งซื้อหรือไม่
	var orderCount int64
	db.Model(&Order{}).Where("product_id = ?", id).Count(&orderCount)
	if orderCount > 0 {
		err := fmt.Errorf("ไม่สามารถลบสินค้า ID %d ถาวรได้ เนื่องจากมีประวัติอยู่ในคำสั่งซื้อ %d รายการ", id, orderCount)
		fmt.Printf("[ป้องกันการลบ] %v\n", err)
		return err
	}

	// ลบออกจากฐานข้อมูลถาวร
	result := db.Unscoped().Delete(&product)
	if result.Error != nil {
		fmt.Printf("[Error] hardDeleteProductDB: %v\n", result.Error)
		return fmt.Errorf("ไม่สามารถลบสินค้าถาวรได้: %w", result.Error)
	}

	fmt.Printf("[Hard Delete สำเร็จ] ลบสินค้า ID %d ('%s') ออกจากฐานข้อมูลถาวรเรียบร้อยแล้ว\n", id, product.Name)
	return nil
}

// ฟังก์ชันดึงสินค้าที่ถูก Soft Delete (Read)
func getDeletedProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	result := db.Unscoped().Where("deleted_at IS NOT NULL").Preload("Category").Find(&products)
	if result.Error != nil {
		fmt.Printf("[Error] getDeletedProductsDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- รายการสินค้าที่ถูก Soft Delete (%d รายการ) ---\n", len(products))
	for _, dp := range products {
		fmt.Printf("- Soft Deleted ID %d: %-30s | วันที่ลบ: %v\n", dp.ID, dp.Name, dp.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}
	return products, nil
}

// ฟังก์ชันกู้คืนสินค้าที่ถูก Soft Delete (Restore)
func restoreProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.Unscoped().First(&product, id).Error; err != nil {
		fmt.Printf("[Error] restoreProductDB: ไม่พบสินค้า ID %d\n", id)
		return fmt.Errorf("ไม่พบสินค้า ID %d", id)
	}

	if !product.DeletedAt.Valid {
		err := fmt.Errorf("สินค้า ID %d ไม่ได้อยู่ในสถานะถูกลบ", id)
		fmt.Printf("[เตือน] %v\n", err)
		return err
	}

	// ตั้งค่า deleted_at ให้เป็น NULL เพื่อคืนสถานะกลับมา Active
	result := db.Unscoped().Model(&product).Update("deleted_at", nil)
	if result.Error != nil {
		fmt.Printf("[Error] restoreProductDB: %v\n", result.Error)
		return fmt.Errorf("ไม่สามารถกู้คืนสินค้าได้: %w", result.Error)
	}

	fmt.Printf("[Restore สำเร็จ] กู้คืนสินค้า ID %d ('%s') กลับมาเป็น Active เรียบร้อยแล้ว\n", id, product.Name)
	return nil
}
