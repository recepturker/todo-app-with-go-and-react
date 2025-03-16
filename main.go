package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   error       `json:"error"`
}

type ToDo struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Completed   bool               `json:"completed"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
}

var PORT string = "8001"
var client *mongo.Client
var collection *mongo.Collection
var app *fiber.App

func main() {
	setConfigurations()

	defer client.Disconnect(context.Background())

	app.Get("/", func(c *fiber.Ctx) error {
		var response Response = Response{}

		response.Success = true
		response.Message = "Welcome to the API"

		return c.Status(fiber.StatusOK).JSON(response)
	})
	app.Get("/api/todos", getTodos)
	// app.Get("/api/todo/:id", getTodo)
	app.Post("/api/todo", createTodo)
	app.Patch("/api/todo/:id/:complete", updateTodo)
	app.Delete("/api/todo/:id", deleteTodo)

	fmt.Println("App is running on http://127.0.0.1:" + PORT)
	log.Fatal(app.Listen("127.0.0.1:" + PORT))
}

func setConfigurations() {
	var err error = godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var TODOAPP_MONGODB_URI string = os.Getenv("TODOAPP_MONGODB_URI")
	var clientOptions *options.ClientOptions = options.Client().ApplyURI(TODOAPP_MONGODB_URI)
	client, err = mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database("todoapp").Collection("todos")

	PORT = os.Getenv("PORT")
	app = fiber.New()
}

func getTodos(c *fiber.Ctx) error {
	var todos []ToDo = []ToDo{}

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Success: false,
			Message: "Error fetching todos. Error: " + err.Error(),
			Data:    nil,
			Error:   err,
		})
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo ToDo

		if err := cursor.Decode(&todo); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(Response{
				Success: false,
				Message: "Error decoding todos. Error: " + err.Error(),
				Data:    nil,
				Error:   err,
			})
		}

		todos = append(todos, todo)
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: "Todos listed successfully",
		Data:    todos,
		Error:   nil,
	})
}

func createTodo(c *fiber.Ctx) error {
	var todo *ToDo = new(ToDo)
	// var todo ToDo

	if err := c.BodyParser(todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Success: false,
			Message: "Error parsing request body. Error: " + err.Error(),
			Data:    nil,
			Error:   err,
		})
	}

	if todo.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Success: false,
			Message: "Title is required",
			Data:    nil,
			Error:   nil,
		})
	}

	result, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Success: false,
			Message: "Error creating todo. Error: " + err.Error(),
			Data:    nil,
			Error:   err,
		})
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: "Todo created successfully",
		Data:    todo,
		Error:   nil,
	})
}

func updateTodo(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Success: false,
			Message: "Invalid ID. Error: " + err.Error(),
			Data:    nil,
			Error:   err,
		})
	}

	complete, err := strconv.ParseBool(c.Params("complete"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Success: false,
			Message: "Invalid complete value. Error: " + err.Error(),
			Data:    nil,
			Error:   err,
		})
	}

	update := bson.M{
		"$set": bson.M{
			"completed": complete,
		},
	}

	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Success: false,
			Message: "Error updating todo. Error: " + err.Error(),
			Data:    nil,
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: "Todo updated successfully",
		Data:    result,
		Error:   nil,
	})
}

func deleteTodo(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Success: false,
			Message: "Invalid ID. Error: " + err.Error(),
			Data:    nil,
			Error:   err,
		})
	}

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Success: false,
			Message: "Error deleting todo. Error: " + err.Error(),
			Data:    nil,
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: "Todo deleted successfully",
		Data:    result,
		Error:   nil,
	})
}
