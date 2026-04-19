package main

import (
	"database/sql",
	"log",
	"net/http",
	"os",
	"strconv",
	"github.com/gin-gonic/gin"
	 "github.com/lib/pq"
	"github.com/joho/godotenv"
)

type Todo struct{
	ID int `json:"id"`
	Task string `json:"task"`
	Done bool `json:"done"`
}

var db *sql.DB

func main(){
	err := godotenv.Load()
	if err ! := nil {
		log.fatal("error loading .env")
	}

	dbURL := os.Getenv("DB_URL")
	port := os.Getenv("PORT")

	db, err = sql.Open("postgres", dbURL)
	if err ! := nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err ! := nil {
		log.Fatal("database connection failed:", err)
	}

	r := gin.Default()

	r.GET("/todos", getTodos)
	r.POST("/todos", createTodo)
	r.PUT("/todos", updateTodo)
	r.DELETE("/todos", deleteTodo)

	log.Println("server running on:" + port)
	r.Run(":" + port)
}