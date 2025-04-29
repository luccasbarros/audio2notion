package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func newOpenAIClient(apiKey string) *openai.Client {
	return openai.NewClient(apiKey)
}

func transcribeAudio(ctx context.Context, client *openai.Client, filepath string) (string, error) {
	fmt.Println("📝 Transcribing with Whisper…")
	resp, err := client.CreateTranscription(ctx, openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filepath,
	})
	if err != nil {
		return "", err
	}
	fmt.Println("\n✅ Transcript received (preview):")
	fmt.Println(resp.Text[:min(len(resp.Text), 200)], "…")
	return resp.Text, nil
}

func analyzeLecture(ctx context.Context, client *openai.Client, transcript string) (string, error) {
	systemPrompt := `You are an expert educational analyst and technical documentation specialist. When analyzing lecture transcripts in Portuguese, provide a comprehensive response in JSON format that follows this structure:
{
  "lecture_analysis": {
    "detailed_notes": {
      "main_topics": [...],
      "key_concepts": {
        "concept_name": {
          "explanation": "...",
          "related_concepts": [...],
          "additional_resources": [...]
        }
      },
      "technical_details": [...],
      "important_examples": [...],
      "cited_references": [...]
    },
    "supplementary_information": {
      "topic": "...",
      "external_resources": [...],
      "related_technologies": [...],
      "industry_applications": [...]
    }
  },
  "exercises": {
    "basic": [
      {
        "title": "...",
        "objective": "...",
        "prerequisites": [...],
        "steps": [...],
        "resources_needed": [...],
        "external_references": [...],
        "expected_outcome": "...",
        "validation_criteria": "..."
      }
    ],
    "medium": [
      // Same structure as basic
    ],
    "advanced": [
      // Same structure as basic
    ]
  },
  "vocabulary": {
    "technical_terms": {
      "portuguese_term": {
        "english_translation": "...",
        "technical_definition": "...",
        "usage_context": "...",
        "related_terms": [...]
      }
    }
  },
  "study_recommendations": {
    "key_points_for_exam": [...],
    "suggested_practice_areas": [...],
    "common_pitfalls": [...],
    "additional_reading": [...]
  }
}

Important guidelines:
1. Analyze the content as if preparing comprehensive study materials for an advanced technical exam
2. Include all technical details, formulas, and methodologies mentioned
3. Add relevant external information that complements the lecture content
4. For exercises, provide detailed step-by-step instructions with real-world applications
5. Include links to documentation, tutorials, and relevant external resources
6. Explain complex concepts with practical examples
7. Identify and elaborate on any industry best practices mentioned
8. Cross-reference related concepts and terms
9. Add technical context to vocabulary terms
10. Highlight any critical concepts that require special attention

Remember to maintain technical accuracy and provide practical, actionable information that would be valuable for both study and real-world application.`

	userPrompt := fmt.Sprintf("Lecture transcript (in Portuguese):\n%s", transcript)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4oMini,
		Temperature: 0.7,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
	})
	if err != nil {
		return "", err
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	return content, nil
}

func parseLectureOutput(jsonContent string) (*LectureOutput, error) {
	var out LectureOutput
	if err := json.Unmarshal([]byte(jsonContent), &out); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return &out, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
