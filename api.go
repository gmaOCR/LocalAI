package main

import (
	"embed"
	"net/http"
	"strconv"

	llama "github.com/go-skynet/llama/go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed index.html
var indexHTML embed.FS

func api(l *llama.LLama, listenAddr string, threads int) error {
	app := fiber.New()
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(indexHTML),
		NotFoundFile: "index.html",
	}))
	/*
		curl --location --request POST 'http://localhost:8080/predict' --header 'Content-Type: application/json' --data-raw '{
		    "text": "What is an alpaca?",
		    "topP": 0.8,
		    "topK": 50,
		    "temperature": 0.7,
		    "tokens": 100
		}'
	*/
	// Serve the index.html file
	app.Get("/", func(c *fiber.Ctx) error {
		data, err := indexHTML.ReadFile("index.html")
		if err != nil {
			return err
		}
		return c.Send(data)
	})

	// Endpoint to generate the prediction
	app.Post("/predict", func(c *fiber.Ctx) error {
		// Get input data from the request body
		input := new(struct {
			Text string `json:"text"`
		})
		if err := c.BodyParser(input); err != nil {
			return err
		}

		// Set the parameters for the language model prediction
		topP, err := strconv.ParseFloat(c.Query("topP", "0.9"), 64) // Default value of topP is 0.9
		if err != nil {
			return err
		}

		topK, err := strconv.Atoi(c.Query("topK", "40")) // Default value of topK is 40
		if err != nil {
			return err
		}

		temperature, err := strconv.ParseFloat(c.Query("temperature", "0.5"), 64) // Default value of temperature is 0.5
		if err != nil {
			return err
		}

		tokens, err := strconv.Atoi(c.Query("tokens", "128")) // Default value of tokens is 128
		if err != nil {
			return err
		}

		// Generate the prediction using the language model
		prediction, err := l.Predict(
			input.Text,
			llama.SetTemperature(temperature),
			llama.SetTopP(topP),
			llama.SetTopK(topK),
			llama.SetTokens(tokens),
			llama.SetThreads(threads),
		)
		if err != nil {
			return err
		}

		// Return the prediction in the response body
		return c.JSON(struct {
			Prediction string `json:"prediction"`
		}{
			Prediction: prediction,
		})
	})

	// Start the server
	app.Listen(":8080")
	return nil
}
