package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/ringtho/inventory/initializers"
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
	conn := initializers.ConnectToDatabase()
	defer conn.Close()

	address := ":" + port
	log.Printf("Server running on port %s\n", port)
	http.ListenAndServe(address, routers.Router())
}