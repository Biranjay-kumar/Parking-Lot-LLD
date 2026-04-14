package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"parking_slot/handlers"
	"parking_slot/repo"
	"parking_slot/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:admin@localhost:5432/practise?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repoLayer := repo.NewRepo(pool)
	serviceLayer := services.NewService(repoLayer)
	handlerLayer := handlers.NewHandler(serviceLayer)

	http.HandleFunc("/park", handlerLayer.Park)
	http.HandleFunc("/unpark", handlerLayer.UnPark)
	http.HandleFunc("/get-slots", handlerLayer.Available)

	fmt.Println("Server running on :8082")
	if err = http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
