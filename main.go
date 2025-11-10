package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://postgres:admin@db:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var now string
		err := db.QueryRow("SELECT NOW()").Scan(&now)
		if err != nil {
			http.Error(w, "DB Error", 500)
			return
		}
		fmt.Fprintf(w, "Hello from Go! Database time: %s\n", now)
	})

	log.Println("Server running on port 9091...")
	log.Fatal(http.ListenAndServe(":9091", nil))
}
