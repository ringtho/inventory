package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/ringtho/inventory/initializers"
)

func init() {
	initializers.LoadDotEnvFile()
}

func main() {
	server, _, conn := initializers.SetupServer()
	defer conn.Close()
	log.Printf("Server running on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
