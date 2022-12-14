package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prophittcorey/sexmachine"
)

var (
	errNoFiles = fmt.Errorf("error: one or more files are required for training")
)

func main() {
	var classifier string
	var data string

	flag.StringVar(&classifier, "classifier", "", "a classifier to create (or update if it already exists)")
	flag.StringVar(&data, "data", "", "one or more files to train a classifier with (csv: name,sex,frequency")

	flag.Parse()

	if len(data) > 0 && len(classifier) > 0 {
		if err := train(classifier, data); err != nil {
			log.Fatal(err)
		}

		return
	}

	flag.Usage()
}

func train(classifierpath, glob string) error {
	files, err := filepath.Glob(glob)

	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errNoFiles
	}

	classifier := sexmachine.New()

	// TODO: load an existing classifier?

	for _, f := range files {
		(func() {
			f, err := os.Open(f)

			if err != nil {
				log.Fatal(err)
			}

			defer f.Close()

			reader := csv.NewReader(f)

			data, err := reader.ReadAll()

			if err != nil {
				log.Fatal(err)
			}

			sexes := map[byte]int{
				'f': sexmachine.Female,
				'm': sexmachine.Male,
			}

			for _, row := range data {
				if len(row) < 3 {
					continue
				}

				name, sex, frequency := row[0], strings.ToLower(row[1]), row[2]

				if len(sex) == 0 || (sex[0] != 'f' && sex[0] != 'm') {
					continue
				}

				if i, err := strconv.ParseInt(frequency, 10, 32); err == nil {
					classifier.Observe(sexes[sex[0]], name, int(i))
				}
			}
		})()
	}

	if err = classifier.SaveFile(classifierpath); err != nil {
		return err
	}

	return nil
}
