package main

import (
	"log"
	"net/http"

	"urlshortener/internal/db"
	"urlshortener/internal/routes"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	db.InitDB()
	db.GetConnection()
	// It is here just for local debugging without docker
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file located, trying run from env variables")
	}

	router := httprouter.New()
	router.GET("/:hash", routes.GetSourceUrl)
	router.POST("/", routes.CreateShortenedUrl)

	server := http.Server{
		Handler: router,
		Addr:    ":8080"}

	log.Fatal(server.ListenAndServe())
	defer db.CloseConnection()

}
