package main

import (
	"fmt"
	"github.com/nasermirzaei89/realworld-go/handlers"
	"github.com/nasermirzaei89/realworld-go/repositories/inmem"
	"log"
	"net/http"
	"os"
)

func main() {
	// repositories
	userRepo := inmem.NewUserRepository()
	articleRepo := inmem.NewArticleRepository()

	// handler
	h := handlers.NewHandler(userRepo, articleRepo)

	// serve
	err := http.ListenAndServe(addr(), h)
	if err != nil {
		log.Fatalln(fmt.Errorf("error on listen and serve http: %w", err))
	}
}

func addr() string {
	if env, ok := os.LookupEnv("APIURL"); ok {
		return env
	}

	return "http://localhost:3000/api"
}
