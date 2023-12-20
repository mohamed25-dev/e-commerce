package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetDbConnectionString() (string, error) {
	user := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Printf("Error parsing port number: %v\n", err)
		return "", err
	}

	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbName), nil
}

func PgNumericToFloat32(pgNumber pgtype.Numeric) (float32, error) {
	v, err := pgNumber.Value()
	if err != nil {
		fmt.Println("Something went wrong while converting pg numeric to float, err: ", err)
		return 0, err
	}

	stringVal := v.(string)
	val, err := strconv.ParseFloat(stringVal, 32)
	if err != nil {
		fmt.Println("Something went wrong while converting pg numeric to float, err: ", err)
		return 0, err
	}

	return float32(val), nil
}

func Float32ToPgNumeric(f float32) (pgtype.Numeric, error) {
	strValue := strconv.FormatFloat(float64(f), 'f', -1, 32)
	var val pgtype.Numeric
	//TODO: fix conversion issue which happens here.
	if err := val.Scan(strValue); err != nil {
		fmt.Println("scanning error, err: ", err)
		return pgtype.Numeric{}, err
	}

	return val, nil
}
