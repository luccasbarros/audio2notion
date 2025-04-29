package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jomei/notionapi"
)

func createNotionPage(ctx context.Context, output LectureOutput) error {
	notionToken := os.Getenv("NOTION_API_TOKEN")
	if notionToken == "" {
		return fmt.Errorf("🚨 NOTION_API_TOKEN environment variable not set")
	}

	databaseID := os.Getenv("NOTION_DATABASE_ID")
	if databaseID == "" {
		return fmt.Errorf("🚨 NOTION_DATABASE_ID environment variable not set")
	}

	notionClient := notionapi.NewClient(notionapi.Token(notionToken))

	blocks := buildLectureBlocks(output)

	_, err := notionClient.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: notionapi.DatabaseID(databaseID),
		},
		Properties: notionapi.Properties{
			"Name": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{Text: &notionapi.Text{Content: fmt.Sprintf("Aula - %s", time.Now().Format("2006-01-02"))}},
				},
			},
		},
		Children: blocks,
	})

	return err
}

func buildLectureBlocks(output LectureOutput) []notionapi.Block {
	var blocks []notionapi.Block

	// Lecture Analysis
	blocks = append(blocks, newHeading1("📚 Lecture Analysis"))

	// Main Topics
	blocks = append(blocks, newHeading2("Main Topics"))
	for _, topic := range output.LectureAnalysis.DetailedNotes.MainTopics {
		blocks = append(blocks, newBulletPoint(topic))
	}

	// Key Concepts
	blocks = append(blocks, newHeading2("Key Concepts"))
	for concept, details := range output.LectureAnalysis.DetailedNotes.KeyConcepts {
		blocks = append(blocks, newBulletPoint(fmt.Sprintf("%s: %s", concept, details.Explanation)))
		// related concepts
		if len(details.RelatedConcepts) > 0 {
			blocks = append(blocks, newParagraph("Related Concepts: "+joinList(details.RelatedConcepts)))
		}
		// additional resources
		if len(details.AdditionalResources) > 0 {
			blocks = append(blocks, newParagraph("Additional Resources: "+joinList(details.AdditionalResources)))
		}
	}

	// Technical Details
	blocks = append(blocks, newHeading2("Technical Details"))
	for _, detail := range output.LectureAnalysis.DetailedNotes.TechnicalDetails {
		blocks = append(blocks, newBulletPoint(detail))
	}

	// Important Examples
	blocks = append(blocks, newHeading2("Important Examples"))
	for _, example := range output.LectureAnalysis.DetailedNotes.ImportantExamples {
		blocks = append(blocks, newBulletPoint(example))
	}

	// Cited References
	blocks = append(blocks, newHeading2("Cited References"))
	for _, ref := range output.LectureAnalysis.DetailedNotes.CitedReferences {
		blocks = append(blocks, newBulletPoint(ref))
	}

	// Supplementary Information
	blocks = append(blocks, newHeading1("🧩 Supplementary Information"))
	sup := output.LectureAnalysis.SupplementaryInformation
	blocks = append(blocks, newParagraph("Topic: "+sup.Topic))
	blocks = append(blocks, newParagraph("External Resources: "+joinList(sup.ExternalResources)))
	blocks = append(blocks, newParagraph("Related Technologies: "+joinList(sup.RelatedTechnologies)))
	blocks = append(blocks, newParagraph("Industry Applications: "+joinList(sup.IndustryApplications)))

	// Exercises
	blocks = append(blocks, newHeading1("💪 Exercises"))
	blocks = append(blocks, buildExerciseSection("Basic", output.Exercises.Basic)...)
	blocks = append(blocks, buildExerciseSection("Medium", output.Exercises.Medium)...)
	blocks = append(blocks, buildExerciseSection("Advanced", output.Exercises.Advanced)...)

	// Vocabulary
	blocks = append(blocks, newHeading1("📖 Vocabulary"))
	blocks = append(blocks, buildVocabularyTable(output))

	// Study Recommendations
	blocks = append(blocks, newHeading1("📝 Study Recommendations"))
	rec := output.StudyRecommendations
	blocks = append(blocks, newParagraph("Key Points for Exam: "+joinList(rec.KeyPointsForExam)))
	blocks = append(blocks, newParagraph("Suggested Practice Areas: "+joinList(rec.SuggestedPracticeAreas)))
	blocks = append(blocks, newParagraph("Common Pitfalls: "+joinList(rec.CommonPitfalls)))
	blocks = append(blocks, newParagraph("Additional Reading: "+joinList(rec.AdditionalReading)))

	return blocks
}

