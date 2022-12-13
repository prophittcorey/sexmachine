// Package sexmachine enables you to easily train and classify names to sexes.
package sexmachine

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io"
	"log"
	"os"
	"strings"
)

const (
	Male = iota
	Female
	Unknown
)

const defaultProb = 0.000001

type sex = int

type names struct {
	freqs map[string]float64
	total float64
}

func (n names) GobEncode() ([]byte, error) {
	w := &bytes.Buffer{}

	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(n.freqs); err != nil {
		return nil, err
	}

	if err := encoder.Encode(n.total); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (n *names) GobDecode(buf []byte) error {
	decoder := gob.NewDecoder(bytes.NewBuffer(buf))

	if err := decoder.Decode(&n.freqs); err != nil {
		return err
	}

	return decoder.Decode(&n.total)
}

func (n names) probability(name string) float64 {
	if val, ok := n.freqs[name]; ok {
		return val / n.total
	}

	return defaultProb
}

// Classifier is used to store labeled data for classification.
type Classifier struct {
	data map[sex]*names
}

func (c Classifier) GobEncode() ([]byte, error) {
	w := &bytes.Buffer{}

	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(c.data); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (c *Classifier) GobDecode(buf []byte) error {
	return gob.NewDecoder(bytes.NewBuffer(buf)).Decode(&c.data)
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

// Load takes an io.Reader and decodes it. This can be used to read a
// classifier from a file, virtual file, bytes, etc.
func (c *Classifier) Load(reader io.Reader) error {
	gz, err := gzip.NewReader(reader)

	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(gz)

	if err := gz.Close(); err != nil {
		return err
	}

	return decoder.Decode(c)
}

// LoadFile wraps Load.
func (c *Classifier) LoadFile(file string) error {
	f, err := os.Open(file)

	if err != nil {
		return err
	}

	defer (func() {
		if err := f.Close(); err != nil {
			log.Printf("warning: failed to close file; %s", err)
		}
	})()

	return c.Load(f)
}

// Save serializes a classifier and writes it out to an io.Writer.
func (c Classifier) Save(writer io.Writer) error {
	gz := gzip.NewWriter(writer)

	if err := gob.NewEncoder(gz).Encode(&c); err != nil {
		return err
	}

	return gz.Close()
}

// SaveFile wraps Save.
func (c Classifier) SaveFile(file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer (func() {
		if err := f.Close(); err != nil {
			log.Printf("warning: failed to close file; %s", err)
		}
	})()

	return c.Save(f)
}

func (c Classifier) priors() []float64 {
	mprior := c.data[Male].total
	fprior := c.data[Female].total

	if total := c.data[Male].total + c.data[Female].total; total > 0 {
		mprior /= total
		fprior /= total
	}

	return []float64{mprior, fprior}
}

// Predict takes a name and attempts to predict the person's sex
// based on the data it was trained with.
func (c Classifier) Predict(name string) (sex, float64) {
	name = parsename(name)

	scores := []float64{0.0, 0.0}
	priors := c.priors()
	total := 0.0

	scores[Male] = priors[Male] * c.data[Male].probability(name)

	total += scores[Male]

	scores[Female] = priors[Female] * c.data[Female].probability(name)

	total += scores[Female]

	if total > 0 {
		scores[Male] /= total
		scores[Female] /= total
	}

	if scores[Male] > scores[Female] {
		return Male, scores[Male]
	}

	// NOTE: What about an "ambiguous" score?

	return Female, scores[Female]
}

// New creates a new classifier that is ready to be used.
func New() *Classifier {
	return &Classifier{
		data: map[sex]*names{
			Male: &names{
				freqs: map[string]float64{},
			},
			Female: &names{
				freqs: map[string]float64{},
			},
		},
	}
}

// Sex converts a sex's integer representation to a string.
func Sex(s sex) string {
	switch s {
	case Male:
		return "male"
	case Female:
		return "female"
	default:
		return "unknown"
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
