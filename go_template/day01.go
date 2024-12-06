package main

import (
	"errors"
	"os"
	"sort"

	"github.com/charmbracelet/log"
	"github.com/jessevdk/go-flags"
)

type ParsedData struct {
}

func parseData(data []byte) (ParsedData, error) {
  return ParsedData{}, nil 
}


func part1Calculation(data ParsedData) {

}

func part2Calculation(data ParsedData) {

}

func main() {
	var opts struct {
		Part2 bool `short:"p" long:"part" description:"Whether to calculate for part 2"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("Cannot parse cli args", "err", err)
	}

	if opts.Part2 {
		log.Info("Part 2 of the problem")
	} else {
		log.Info("Part 1 of the problem")
	}

	log.Info("Reading data....")
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal("Cannot read the data from the file")
	}

	log.Info("Parsing data...")
	data, err := parseData(bytes)
	if err != nil {
		log.Fatal("Cannot parse data", "err", err)
	}

	log.Info("Calculating...")
	if opts.Part2 {
		part2Calculation(data)
	} else {
		part1Calculation(data)
	}
}
