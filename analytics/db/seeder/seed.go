package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbUser     = "postgres"
	dbPassword = "password"
	dbName     = "analytics_db"
	dbHost     = "localhost"
	dbPort     = 5432
)

func main() {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`delete from transactions where id != '1'`)

	_, err = db.Exec(`INSERT INTO public.transactions (id,"created_at", customer_id, product_id,quantity, total_price)
		VALUES
			('123', NOW(), 345, 4324, 8, 452), 
			('124', NOW(), 345, 5265, 5, 263),
			('125', NOW(), 344, 9152, 2, 61)`)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transactions inserted successfully!")
}
