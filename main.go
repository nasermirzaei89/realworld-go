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
	h := handlers.NewHandler(userRepo, articleRepo, secret())

	// serve
	err := http.ListenAndServe(addr(), h)
	if err != nil {
		log.Fatalln(fmt.Errorf("error on listen and serve http: %w", err))
	}
}

func secret() []byte {
	if env, ok := os.LookupEnv("SECRET"); ok {
		return []byte(env)
	}

	return []byte("secret")
}

func addr() string {
	if env, ok := os.LookupEnv("APIURL"); ok {
		return env
	}

	return "http://localhost:3000/api"
}
