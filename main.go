package main

import (
	"fmt"
	"notification-service/controller"
	"notification-service/models"
	redisclient "notification-service/redis"
	"notification-service/routes"
	"notification-service/worker"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port =%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("No se pudo conecatar a la base de datos")

	}
	fmt.Print("Base de datos conectada")

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

	routes.SetRoutes(r)
	port := os.Getenv("PORT")

	if port == "" {
		port = "9000"
	}

	r.Run(":" + port)
}
