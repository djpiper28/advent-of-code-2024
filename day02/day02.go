package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jessevdk/go-flags"
	"golang.org/x/exp/constraints"
)

type ParsedData struct {
	Data []int
}

const UNSET = -1

func parseData(data []byte) ([]ParsedData, error) {
	parsedLines := make([]ParsedData, 0)
	currentLine := make([]int, 0)
	currentNum := UNSET

	appendCurrentNum := func() {
		if currentNum != UNSET {
			currentLine = append(currentLine, currentNum)
			currentNum = UNSET
		}
	}

	appendCurrentLine := func() {
		appendCurrentNum()
		parsedLines = append(parsedLines, ParsedData{Data: currentLine})
		currentLine = make([]int, 0)
	}

	for _, b := range data {
		if b >= '0' && b <= '9' {
			if currentNum == UNSET {
				currentNum = 0
			}

			currentNum *= 10
			currentNum += int(b - '0')
		} else if b == ' ' {
			appendCurrentNum()
		} else if b == '\n' {
			appendCurrentLine()
		} else {
			return nil, fmt.Errorf("Cannot parse %c - unrecognised char", b)
		}
	}

	if currentNum != UNSET {
		appendCurrentLine()
	}

	log.Info("Finished parsing data", "len", len(parsedLines))

	return parsedLines, nil
}

func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func validateLineIsSafe(input ParsedData) error {
	isUp := true

	for i, currentNum := range input.Data {
		if i == 1 {
			isUp = (currentNum - input.Data[0]) > 0
		}

		if i > 0 {
			if currentNum == input.Data[i-1] {
				return errors.New("Equal values are 'unsafe'")
			} else {
				difference := currentNum - input.Data[i-1]
				localIsUp := difference > 0

				if localIsUp != isUp {
					return errors.New("Differing isUp values are 'unsafe'")
				}

				if Abs(difference) > 3 {
					return errors.New("Large differences (>3) are 'unsafe'")
				}
			}
		}
	}

	return nil
}

func calculatePart1(data []ParsedData) int {
	sum := 0

	for _, line := range data {
		err := validateLineIsSafe(line)
		if err == nil {
			sum += 1
		}
	}

	return sum
}

func calculateLine2(line ParsedData) error {
	err := validateLineIsSafe(line)
	if err == nil {
		return nil
	}

	log.Debug("Unsafe line detected, trying to brute force removing one entity")
	for i := range line.Data {
		newLine := make([]int, 0)

		for j, num := range line.Data {
			if j != i {
				newLine = append(newLine, num)
			}
		}

		err = validateLineIsSafe(ParsedData{Data: newLine})
		if err == nil {
			return nil
		}
	}

	return errors.New("Cannot find a 'safe' variant")
}

func calculatePart2(data []ParsedData) int {
	sum := 0

	for _, line := range data {
		err := calculateLine2(line)
		if err == nil {
			sum += 1
		}
	}

	return sum
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

	log.Info("Reading file...")
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal("Cannot read file", "err", err)
	}

	log.Info("Parsing data...")
	data, err := parseData(bytes)
	if err != nil {
		log.Fatal("Cannot parse data", "err", err)
	}

	log.Info("Processing data...")
	if opts.Part2 {
		output := calculatePart2(data)
		log.Info("Processing done", "output", output)
	} else {
		output := calculatePart1(data)
		log.Info("Processing done", "output", output)
	}
}
