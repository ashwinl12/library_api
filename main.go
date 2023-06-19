package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type data struct {
	db *gorm.DB
}

func main() {
	dsn := "host=localhost user=onito_his_user password=onito@db dbname=library_db port=5432 TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto-migrate the user model
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	dbData := data {
		db: db,
	}

	// Define routes
	r.POST("/register", dbData.registerUser)
	r.POST("/login", dbData.loginUser)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func (dt data) registerUser(c *gin.Context) {
	var user User

	// Bind JSON request body to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Create user in the database
	if err := dt.db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully"})
}

func (dt data) loginUser(c *gin.Context) {
	var user User

	// Bind JSON request body to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Check if user exists in the database
	if err := dt.db.Where(&user).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(200, gin.H{"message": "User logged in successfully"})
}
