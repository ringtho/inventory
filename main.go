package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/ringtho/inventory/db"
	"github.com/ringtho/inventory/initializers"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/routers"
)

func init() {
	initializers.LoadDotEnvFile()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in the environment")
	}
	conn := db.ConnectToDatabase()
	defer conn.Close()

	DB := database.New(conn)

	address := ":" + port
	log.Printf("Server running on port %s\n", port)
	http.ListenAndServe(address, routers.Router(DB))
}