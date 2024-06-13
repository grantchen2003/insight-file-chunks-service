package main

import (
	"fmt"
	"log"
	"os"

	serverPackage "github.com/grantchen2003/insight/filechunks/internal/server"

	"github.com/grantchen2003/insight/filechunks/internal/config"
)

func main() {
	env := os.Getenv("ENV")

	log.Printf("ENV=%s", env)

	if err := config.LoadEnvVars(env); err != nil {
		log.Fatalf("failed to load env vars")
	}

	address := fmt.Sprintf("%s:%s", os.Getenv("DOMAIN"), os.Getenv("PORT"))

	server := serverPackage.NewServer()

	if err := server.Start(address); err != nil {
		log.Fatalf("failed to start server")
	}
}
