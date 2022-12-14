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
	var classifierfile string
	var datafile string
	var testfile string
	var check string

	flag.StringVar(&classifierfile, "classifier", "", "a path to a classifier")
	flag.StringVar(&datafile, "data", "", "one or more files to train a classifier with (csv: name,sex,frequency")
	flag.StringVar(&testfile, "test", "", "a file to test the classifier against (csv: name,sex,frequency)")
	flag.StringVar(&check, "check", "", "a name to check using the specified classifier")

	flag.Parse()

	if len(datafile) > 0 && len(classifierfile) > 0 {
		if err := train(classifierfile, datafile); err != nil {
			log.Fatal(err)
		}

		return
	}

	if len(check) > 0 && len(classifierfile) > 0 {
		classifier := sexmachine.New()

		if err := classifier.LoadFile(classifierfile); err == nil {
			sex, prob := classifier.Predict(check)

			fmt.Printf("%s is %s (%.2f%%)\n", check, sexmachine.Sex(sex), prob*100)
		}

		return
	}

	if len(testfile) > 0 && len(classifierfile) > 0 {
		classifier := sexmachine.New()

		if err := classifier.LoadFile(classifierfile); err == nil {
			if err := test(classifier, testfile); err != nil {
				log.Fatal(err)
			}
		}

		return
	}

	flag.Usage()
}

func train(classifierfile, glob string) error {
	files, err := filepath.Glob(glob)

	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errNoFiles
	}

	classifier := sexmachine.New()

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

	if err = classifier.SaveFile(classifierfile); err != nil {
		return err
	}

	return nil
}

func test(classifier *sexmachine.Classifier, testfile string) error {
	f, err := os.Open(testfile)

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

	var total int
	var hit int
	var miss int

	for _, row := range data {
		if len(row) < 2 {
			continue
		}

		name, sex := row[0], strings.ToLower(row[1])

		if len(sex) == 0 || (sex[0] != 'f' && sex[0] != 'm') {
			continue
		}

		total++

		prediction, _ := classifier.Predict(name)

		if prediction == sexes[sex[0]] {
			hit++
		} else {
			miss++
		}
	}

	fmt.Printf("Got %d hits and %d misses (%.2f%% success).\n", hit, miss, (float64(hit)/float64(total))*100)

	return nil
}
