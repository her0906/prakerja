package main

import (
	"kasir-api/handlers"
	"kasir-api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Inisialisasi koneksi database SQLite
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Gagal menghubungkan database")
	}

	db.AutoMigrate(&models.Product{}, &models.Transaction{})

	// Inisialisasi klien Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Alamat server Redis Anda
	})

	// Inisialisasi router
	router := gin.Default()

	// Atur koneksi database dan klien Redis di dalam handler
	// handlers.SetDB(db)
	handlers.SetRedisClient(rdb)

	// Endpoint untuk mendapatkan daftar semua produk
	router.GET("/products", handlers.GetProducts)

	// Endpoint untuk mendapatkan detail produk berdasarkan ID
	router.GET("/products/:id", handlers.GetProductByID)

	// Endpoint untuk membuat transaksi baru
	router.POST("/transactions", handlers.CreateTransaction)

	// Jalankan server
	router.Run(":8080")
}
