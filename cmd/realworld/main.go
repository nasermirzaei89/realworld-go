package main

import (
	"fmt"
	"github.com/nasermirzaei89/realworld-go/internal/handlers"
	"github.com/nasermirzaei89/realworld-go/internal/repositories/inmem"
	"log"
	"net/http"
	"os"
)

func main() {
	// repositories
	userRepo := inmem.NewUserRepository()
	articleRepo := inmem.NewArticleRepository()

	// handler
	h := handlers.NewHandler(userRepo, articleRepo, secret())

	// serve
	err := http.ListenAndServe(addr(), h)
	if err != nil {
		log.Fatalln(fmt.Errorf("error on listen and serve http: %w", err))
	}
}

func secret() []byte {
	if env, ok := os.LookupEnv("JWT_SECRET"); ok {
		return []byte(env)
	}

	return []byte("secret")
}

func addr() string {
	if env, ok := os.LookupEnv("API_ADDRESS"); ok {
		return env
	}

	return "0.0.0.0:8080"
}
