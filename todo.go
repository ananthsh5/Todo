package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin" // Import the Gin web framework
	"gorm.io/driver/sqlite" // Import the SQLite driver for GORM
	"gorm.io/gorm" // Import the GORM library
	"github.com/gin-contrib/cors" // Import the CORS middleware for Gin
)

var db *gorm.DB // Declare a global variable to hold the GORM database connection

// Define a struct representing the Todo model
type Todo struct {
	ID        uint   `json:"ID"` // ID field of type uint for identifying a Todo
	Title     string `json:"Title"` // Title field of type string for Todo title
	Completed bool   `json:"Completed"` // Completed field of type boolean to represent the completion status of the Todo
}

// Handler function to get all Todos from the database
func getTodos(c *gin.Context) {
	var todos []Todo
	db.Find(&todos)
	log.Printf("Todos: %+v\n", todos) // Print the todos to the console
	c.JSON(http.StatusOK, todos) // Return the list of todos as JSON response
}

// Handler function to create a new Todo
func createTodo(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err == nil {
		db.Create(&todo)
		log.Printf("New todo created: %+v\n", todo) // Print the created todo to the console
		c.JSON(http.StatusCreated, todo) // Return the created todo as JSON response with 201 status code
	} else {
		log.Printf("Error creating todo: %s\n", err.Error()) // Print the error to the console
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Return the error as JSON response with 400 status code
	}
}

// Handler function to update an existing Todo
func updateTodo(c *gin.Context) {
	var todo Todo
	id := c.Param("ID") // Get the Todo ID from the request parameter
	if err := db.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"}) // Return an error response if the Todo is not found
		return
	}
	if err := c.ShouldBindJSON(&todo); err == nil {
		db.Save(&todo)
		c.JSON(http.StatusOK, todo) // Return the updated todo as JSON response with 200 status code
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Return the error as JSON response with 400 status code
	}
}

// Handler function to delete a Todo
func deleteTodo(c *gin.Context) {
	var todo Todo
	id := c.Param("ID") // Get the Todo ID from the request parameter
	db.Delete(&todo, id) // Delete the Todo from the database
	c.JSON(http.StatusOK, gin.H{"message": "todo deleted"}) // Return a success message as JSON response with 200 status code
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{}) // Connect to the SQLite database
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Todo{}) // Auto migrate the Todo model

	r := gin.Default() // Initialize a new Gin router with the default middleware

	// Enable CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config)) // Use the CORS middleware with the specified configuration

	// Define the route handlers
	r.GET("/todos", getTodos) // Handle GET requests for fetching all Todos
	r.POST("/todos", createTodo) // Handle POST requests for creating a new Todo
	r.PUT("/todos/:ID", updateTodo) // Handle PUT requests for updating an existing Todo
	r.DELETE("/todos/:ID", deleteTodo) // Handle DELETE requests for deleting a Todo

	if err := http.ListenAndServe(":8080", r); err != nil { // Start the HTTP server on port 8080
		log.Fatal(err.Error()) // Log any errors that occur during server start-up
	}
}
