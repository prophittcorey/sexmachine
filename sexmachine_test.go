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

func TestUnknowns(t *testing.T) {
	classifier := New()

	classifier.Train(Male, "bob", "john", "tim", "tom", "alex", "alex", "joey")
	classifier.Train(Female, "sara", "sally", "abbey", "alex", "joey")

	if sex, prob := classifier.Predict("Maximus"); sex != Unknown {
		t.Fatalf("failed to classify an unknown; got %v at %f", Sex(sex), prob)
	}
}

func TestEqualProbability(t *testing.T) {
	classifier := New()

	classifier.Train(Male, "tim", "tom", "alex", "tim", "joey")
	classifier.Train(Female, "sara", "sally", "abbey", "alex", "joey")

	if sex, prob := classifier.Predict("Alex"); sex != Male {
		t.Fatalf("failed to classify Alex as Male with equal probabilities; got %v at %f", Sex(sex), prob)
	}
}

func TestClassification(t *testing.T) {
	classifier := New()

	classifier.Train(Male, "bob", "john", "tim", "tom", "alex", "alex", "joey")
	classifier.Train(Female, "sara", "sally", "abbey", "alex", "joey")

	if sex, prob := classifier.Predict("Alex"); sex != Male {
		t.Fatalf("failed to classify Alex as Male; got %v at %f", Sex(sex), prob)
	}

	if sex, prob := classifier.Predict("Sara"); sex != Female {
		t.Fatalf("failed to classify Sara as Female; got %v at %f", Sex(sex), prob)
	}
}

func TestClassificationWithObserve(t *testing.T) {
	classifier := New()

	classifier.Observe(Male, "Alex", 100)
	classifier.Observe(Male, "Bob", 10)
	classifier.Observe(Male, "Sam", 4)

	classifier.Observe(Female, "Alex", 5)
	classifier.Observe(Female, "Sally", 100)
	classifier.Observe(Female, "Same", 1)

	if sex, prob := classifier.Predict("Alex"); sex != Male {
		t.Fatalf("failed to classify Alex as Male; got %v at %f", Sex(sex), prob)
	}

	if sex, prob := classifier.Predict("Sally"); sex != Female {
		t.Fatalf("failed to classify Sara as Female; got %v at %f", Sex(sex), prob)
	}
}

func TestSuite(t *testing.T) {
	classifier := New()

	/* train */

	classifier.Train(Male, "joey", "joey", "joey", "nick", "sam", "brent")
	classifier.Train(Female, "tory", "tara", "joey", "sara", "joey")

	/* save */

	if err := classifier.SaveFile("/tmp/classifier-test.bin"); err != nil {
		t.Fatalf("failed to write test file; %s", err)
	}

	/* load */

	classifier = New()

	if err := classifier.LoadFile("/tmp/classifier-test.bin"); err != nil {
		t.Fatalf("failed to load test file; %s", err)
	}

	/* test */

	if sex, prob := classifier.Predict("Joey"); sex != Male {
		t.Fatalf("failed to classify Joey as Male; got %v at %f", Sex(sex), prob)
	}

	if sex, prob := classifier.Predict("tara"); sex != Female {
		t.Fatalf("failed to classify Tara as Female; got %v at %f", Sex(sex), prob)
	}
}
