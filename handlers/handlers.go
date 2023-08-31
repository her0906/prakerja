package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"kasir-api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func SetRedisClient(client *redis.Client) {
	rdb = client
}

// Handler untuk mendapatkan daftar semua produk
func GetProducts(c *gin.Context) {
	productKeys, err := rdb.Keys(ctx, "product:*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil produk"})
		return
	}

	var products []models.Product
	for _, key := range productKeys {
		val, err := rdb.Get(ctx, key).Result()
		if err == nil {
			var product models.Product
			if err := json.Unmarshal([]byte(val), &product); err == nil {
				products = append(products, product)
			}
		}
	}

	c.JSON(http.StatusOK, products)
}

// Handler untuk mendapatkan detail produk berdasarkan ID
func GetProductByID(c *gin.Context) {
	productIDStr := c.Param("id")
	productKey := "product:" + productIDStr
	val, err := rdb.Get(ctx, productKey).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	var product models.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data produk"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Handler untuk membuat transaksi baru
func CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.BindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	productKey := "product:" + strconv.Itoa(transaction.ProductID)
	val, err := rdb.Get(ctx, productKey).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	var product models.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data produk"})
		return
	}

	// Hitung total harga transaksi
	total := product.Price * transaction.Quantity
	transaction.TotalPrice = total

	// Simpan transaksi ke dalam Redis
	transactionJSON, _ := json.Marshal(transaction)
	transactionKey := "transaction:" + strconv.Itoa(transaction.ID)
	rdb.Set(ctx, transactionKey, transactionJSON, 0)

	c.JSON(http.StatusCreated, transaction)
}
