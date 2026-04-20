package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

var db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.fatal("error loading .env")
	}

	dbURL := os.Getenv("DB_URL")
	port := os.Getenv("PORT")

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
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

func getTodos(c *gin.Context) {
	rows, err := db.Query("SELECT id, task, done FROM todos ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Task, &todo.Done)
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
	var todo Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := db.QueryRow(
		"INSERT INTO todos (task,done) VALUES ($1,$2) RETURNING id",
		todo.Task,
		todo.Done,
	).Scan(&todo.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func updateTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Input"})
		return
	}

	_, err := db.Exec(
		"UPDATE todos SET task=$1, done=$2 WHERE id=$3",
		todo.Task,
		todo.Done,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	todo.ID = id
	c.JSON(http.StatusOK, todo)
}

func deleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := db.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "todo deleted",
	})

}
