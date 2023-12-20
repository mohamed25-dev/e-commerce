package main

import (
	"context"
	"ecommerce/transactions/utils"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("transactions/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr, err := utils.GetDbConnectionString()
	if err != nil {
		log.Fatal("could not get DB connection string, err: ", err)
	}

	fmt.Println(connStr)
	// Create a connection pool
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close(context.Background())

	pool.Exec(context.Background(), `delete from transactions where id != '1'`)
	pool.Exec(context.Background(), `delete from products where id != '1'`)
	pool.Exec(context.Background(), `delete from customers where id != '1'`)

	_, err = pool.Exec(context.Background(), `INSERT INTO transactions (id,"created_at", customer_id, product_id,quantity, total_price)
		VALUES
			('123', NOW(), '345', '4324', 8, 452), 
			('124', NOW(), '345', '5265', 5, 263),
			('125', NOW(), '344', '9152', 2, 610)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = pool.Exec(context.Background(), `INSERT INTO customers (id, customer_name)
		VALUES
			('123','Mohamed Mirghani'),
			('124','Omer Babiker'), 
			('125','Murtada Mirhgani');`,
	)

	if err != nil {
		log.Fatal(err)
	}

	_, err = pool.Exec(context.Background(), `INSERT INTO products (id,product_name,price) 
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
