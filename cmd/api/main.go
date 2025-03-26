package main

import (
	"log"
	"net/http"
	"session-auth/internal/data"
	"session-auth/internal/database"

	"gorm.io/gorm"
)

type application struct {
	config *config
	models *data.Data
}

type config struct {
	db          *gorm.DB
	redisClient *database.RedisClient
}

func main() {
	// init db
	db := database.Setup()
	redisClient := database.SetupRedis()

	// init config
	cfg := &config{
		db:          db,
		redisClient: redisClient,
	}
	//init app
	app := application{
		config: cfg,
		models: data.NewModels(db),
	}

	if err := app.models.MigrateDB(); err != nil {
		panic("Failed to auto migrate the models..")
	}

	// config server
	srv := http.Server{
		Addr:    ":3000",
		Handler: app.setup(),
	}

	// start server
	log.Println("Server is running on port 3000")
	log.Fatal(srv.ListenAndServe())

}
