package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ฟังก์ชันเพิ่มประเภทสินค้า (Create)
func createCategoryDB(db *gorm.DB, name string, description string) (*Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		err := errors.New("ชื่อประเภทสินค้าต้องไม่ว่างเปล่า")
		fmt.Printf("[Error] createCategoryDB: %v\n", err)
		return nil, err
	}

	category := &Category{
		Name:        name,
		Description: strings.TrimSpace(description),
	}

	// คำสั่งบันทึกข้อมูลลง table
	result := db.Create(category)
	if result.Error != nil {
		fmt.Printf("[Error] createCategoryDB: %v\n", result.Error)
		return nil, fmt.Errorf("ไม่สามารถเพิ่มประเภทสินค้าได้: %w", result.Error)
	}

	fmt.Printf("[Create Category สำเร็จ] ID: %d | ชื่อประเภท: '%s' | รายละเอียด: %s\n",
		category.ID, category.Name, category.Description)
	return category, nil
}

// ฟังก์ชันลบประเภทสินค้าโดยใช้ id (Delete)
// กฎ: ห้ามลบประเภทสินค้าที่ยังมีสินค้าผูกอยู่
func deleteCategoryByIDDB(db *gorm.DB, id uint) error {
	var category Category
	// ค้นหาประเภทสินค้าว่ามีอยู่หรือไม่
	if err := db.First(&category, id).Error; err != nil {
		fmt.Printf("[Error] deleteCategoryByIDDB: ไม่พบประเภทสินค้า ID %d\n", id)
		return fmt.Errorf("ไม่พบประเภทสินค้า ID %d", id)
	}

	// ตรวจสอบว่ามีสินค้าที่ใช้ประเภทนี้อยู่หรือไม่
	var productCount int64
	db.Model(&Product{}).Where("category_id = ?", id).Count(&productCount)
	if productCount > 0 {
		err := fmt.Errorf("ไม่สามารถลบประเภทสินค้า ID %d ได้ เนื่องจากยังมีสินค้าอยู่ในประเภทนี้ %d รายการ", id, productCount)
		fmt.Printf("[ป้องกันการลบ] %v\n", err)
		return err
	}

	// คำสั่งลบประเภทสินค้า
	result := db.Delete(&category)
	if result.Error != nil {
		fmt.Printf("[Error] deleteCategoryByIDDB: %v\n", result.Error)
		return fmt.Errorf("ไม่สามารถลบประเภทสินค้าได้: %w", result.Error)
	}

	fmt.Printf("[Delete Category สำเร็จ] ลบประเภทสินค้า ID %d ('%s') เรียบร้อยแล้ว\n", id, category.Name)
	return nil
}

// ฟังก์ชันดึงข้อมูลประเภทสินค้าทั้งหมด (Read)
func getAllCategoriesDB(db *gorm.DB) ([]Category, error) {
	var categories []Category
	result := db.Find(&categories)
	if result.Error != nil {
		fmt.Printf("[Error] getAllCategoriesDB: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("\n--- รายการประเภทสินค้าทั้งหมด (%d หมวดหมู่) ---\n", len(categories))
	for _, c := range categories {
		fmt.Printf("- หมวด ID %d: %-30s | รายละเอียด: %s\n", c.ID, c.Name, c.Description)
	}
	return categories, nil
}

// ฟังก์ชันดึงข้อมูลประเภทสินค้าตาม ID (Read)
func getCategoryByIDDB(db *gorm.DB, id uint) (*Category, error) {
	var category Category
	result := db.First(&category, id)
	if result.Error != nil {
		fmt.Printf("[Error] getCategoryByIDDB: ไม่พบประเภทสินค้า ID %d\n", id)
		return nil, result.Error
	}

	fmt.Printf("\n--- ข้อมูลประเภทสินค้า ID %d ---\n", id)
	fmt.Printf("ID: %d | ชื่อหมวด: %s | รายละเอียด: %s\n", category.ID, category.Name, category.Description)
	return &category, nil
}
