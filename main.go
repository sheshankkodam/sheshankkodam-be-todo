package main

import (
	"log"
	"os"
	"sync"

	"github.com/heroku/sheshankkodam-be-todo/pkg/api"
	"github.com/heroku/sheshankkodam-be-todo/pkg/database"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("database url must be set")
	}

	db, dbErr := database.Open(dbUrl)
	if dbErr != nil {
		log.Fatalf("Unable to open database session, error=%s", dbErr)
	}

	if dbInitErr := db.Init(); dbInitErr != nil {
		log.Fatalf("Unable to initialise database session, error=%s", dbInitErr)
	}

	wg := sync.WaitGroup{}
	httpServer := api.NewServer(port, db)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if httpErr := httpServer.Run(); httpErr != nil {
			log.Fatalf("Error starting HTTP server, error=%s", httpErr)
		}
	}()

	defer httpServer.Close()
	// Wait for everything to shutdown
	wg.Wait()
}
