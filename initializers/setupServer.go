package initializers

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ringtho/inventory/db"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/routers"
)

func SetupServer()(*http.Server, *database.Queries, *sql.DB){
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in the environment")
	}

	conn := db.ConnectToDatabase()
	DB := database.New(conn)

	address := ":" + port
	server := &http.Server{
		Addr:    address,
		Handler: routers.Router(DB),
	}
	return server, DB, conn
}