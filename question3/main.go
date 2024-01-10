package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

func main() {
	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", "user=postgres password=123 dbname=sorting sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	fmt.Println("Connected to the database successfully")

	// Create a temporary table with updated seat IDs
	_, err = db.Exec(`
		CREATE TEMPORARY TABLE IF NOT EXISTS temp_seat AS
		SELECT 
			CASE 
				WHEN id % 2 = 0 THEN id - 1
				ELSE id
			END AS id, 
			student
		FROM Seat
		ORDER BY id
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Replace the original table with the temporary table
	_, err = db.Exec("ALTER TABLE temp_seat RENAME TO Seat")
	if err != nil {
		log.Fatal(err)
	}

	// Query the result table ordered by id in ascending order
	rows, err := db.Query("SELECT * FROM Seat ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Print the result
	fmt.Printf("+----+---------+\n| id | student |\n+----+---------+\n")
	rowNumber := 1
	for rows.Next() {
		var id int
		var student string
		if err := rows.Scan(&id, &student); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("| %d | %s |\n", rowNumber, student)
		rowNumber++
	}
	fmt.Printf("+----+---------+\n")
}
