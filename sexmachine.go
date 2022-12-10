package sexmachine

import "strings"

// ISO/IEC 5218 (see: https://en.wikipedia.org/wiki/ISO/IEC_5218)
const (
	Unknown = iota
	Male
	Female
)

// Classifier is used to store labeled data for classification.
type Classifier struct {
	Freqs map[int]map[string]float64 /* (Fe)Male -> { bob->2.0, ... } */
	Total int                        /* size of entire population/observations */
}

// Train labels data and adds it to the classifier.
func (c *Classifier) Train(label int, names ...string) {
	for _, name := range names {
		c.Freqs[label][name]++
		c.Total++
	}
}

// Predict takes a name and attempts to predict the person's sex using the classifier.
func (c *Classifier) Predict(name string) (int, float64) {
	return 0, 0.0
}

// New creates a new classifier.
func New() *Classifier {
	return &Classifier{
		Freqs: map[int]map[string]float64{
			Male:   {},
			Female: {},
		},
	}
}

func parsename(name string) string {
	if fields := strings.Fields(name); len(fields) > 0 {
		return normalize(fields[0])
	}

	return ""
}

func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
