package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"otpapi/api/handlers"
	"otpapi/db"
)

func testDBConnection(connString string) error {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("Unable to connect to the database: %v", err)
	}
	defer pool.Close()

	var result string
	err = pool.QueryRow(context.Background(), "SELECT 'Database is connected'").Scan(&result)
	if err != nil {
		return fmt.Errorf("Error executing query: %v", err)
	}

	fmt.Println(result)
	return nil
}

func main() {
	connString := "user=postgres password=123 dbname=otpapi sslmode=disable"

	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Unable to connect to the database:", err)
	}
	defer pool.Close()


	if err := testDBConnection(connString); err != nil {
		log.Fatal("Database connection test failed:", err)
	}
	
	dbAccess := db.NewDB(pool)
	
	router := gin.Default()
	
	router.POST("/api/users", handlers.CreateUser(dbAccess))
	router.POST("/api/users/generateotp", handlers.GenerateOTP(dbAccess))
	router.POST("/api/users/verifyotp", handlers.VerifyOTP(dbAccess))
	
	port := 8080
	serverAddr := fmt.Sprintf(":%d", port)
	log.Printf("Server is running on http://127.0.0.1%s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}