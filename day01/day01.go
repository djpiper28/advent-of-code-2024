package main

import (
	"errors"
	"os"
	"sort"

	"github.com/charmbracelet/log"
	"github.com/jessevdk/go-flags"
	"golang.org/x/exp/constraints"
)

type ParsedData struct {
	List1, List2 []int
}

func parseData(data []byte) (ParsedData, error) {
	list1 := make([]int, 0)
	list2 := make([]int, 0)

	isListA := true
	var temp int

	for _, b := range data {
		if b == '\n' {
			if isListA {
				return ParsedData{}, errors.New("Illegal input - expected another number")
			}
			list2 = append(list2, temp)
			isListA = true
			temp = 0
			// Theree is a potential edge case below where data is `0     123123`
		} else if b == ' ' && temp != 0 {
			list1 = append(list1, temp)
			isListA = false
			temp = 0
		} else if b >= '0' && b <= '9' {
			temp *= 10
			temp += int(b - '0')
		}
	}

	if len(list2)+1 == len(list1) {
		list2 = append(list2, temp)
	}

	log.Info("Parsed data", "list1 length", len(list1), "list2 length", len(list2))

	return ParsedData{List1: list1, List2: list2}, nil
}

type compareFunc = func(i, j int) bool

func intCompare(data []int) compareFunc {
	return func(i, j int) bool {
		return data[i] < data[j]
	}
}

func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func part1Calculation(data ParsedData) {
	sum := 0

	for i, a := range data.List1 {
		b := data.List2[i]

		difference := Abs(a - b)
		sum += difference
	}

	log.Info("Complete!", "output", sum)
}

func part2Calculation(data ParsedData) {
	countMap := make(map[int]int)

	countArray := func(data []int) {
		for _, num := range data {
			count, found := countMap[num]
			if !found {
				countMap[num] = 1
			} else {
				countMap[num] = count + 1
			}
		}
	}

	countArray(data.List2)

	sum := 0

	for _, num := range data.List1{
    count, _ := countMap[num]
		difference := num * count
		sum += difference
	}

	log.Info("Complete!", "output", sum)
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

	log.Info("Sorting lists...")
	sort.Slice(data.List1, intCompare(data.List1))
	sort.Slice(data.List2, intCompare(data.List2))

	if opts.Part2 {
		part2Calculation(data)
	} else {
		part1Calculation(data)
	}
}
