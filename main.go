package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Sheet struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Description string
}

var db *gorm.DB

func initDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.AutoMigrate(&Sheet{})

}

func main() {
	godotenv.Load()
	initDB()
	router := gin.Default()

	router.GET("/sheets", func(ctx *gin.Context) {
		var sheets []Sheet
		db.Find(&sheets)
		ctx.JSON(http.StatusOK, sheets)

	})
	router.GET("/sheets/:id", func(ctx *gin.Context) {
		var sheet Sheet
		sheetID := ctx.Param("id")
		ID, err := strconv.Atoi(sheetID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		result := db.First(&sheet, ID)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Sheet not found"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return

		}

		ctx.JSON(http.StatusOK, sheet)
	})
	router.POST("/sheets", func(ctx *gin.Context) {
		var newSheet Sheet
		if err := ctx.ShouldBindJSON(&newSheet); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		result := db.Create(&newSheet)
		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sheet"})
			return
		}
		ctx.JSON(http.StatusCreated, newSheet)

	})
	router.DELETE("/sheets/:id", func(ctx *gin.Context) {
		sheetID := ctx.Param("id")
		ID, err := strconv.Atoi(sheetID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var sheet Sheet
		result := db.First(&sheet, ID)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Sheet not found"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		if err := db.Delete(&sheet).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sheet"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Sheet deleted successfully"})
	})

	router.Run("localhost:3000")
}
