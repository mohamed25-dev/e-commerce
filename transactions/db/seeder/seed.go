package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// TODO: use env variables
const (
	dbUser     = "postgres"
	dbPassword = "password"
	dbName     = "transactions_db"
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
			('123', NOW(), '345', '4324', 8, 452), 
			('124', NOW(), '345', '5265', 5, 263),
			('125', NOW(), '344', '9152', 2, 610)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO public.customers (id, customer_name)
		VALUES
			('123','Mohamed Mirghani'),
			('124','Omer Babiker'), 
			('125','Murtada Mirhgani');`,
	)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO public.products (id,product_name,price) 
		VALUES
			('123','Pixel 8',1299.00),
			('124','Galaxy Tab s9', 1499.00),
			('125','One Plus 10', 1100.00);`,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data inserted successfully!")
}
