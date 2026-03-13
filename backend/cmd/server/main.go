package main

import (
	"fmt"
	"log"
	"time"

	"aiimageproject/internal/middleware"
	"aiimageproject/internal/model"
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
	autoMigrate(db)

	jwtService := jwt.NewJWT(
		viper.GetString("jwt.secret"),
		viper.GetDuration("jwt.expire"),
	)

	r := setupRouter(db, jwtService)

	addr := fmt.Sprintf(":%d", viper.GetInt("server.port"))
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
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
		log.Fatalf("Failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Image{},
		&model.Category{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func setupRouter(db *gorm.DB, jwtService *jwt.JWT) *gin.Engine {
	gin.SetMode(viper.GetString("server.mode"))

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

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
	}

	return r
}
