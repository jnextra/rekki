package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/josephn123/rekki/pkg/api"
	"github.com/josephn123/rekki/pkg/validators"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	vd := validators.NewValidators()
	emailAPI := api.NewEmailAPI(vd)

	log.Printf("Listening on port %v", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), emailAPI)
}
