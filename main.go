package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

const outputFile = "lecture.mp3"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("⚠️  No .env file found, relying on environment variables")
	}

	recordCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	deviceID, err := getDefaultAudioOutput()
	if err != nil {
		log.Fatalf("failed to get audio source: %v", err)
	}

	if err := recordAudio(recordCtx, deviceID, outputFile); err != nil {
		log.Fatalf("recording failed: %v", err)
	}

	ctx := context.Background()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("please set OPENAI_API_KEY environment variable")
	}
	client := newOpenAIClient(apiKey)

	transcript, err := transcribeAudio(ctx, client, outputFile)
	if err != nil {
		log.Fatalf("transcription failed: %v", err)
	}

	content, err := analyzeLecture(ctx, client, transcript)
	if err != nil {
		log.Fatalf("lecture analysis failed: %v", err)
	}

	lectureOutput, err := parseLectureOutput(content)
	if err != nil {
		log.Fatalf("parsing output failed: %v", err)
	}

	if err := createNotionPage(ctx, *lectureOutput); err != nil {
		log.Printf("failed to create Notion page: %v", err)
	} else {
		fmt.Println("✅ Notion page created successfully!")
	}
}
