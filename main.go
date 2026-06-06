package main

import (
	"fmt"
	"notification-service/controller"
	"notification-service/models"
	redisclient "notification-service/redis"
	"notification-service/routes"
	"notification-service/worker"
	"os"

	_ "notification-service/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Notification Service API
// @version 1.0.0
// @description Sistema de notificacion en tiempo real con WebSockets y Redis
// @host localhost:9000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	godotenv.Load()

	// usa DB_URL si existe, si no construye el DSN
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	}

	fmt.Println("DSN:", dsn) // log temporal

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("No se pudo conectar a la base de datos")
	}
	fmt.Println("Base de datos conectada")

	controller.DB = db

	models.MigrateTables(db)

	redisclient.ConectarRedis()

	hub := worker.NewHub()
	controller.HubGlobal = hub

	worker.StartWorker(hub)

	r := gin.Default()

	r.Use(cors.New(cors.Config{

		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetRoutes(r)
	port := os.Getenv("PORT")

	if port == "" {
		port = "9000"
	}

	r.Run(":" + port)
}
