package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   error       `json:"error"`
}

type ToDo struct {
	ID          int    `json:"id"`
	Completed   bool   `json:"completed"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	// test
	var err error = godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var PORT string = os.Getenv("PORT")

	var app *fiber.App = fiber.New()

	var todos []ToDo = []ToDo{}

	app.Get("/", func(c *fiber.Ctx) error {
		var response Response = Response{}

		response.Success = true
		response.Message = "Welcome to the API"

		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Get("/api/todo", func(c *fiber.Ctx) error {
		var response Response = Response{}

		response.Success = true
		response.Message = "TODOS listed successfully"
		response.Data = todos

		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Post("/api/todo", func(c *fiber.Ctx) error {
		var response Response = Response{}
		var todo *ToDo = &ToDo{}

		if err := c.BodyParser(todo); err != nil {
			response.Success = false
			response.Message = "Invalid JSON | " + err.Error()
			response.Error = err

			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		if todo.Title == "" {
			response.Success = false
			response.Message = "Title is required"

			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		response.Success = true
		response.Message = "Todo created successfully"
		response.Data = *todo

		return c.Status(fiber.StatusCreated).JSON(response)
	})

	app.Patch("/api/todo/:id", func(c *fiber.Ctx) error {
		var response Response = Response{}

		id, err := c.ParamsInt("id")

		if err != nil {
			response.Success = false
			response.Message = "Invalid ID | " + err.Error()
			response.Error = err

			return c.Status(fiber.StatusBadRequest).JSON(response)
		} else if id < 1 || id > len(todos) {
			response.Success = false
			response.Message = "Todo not found"

			return c.Status(fiber.StatusNotFound).JSON(response)
		}

		// var todo *ToDo = &todos[id-1]
		var todo *ToDo = &ToDo{}

		if err := c.BodyParser(todo); err != nil {
			response.Success = false
			response.Message = "Invalid JSON | " + err.Error()
			response.Error = err

			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		for i, todo := range todos {
			if todo.ID == id {
				todos[i].Completed = todo.Completed
				todos[i].Title = todo.Title
				todos[i].Description = todo.Description

				response.Success = true
				response.Message = "Todo updated successfully"
				response.Data = todos[i]

				return c.Status(fiber.StatusOK).JSON(response)
			}
		}

		response.Success = false
		response.Message = "Todo not found"

		return c.Status(fiber.StatusNotFound).JSON(response)
	})

	app.Patch("/api/todo/:id/:complete", func(c *fiber.Ctx) error {
		var response Response = Response{}

		id, err := c.ParamsInt("id")

		if err != nil {
			response.Success = false
			response.Message = "Invalid ID | " + err.Error()
			response.Error = err

			return c.Status(fiber.StatusBadRequest).JSON(response)
		} else if id < 1 || id > len(todos) {
			response.Success = false
			response.Message = "Todo not found"

			return c.Status(fiber.StatusNotFound).JSON(response)
		}

		completed, err2 := strconv.ParseBool(c.Params("complete"))

		if err2 != nil {
			response.Success = false
			response.Message = "Invalid complete | " + err.Error()
			response.Error = err

			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		for i, todo := range todos {
			if todo.ID == id {
				todos[i].Completed = completed

				response.Success = true
				response.Message = "Todo updated successfully"
				response.Data = todos[i]

				return c.Status(fiber.StatusOK).JSON(response)
			}
		}

		response.Success = false
		response.Message = "Todo not found"

		return c.Status(fiber.StatusNotFound).JSON(response)
	})

	app.Delete("/api/todo/:id", func(c *fiber.Ctx) error {
		var response Response = Response{}

		id, err := c.ParamsInt("id")

		if err != nil {
			response.Success = false
			response.Message = "Invalid ID | " + err.Error()
			response.Error = err

			return c.Status(fiber.StatusBadRequest).JSON(response)
		} else if id < 1 || id > len(todos) {
			response.Success = false
			response.Message = "Todo not found"

			return c.Status(fiber.StatusNotFound).JSON(response)
		}

		for i, todo := range todos {
			if todo.ID == id {
				todos = append(todos[:i], todos[i+1:]...)

				response.Success = true
				response.Message = "Todo deleted successfully"

				return c.Status(fiber.StatusOK).JSON(response)
			}
		}

		response.Success = false
		response.Message = "Todo not found"

		return c.Status(fiber.StatusNotFound).JSON(response)
	})

	fmt.Println("App is running on http://127.0.0.1:" + PORT)
	log.Fatal(app.Listen("127.0.0.1:" + PORT))
}
