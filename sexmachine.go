package sexmachine

import "strings"

const (
	Male = iota
	Female
)

const defaultProb = 0.00000000001

type sex = int

type names struct {
	freqs map[string]float64
	total float64
}

func (n names) probability(name string) float64 {
	if val, ok := n.freqs(name); ok {
		return value / n.total
	}

	return defaultProb
}

// Classifier is used to store labeled data for classification.
type Classifier struct {
	data map[sex]names
}

func (c *Classifier) priors() []float64 {
	mprior := c.data[Male].total
	fprior := c.data[Female].total

	if total := c.data[Male].total + c.data[Female].total; total > 0 {
		mprior /= total
		fprior /= total
	}

	return []float64{mprior, fprior}
}

// Train labels data and adds it to the classifier.
func (c *Classifier) Train(label sex, names ...string) {
	for _, name := range names {
		if d, ok := c.data[label]; ok {
			d.freqs[normalize(name)]++
			d.total++
		}
	}
}

// Predict takes a name and attempts to predict the person's sex
// based on the data it was trained with.
func (c *Classifier) Predict(name string) (sex, float64) {
	name = parsename(name)

	return 0, 0.0
}

// New creates a new classifier that is ready to be used.
func New() *Classifier {
	return &Classifier{
		data: map[sex]names{
			Male: names{
				freqs: map[string]float64{},
			},
			Female: names{
				freqs: map[string]float64{},
			},
		},
	}
}

// A helper function to naively parse a person's name. Assumes if there's more
// than one name given, the first whole word is the person's given name. The
// name is also normalized before being returned.
func parsename(name string) string {
	if fields := strings.Fields(name); len(fields) > 0 {
		return normalize(fields[0])
	}

	return ""
}

// A helper function to normalize a name (removed leading and trailing white
// space and downcases the name).
func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
