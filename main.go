package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ringtho/inventory/initializers"
	"github.com/ringtho/inventory/routers"
)

func init() {
	initializers.LoadDotEnvFile()
}


func main() {
	port := ":" + os.Getenv("PORT")
	fmt.Printf("Server running on port %s\n", port)
	http.ListenAndServe(port, routers.Router())
}