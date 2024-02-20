package main

import (
	"log"
	"rinha-backend-2024-q1/controllers"
	"rinha-backend-2024-q1/lib/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	dbInstance, err := db.OpenConn()
	if err != nil {
		log.Fatal(err)
	}

	defer dbInstance.Close()

	router := gin.Default()

	router.POST("/clientes/:id/transacoes", controllers.NewTransaction)
	router.GET("/clientes/:id/extrato", controllers.Extrato)
	router.Run(":8080")
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
