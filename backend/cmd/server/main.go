package main

import (
	"fmt"
	"log"
	"time"

	"aiimageproject/internal/handler"
	"aiimageproject/internal/middleware"
	"aiimageproject/pkg/jwt"
	"aiimageproject/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	initConfig()

	db := initDB()
	if db != nil {
		autoMigrate(db)
	}

	jwtService := jwt.NewJWT(
		viper.GetString("jwt.secret"),
		viper.GetDuration("jwt.expire"),
	)

	r := setupRouter(db, jwtService)

	addr := fmt.Sprintf(":%d", viper.GetInt("server.port"))
	log.Printf("Server starting on %s", addr)
	log.Printf("Image Engine: %s", viper.GetString("image_engine.base_url"))
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.SetDefault("server.port", 8081)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("image_engine.base_url", "http://localhost:8080")
	viper.SetDefault("image_engine.timeout", 120)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Failed to read config file: %v, using defaults", err)
	}

	log.Printf("Config loaded - Image Engine URL: %s", viper.GetString("image_engine.base_url"))
}

func initDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.dbname"),
		viper.GetString("database.charset"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Warning: Failed to connect database: %v", err)
		log.Printf("Running without database connection...")
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Warning: Failed to get sql.DB: %v", err)
		return nil
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("Database connected successfully")
	return db
}

func autoMigrate(db *gorm.DB) {
	// TODO: Add models when needed
}

func setupRouter(db *gorm.DB, jwtService *jwt.JWT) *gin.Engine {
	gin.SetMode(viper.GetString("server.mode"))

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	r.GET("/health", handler.HealthCheck)

	r.GET("/", func(c *gin.Context) {
		response.Success(c, gin.H{
			"name":    "AI Image Project API",
			"version": "1.0.0",
		})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			response.Success(c, gin.H{"message": "pong"})
		})

		api.GET("/models", handler.GetModels)
		api.GET("/models/capabilities", handler.GetCapabilities)
		api.GET("/models/capability/:capability", handler.GetModelsByCapability)

		api.POST("/text2image", handler.Text2Image)
		api.POST("/image/edit", handler.ImageEdit)
	}

	return r
}