func buildVocabularyTable(output LectureOutput) notionapi.Block {
	tableRows := []notionapi.Block{
		notionapi.TableRowBlock{
			BasicBlock: notionapi.BasicBlock{Object: "block", Type: "table_row"},
			TableRow: notionapi.TableRow{
				Cells: [][]notionapi.RichText{
					{{Text: &notionapi.Text{Content: "Termo (PT)"}}},
					{{Text: &notionapi.Text{Content: "English"}}},
					{{Text: &notionapi.Text{Content: "Definição"}}},
				},
			},
		},
	}

	for term, details := range output.Vocabulary.TechnicalTerms {
		tableRows = append(tableRows, notionapi.TableRowBlock{
			BasicBlock: notionapi.BasicBlock{Object: "block", Type: "table_row"},
			TableRow: notionapi.TableRow{
				Cells: [][]notionapi.RichText{
					{{Text: &notionapi.Text{Content: term}}},
					{{Text: &notionapi.Text{Content: details.EnglishTranslation}}},
					{{Text: &notionapi.Text{Content: details.TechnicalDefinition}}},
				},
			},
		})
	}

	return notionapi.TableBlock{
		BasicBlock: notionapi.BasicBlock{Object: "block", Type: "table"},
		Table: notionapi.Table{
			TableWidth:      3,
			HasColumnHeader: true,
			HasRowHeader:    false,
			Children:        tableRows,
		},
	}
}

func newHeading1(text string) notionapi.Block {
	return notionapi.Heading1Block{
		BasicBlock: notionapi.BasicBlock{Object: "block", Type: "heading_1"},
		Heading1: notionapi.Heading{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: text}},
			},
		},
	}
}

func newHeading2(text string) notionapi.Block {
	return notionapi.Heading2Block{
		BasicBlock: notionapi.BasicBlock{Object: "block", Type: "heading_2"},
		Heading2: notionapi.Heading{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: text}},
			},
		},
	}
}

func newBulletPoint(text string) notionapi.Block {
	return notionapi.BulletedListItemBlock{
		BasicBlock: notionapi.BasicBlock{Object: "block", Type: "bulleted_list_item"},
		BulletedListItem: notionapi.ListItem{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: text}},
			},
		},
	}
}

func newParagraph(text string) notionapi.Block {
	return notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{Object: "block", Type: "paragraph"},
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: text}},
			},
		},
	}
}

func joinList(items []string) string {
	return strings.Join(items, ", ")
}

func buildExerciseSection(level string, exercises []Exercise) []notionapi.Block {
	var blocks []notionapi.Block
	blocks = append(blocks, newHeading2(fmt.Sprintf("Level: %s", level)))

	if len(exercises) == 0 {
		blocks = append(blocks, newParagraph("No exercises available."))
		return blocks
	}

	for _, exercise := range exercises {
		blocks = append(blocks, newBulletPoint(exercise.Title))
		blocks = append(blocks, newParagraph(fmt.Sprintf("Objective: %s", exercise.Objective)))
		blocks = append(blocks, newParagraph("Steps: "+joinList(exercise.Steps)))
		blocks = append(blocks, newParagraph("Expected Outcome: "+exercise.ExpectedOutcome))
	}

	return blocks
}
