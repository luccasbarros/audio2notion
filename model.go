package main

type LectureOutput struct {
	LectureAnalysis struct {
		DetailedNotes struct {
			MainTopics  []string `json:"main_topics"`
			KeyConcepts map[string]struct {
				Explanation         string   `json:"explanation"`
				RelatedConcepts     []string `json:"related_concepts"`
				AdditionalResources []string `json:"additional_resources"`
			} `json:"key_concepts"`
			TechnicalDetails  []string `json:"technical_details"`
			ImportantExamples []string `json:"important_examples"`
			CitedReferences   []string `json:"cited_references"`
		} `json:"detailed_notes"`
		SupplementaryInformation struct {
			Topic                string   `json:"topic"`
			ExternalResources    []string `json:"external_resources"`
			RelatedTechnologies  []string `json:"related_technologies"`
			IndustryApplications []string `json:"industry_applications"`
		} `json:"supplementary_information"`
	} `json:"lecture_analysis"`
	Exercises struct {
		Basic    []Exercise `json:"basic"`
		Medium   []Exercise `json:"medium"`
		Advanced []Exercise `json:"advanced"`
	} `json:"exercises"`
	Vocabulary struct {
		TechnicalTerms map[string]struct {
			EnglishTranslation  string   `json:"english_translation"`
			TechnicalDefinition string   `json:"technical_definition"`
			UsageContext        string   `json:"usage_context"`
			RelatedTerms        []string `json:"related_terms"`
		} `json:"technical_terms"`
	} `json:"vocabulary"`
	StudyRecommendations struct {
		KeyPointsForExam       []string `json:"key_points_for_exam"`
		SuggestedPracticeAreas []string `json:"suggested_practice_areas"`
		CommonPitfalls         []string `json:"common_pitfalls"`
		AdditionalReading      []string `json:"additional_reading"`
	} `json:"study_recommendations"`
}

type Exercise struct {
	Title              string   `json:"title"`
	Objective          string   `json:"objective"`
	Prerequisites      []string `json:"prerequisites"`
	Steps              []string `json:"steps"`
	ResourcesNeeded    []string `json:"resources_needed"`
	ExternalReferences []string `json:"external_references"`
	ExpectedOutcome    string   `json:"expected_outcome"`
	ValidationCriteria string   `json:"validation_criteria"`
}
