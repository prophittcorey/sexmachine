package sexmachine

// ISO/IEC 5218 (see: https://en.wikipedia.org/wiki/ISO/IEC_5218)
const (
	Unknown = iota
	Male
	Female
)

// Classifier is used to store labeled data for classification.
type Classifier struct {
}

// Train labels data and adds it to the classifier.
func (c *Classifier) Train(label int, names ...string) {
}

// Assume takes a name and attempts to guess its sex using the classifier.
func (c *Classifier) Assume(name string) (int, float64) {
	return 0, 0.0
}

// New creates a brand new classifier.
func New() *Classifier {
	return &Classifier{}
}
