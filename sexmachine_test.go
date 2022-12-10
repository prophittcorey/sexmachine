package sexmachine

import (
	"testing"
)

func TestParseName(t *testing.T) {
	if name := parsename(""); name != "" {
		t.Fatalf("failed to parse a blank name; got %s", name)
	}

	if name := parsename("   \t\n "); name != "" {
		t.Fatalf("failed to parse a blank name; got %s", name)
	}

	if name := parsename(" John Smith "); name != "john" {
		t.Fatalf("failed to parse a full name; got %s", name)
	}
}

func TestClassification(t *testing.T) {
	classifier := New()

	classifier.Train(Male, "bob", "john", "tim")
	classifier.Train(Female, "sara", "sally", "abbey")

	if sex, prob := classifier.Predict("Bob"); sex != Male {
		t.Fatalf("failed to classify Bob as Male; got %v at %f", sex, prob)
	}

	if sex, prob := classifier.Predict("Sara"); sex != Female {
		t.Fatalf("failed to classify Sara as Female; got %v at %f", sex, prob)
	}
}
