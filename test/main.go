package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ollama/ollama/api"
)

func main() {
	imgData, err := os.ReadFile("test.jpg")
	if err != nil {
		log.Fatal(err)
	}

	baseURL, err := url.Parse("http://localhost:11434")
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewClient(baseURL, http.DefaultClient)

	req := &api.GenerateRequest{
		Model:  "llama3.2-vision",
		Prompt: "describe this image",
		Images: []api.ImageData{imgData},
	}

	ctx := context.Background()
	startTime := time.Now()

	respFunc := func(resp api.GenerateResponse) error {
		// In streaming mode, responses are partial so we call fmt.Print (and not
		// Println) in order to avoid spurious newlines being introduced. The
		// model will insert its own newlines if it wants.
		fmt.Print(resp.Response)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nGeneration took: %v\n", elapsed)
}
