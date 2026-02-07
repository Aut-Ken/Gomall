package main

import (
	"gomall/backend/internal/config"
	"gomall/backend/internal/database"
	"gomall/backend/internal/model"
	"log"
)

func main() {
	// Initialize config
	if err := config.Init("conf/config-dev.yaml"); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// Products to seed
	products := []model.Product{
		{
			Name:        "Apple Watch Series 9",
			Description: "智能手表，健康监测，GPS定位。",
			Price:       2999.00,
			Stock:       100,
			Category:    "手机数码",
			ImageURL:    "http://localhost:8080/photos/product_watch.jpg",
			Status:      1,
		},
		{
			Name:        "Sony WH-1000XM5",
			Description: "旗舰降噪耳机，沉浸式音效体验。",
			Price:       2499.00,
			Stock:       50,
			Category:    "手机数码",
			ImageURL:    "http://localhost:8080/photos/product_headphone.jpg",
			Status:      1,
		},
		{
			Name:        "Fujifilm X-T5",
			Description: "复古微单相机，4000万像素。",
			Price:       11999.00,
			Stock:       20,
			Category:    "手机数码",
			ImageURL:    "http://localhost:8080/photos/product_camera.jpg",
			Status:      1,
		},
		{
			Name:        "Nike Air Zoom",
			Description: "轻量缓震跑步鞋，透气舒适。",
			Price:       899.00,
			Stock:       200,
			Category:    "服饰鞋包",
			ImageURL:    "http://localhost:8080/photos/product_shoes.jpg",
			Status:      1,
		},
		{
			Name:        "PlayStation 5 Controller",
			Description: "PS5原装无线手柄，触觉反馈。",
			Price:       559.00,
			Stock:       150,
			Category:    "电脑办公",
			ImageURL:    "http://localhost:8080/photos/product_ps5.jpg",
			Status:      1,
		},
		{
			Name:        "iPad Air",
			Description: "10.9英寸平板电脑",
			Price:       4799.00,
			Stock:       80,
			Category:    "电脑办公",
			ImageURL:    "http://localhost:8080/photos/product_ipad.jpg", // New image
			Status:      1,
		},
		{
			Name:        "活着 (余华)",
			Description: "讲述了农村人福贵悲惨的人生遭遇。",
			Price:       45.00,
			Stock:       500,
			Category:    "礼品鲜花", // Temporary mapping to valid category or add "Books" to frontend if needed. Let's use 礼品鲜花 as misc
			ImageURL:    "http://localhost:8080/photos/product_book.jpg",
			Status:      1,
		},
		{
			Name:        "SK-II 神仙水",
			Description: "护肤精华露，修护肌肤。",
			Price:       1540.00,
			Stock:       80,
			Category:    "美妆护肤",
			ImageURL:    "http://localhost:8080/photos/product_cosmetics.jpg",
			Status:      1,
		},
		{
			Name:        "iPhone 15 Pro",
			Description: "钛金属机身，A17 Pro芯片。",
			Price:       7999.00,
			Stock:       100,
			Category:    "手机数码",
			ImageURL:    "http://localhost:8080/photos/product_iphone.jpg",
			Status:      1,
		},
		{
			Name:        "Handcrafted Leather Shoes",
			Description: "手工真皮皮鞋，商务休闲。",
			Price:       1299.00,
			Stock:       60,
			Category:    "服饰鞋包",
			ImageURL:    "http://localhost:8080/photos/product_leather_shoes.jpg",
			Status:      1,
		},
		// Adding more items to ensure variety and image coverage
		{
			Name:        "MacBook Pro 14",
			Description: "M3 Pro芯片，极致性能。",
			Price:       16999.00,
			Stock:       30,
			Category:    "电脑办公",
			ImageURL:    "http://localhost:8080/photos/product_macbook.jpg", // New image
			Status:      1,
		},
	}

	log.Println("Seeding products...")
	for _, p := range products {
		var count int64
		database.DB.Model(&model.Product{}).Where("name = ?", p.Name).Count(&count)
		if count == 0 {
			if err := database.DB.Create(&p).Error; err != nil {
				log.Printf("Failed to create product %s: %v", p.Name, err)
			} else {
				log.Printf("Created product: %s", p.Name)
			}
		} else {
			// Update existing just in case (e.g. image url)
			var existP model.Product
			database.DB.Where("name = ?", p.Name).First(&existP)
			existP.ImageURL = p.ImageURL
			existP.Price = p.Price
			existP.Stock = p.Stock
			existP.Description = p.Description
			existP.Category = p.Category
			database.DB.Save(&existP)
			log.Printf("Updated product: %s", p.Name)
		}
	}
	log.Println("Seeding complete.")
}
