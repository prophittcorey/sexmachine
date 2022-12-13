package main

import (
	"flag"
	"log"
)

func main() {
	var data string
	var classifier string

	flag.StringVar(&data, "data", "", "a data file to train a classifier with")
	flag.StringVar(&classifier, "classifier", "", "a classifier to create (or update if it already exists)")

	flag.Parse()

	if len(data) > 0 && len(classifier) > 0 {
		if err := train(classifier, data); err != nil {
			log.Fatal(err)
		}

		return
	}

	flag.Usage()
}

func train(classifier, data string) error {
	return nil
}
