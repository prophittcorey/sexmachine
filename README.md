# Sex Machine

[![Go Reference](https://pkg.go.dev/badge/github.com/prophittcorey/sexmachine.svg)](https://pkg.go.dev/github.com/prophittcorey/sexmachine)

A zero dependency, pure Golang package for the training and classification of names.

## Package Usage

Training and saving a classifier.

```golang
import "github.com/prophittcorey/sexmachine"

classifier := New()

/* train */

classifier.Train(Male, "joey", "joey", "joey", "nick", "sam", "brent")
classifier.Train(Female, "tory", "tara", "joey", "sara", "joey")

/* save */

if err := classifier.SaveFile(os.TempDir() + "/sexmachine.classifier"); err != nil {
  t.Fatalf("failed to write test file; %s", err)
}
```

Loading and using a classifier.

```bash
import "github.com/prophittcorey/sexmachine"

classifier := New()

/* load */

if err := classifier.LoadFile(os.TempDir() + "/sexmachine.classifier"); err != nil {
  t.Fatalf("failed to load test file; %s", err)
}

/* test */

if sex, probability := classifier.Predict("Joey"); sex == sexmachine.Male {
  fmt.Printf("Was male w/%.2f probability", probability)
}
```

## License

The source code for this repository is licensed under the MIT license, which you can
find in the [LICENSE](LICENSE.md) file.
