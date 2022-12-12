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

	classifier.Train(Male, "bob", "john", "tim", "tom", "alex", "alex", "joey")
	classifier.Train(Female, "sara", "sally", "abbey", "alex", "joey")

	if sex, prob := classifier.Predict("Alex"); sex != Male {
		t.Fatalf("failed to classify Alex as Male; got %v at %f", sex, prob)
	}

	if sex, prob := classifier.Predict("Sara"); sex != Female {
		t.Fatalf("failed to classify Sara as Female; got %v at %f", sex, prob)
	}
}

func TestLoad(t *testing.T) {
	classifier := New()

	if err := classifier.Load("testdata/classifier.test.bin"); err != nil {
		t.Fatalf("failed to load classifier; %s", err)
	}

	if sex, prob := classifier.Predict("Alex"); sex != Male {
		t.Fatalf("failed to classify Alex as Male; got %v at %f", sex, prob)
	}

	if sex, prob := classifier.Predict("Sara"); sex != Female {
		t.Fatalf("failed to classify Sara as Female; got %v at %f", sex, prob)
	}
}
