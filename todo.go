package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/gin-contrib/cors"
)

var db *gorm.DB

type Todo struct {
	ID        uint   `json:"ID"`
	Title     string `json:"Title"`
	Completed bool   `json:"Completed"`
}

func getTodos(c *gin.Context) {
	var todos []Todo
	db.Find(&todos)
	log.Printf("Todos: %+v\n", todos) // Print the todos to the console
	c.JSON(http.StatusOK, todos)
}


func createTodo(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err == nil {
		db.Create(&todo)
		log.Printf("New todo created: %+v\n", todo) // Print the created todo to the console
		c.JSON(http.StatusCreated, todo)
	} else {
		log.Printf("Error creating todo: %s\n", err.Error()) // Print the error to the console
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func updateTodo(c *gin.Context) {
	var todo Todo
	id := c.Param("ID")
	if err := db.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	if err := c.ShouldBindJSON(&todo); err == nil {
		db.Save(&todo)
		c.JSON(http.StatusOK, todo)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func deleteTodo(c *gin.Context) {
	var todo Todo
	id := c.Param("ID")
	db.Delete(&todo, id)
	c.JSON(http.StatusOK, gin.H{"message": "todo deleted"})
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Todo{})

	r := gin.Default()

	// Enable CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.GET("/todos", getTodos)
	r.POST("/todos", createTodo)
	r.PUT("/todos/:ID", updateTodo)
	r.DELETE("/todos/:ID", deleteTodo)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err.Error())
	}
}
