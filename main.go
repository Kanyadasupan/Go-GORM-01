package main

import (
	"fmt"
	"time"
)

// รองรับการเรียกใช้ package time เมื่อเปิดคอมเมนต์ทดสอบ
var _ = time.Now

func main() {
	// 1. เชื่อมต่อฐานข้อมูลและ AutoMigrate ตาราง (กำหนดอยู่ใน models.go)
	db := initDB()

	// 2. เตรียมข้อมูลทดสอบขั้นต่ำ (หากยังไม่มีข้อมูล)
	seedDatabase(db)

	fmt.Println()
	fmt.Println("===========================================================================")
	fmt.Println("            ระบบจัดการร้านค้า (Store Management System) - Go + GORM        ")
	fmt.Println("===========================================================================")

	// =========================================================================
	//  วิธีทดสอบทีละฟังก์ชัน:
	//  -> เอาเครื่องหมาย // ออก หน้าฟังก์ชันที่ต้องการทดสอบ แล้วรัน: go run .
	// =========================================================================

	// ---------------- ส่วนที่ 2: Create (บันทึกข้อมูล) ----------------
	// createCategoryDB(db, "Gaming Gear", "อุปกรณ์เล่นเกมระดับโปร")
	// createProductDB(db, "PlayStation 5 Pro", 29900.00, 10, 4)
	// createCustomerDB(db, "มานะ ใจดี", "mana@gmail.com", "089-111-2222")
	// createOrderDB(db, 1, 1, 2) // สั่งซื้อสินค้า ID 1 จำนวน 1 ชิ้น (ลดสต็อกด้วย Transaction)

	// ---------------- ส่วนที่ 3: Read / Conditional Queries (ค้นหาข้อมูล) ----------------
	getAllProductsDB(db)                                // 1. ดึงสินค้าทั้งหมด
	// getProductByIDDB(db, 1)                             // 2. ดึงสินค้าตาม ID
	// getProductByNameDB(db, "iPhone")                    // 3. ค้นหาชื่อสินค้าด้วย LIKE
	// getProductsByCategoryDB(db, 1)                      // 4. ดึงสินค้าตาม CategoryID
	// getProductsByPriceRangeDB(db, 1000, 10000)          // 5. ดึงสินค้าตามช่วงราคา Min - Max
	// getProductsLowStockDB(db, 5)                        // 6. ดึงสินค้าสต็อกน้อยกว่า 5
	// getProductsByCategoryAndPriceDB(db, 2, 500, 30000)  // 7. ดึงสินค้าตามหมวดและราคา
	// getProductsSortedByPriceDB(db, "ASC")               // 8. เรียงราคาน้อยไปมาก (หรือ "DESC")
	// getTopExpensiveProductsDB(db, 3)                    // 9. ดึงสินค้าที่แพงที่สุด 3 อันดับแรก
	// getCustomerByEmailDB(db, "somchai@gmail.com")       // 10. ค้นหาลูกค้าด้วย Email
	// getOrdersByCustomerDB(db, 1)                        // 11. ดึงคำสั่งซื้อของลูกค้า ID 1
	// getOrdersByStatusDB(db, "Pending")                  // 12. ดึงคำสั่งซื้อตามสถานะ
	// getOrdersByDateRangeDB(db, time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour)) // 13. ดึงคำสั่งซื้อตามช่วงวัน
	// getOrdersWithDetailsDB(db)                          // 14. ดึงคำสั่งซื้อทั้งหมดพร้อมรายละเอียด

	// ---------------- ส่วนที่ 4: Update (แก้ไขข้อมูล) ----------------
	// updateProductByIDDB(db, 1, "MacBook Air M3 15-inch", 44900.00, 8, 1)
	// updateProductPriceDB(db, 1, 38500.00)
	// updateProductStockDB(db, 1, +6)                     // ปรับสต็อก +5 ชิ้น
	// updateCustomerByIDDB(db, 1, "สมชาย สายไอที", "somchai@gmail.com", "081-111-9999")
	// updateOrderStatusDB(db, 1, "Completed")             // สถานะ: Pending, Completed, Cancelled

	// ---------------- ส่วนที่ 5: Delete (ลบและกู้คืนข้อมูล) ----------------
	// deleteCategoryByIDDB(db, 1)                         // ป้องกันการลบหากยังมีสินค้าผูกอยู่
	// deleteCustomerByIDDB(db, 1)                         // ป้องกันการลบหากมีคำสั่งซื้อผูกอยู่
	// hardDeleteProductDB(db, 1)                          // ป้องกันการลบหากอยู่ในประวัติคำสั่งซื้อ
	// softDeleteProductDB(db, 2)                          // Soft Delete สินค้า ID 2
	// getDeletedProductsDB(db)                            // ดึงสินค้าที่ถูก Soft Delete
	// restoreProductDB(db, 2)                             // กู้คืนสินค้า ID 2 กลับมาเป็น Active
}
