package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getDefaultAudioOutput() (string, error) {
	cmd := exec.Command("pactl", "get-default-sink")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get default sink: %v", err)
	}

	sinkName := strings.TrimSpace(string(output))
	monitorCmd := exec.Command("pactl", "list", "sources")
	monitorOutput, err := monitorCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list sources: %v", err)
	}

	sources := strings.Split(string(monitorOutput), "\n")
	var sourceIndex string
	for i, line := range sources {
		if strings.Contains(line, sinkName+".monitor") {
			for j := i; j >= 0; j-- {
				if strings.HasPrefix(sources[j], "Source #") {
					sourceIndex = strings.TrimPrefix(sources[j], "Source #")
					break
				}
			}
			break
		}
	}
	if sourceIndex == "" {
		return "", fmt.Errorf("could not find monitor for default sink")
	}

	return sourceIndex, nil
}

func recordAudio(ctx context.Context, deviceID string, outputFile string) error {
	fmt.Println("🎧 Recording to", outputFile, "… Press Ctrl+C to stop.")
	cmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-y",
		"-f", "pulse",
		"-i", deviceID,
		"-ar", "44100",
		"-ac", "1",
		"-codec:a", "libmp3lame",
		"-q:a", "4",
		"-b:a", "192k",
		outputFile,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil && ctx.Err() != context.Canceled {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}
