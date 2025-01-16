package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
)

// generateNaturalDescription returns the time elapsed for the operation
func generateNaturalDescription(ctx context.Context, client *api.Client, imgData api.ImageData) (time.Duration, error) {
	streamTrue := true
	req := &api.GenerateRequest{
		Model: "llama3.2-vision",
		// Prompt asks for a detailed description including objects, scene, colors, and mood
		Prompt: "请用中文详细描述这张图片的内容，包括图片中的主要对象、场景、颜色和整体氛围。描述长度控制在500字以内。",
		Images: []api.ImageData{imgData},
		Stream: &streamTrue,
		// System prompt ensures Chinese output with accurate and detailed description
		System: "你是一个专业的图像描述专家。请始终使用中文回答，描述要准确、清晰、富有细节，并严格控制在500字以内。",
	}

	startTime := time.Now()
	fmt.Println("Generating natural language description...")

	// Generate and stream the response
	err := client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to generate natural description: %v", err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nNatural description time elapsed: %v\n", elapsed)
	return elapsed, nil
}

// generateStructuredDescription returns the time elapsed for the operation
func generateStructuredDescription(ctx context.Context, client *api.Client, imgData api.ImageData) (time.Duration, error) {
	// ImageDescription defines the structure for JSON output
	type ImageDescription struct {
		MainObjects []string `json:"main_objects"` // Main objects in the image
		Scene       string   `json:"scene"`        // Scene description
		Colors      []string `json:"colors"`       // Main colors in the image
		Mood        string   `json:"mood"`         // Overall mood/atmosphere
		Details     string   `json:"details"`      // Additional details
	}

	streamTrue := true
	req := &api.GenerateRequest{
		Model: "llama3.2-vision",
		// Prompt specifies JSON format with required fields
		Prompt: `请以 JSON 格式描述这张图片，包含以下字段：
- main_objects: 图片中的主要对象（数组，最多5个）
- scene: 场景描述（100字以内）
- colors: 主要颜色（数组，最多5个）
- mood: 整体氛围（50字以内）
- details: 额外细节描述（200字以内）
请确保输出是有效的 JSON 格式，总字数控制在500字以内。`,
		Images: []api.ImageData{imgData},
		Stream: &streamTrue,
		// System prompt enforces JSON structure and Chinese output
		System: `你是一个专业的图像分析专家。请始终使用中文回答，并严格按照以下 JSON 格式输出，确保总字数不超过500字：
{
    "main_objects": ["对象1", "对象2"],
    "scene": "场景描述",
    "colors": ["颜色1", "颜色2"],
    "mood": "氛围描述",
    "details": "细节描述"
}`,
	}

	startTime := time.Now()
	fmt.Println("Generating structured description...")

	// Collect full response for JSON parsing
	var fullResponse string
	err := client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		fullResponse += resp.Response
		fmt.Print(resp.Response)
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to generate structured description: %v", err)
	}

	// Validate JSON format
	var desc ImageDescription
	if err := json.NewDecoder(strings.NewReader(fullResponse)).Decode(&desc); err != nil {
		return 0, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nStructured description time elapsed: %v\n", elapsed)
	return elapsed, nil
}

func main() {
	totalStart := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	imgData, err := os.ReadFile("test.jpg")
	if err != nil {
		log.Fatalf("failed to read image file: %v", err)
	}

	baseURL, err := url.Parse("http://localhost:11434")
	if err != nil {
		log.Fatalf("failed to parse URL: %v", err)
	}
	client := api.NewClient(baseURL, &http.Client{Timeout: 30 * time.Second})

	fmt.Println("\n=== Natural Language Description ===")
	naturalTime, err := generateNaturalDescription(ctx, client, imgData)
	if err != nil {
		log.Printf("natural description failed: %v", err)
	}

	fmt.Println("\n=== Structured Description ===")
	structuredTime, err := generateStructuredDescription(ctx, client, imgData)
	if err != nil {
		log.Printf("structured description failed: %v", err)
	}

	// Print total time and breakdown
	totalTime := time.Since(totalStart)
	fmt.Printf("\n=== Time Summary ===\n")
	fmt.Printf("Natural Description: %v\n", naturalTime)
	fmt.Printf("Structured Description: %v\n", structuredTime)
	fmt.Printf("Total Time: %v\n", totalTime)
}
