package sexmachine

import "testing"

func TestClassification(t *testing.T) {
	classifier := New()

	classifier.Train(Male, "bob", "john", "tim")
	classifier.Train(Female, "sara", "sally", "")

	if sex, prob := classifier.Assume("Bob"); sex != Male {
		t.Fatalf("failed to classify Bob as Male; got %v at %f", sex, prob)
	}

	if sex, prob := classifier.Assume("Sara"); sex != Female {
		t.Fatalf("failed to classify Sara as Female; got %v at %f", sex, prob)
	}
}
