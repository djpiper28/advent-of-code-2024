package main

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/jessevdk/go-flags"
)

type Sequence struct {
	Data []int
}

type ParsedData struct {
	// Maps X|Y where Y relies on X (or as wordered x before y)
	// i.e:
	// NumberDependancies[y] = { x }
	NumberDependancies map[int][]int
	Sequences          []Sequence
}

const (
	DEPENDANT_SEPERATOR = '|'
	SEQUENCE_SEPERATOR  = ','
)

type ParsingState struct {
	Data  []byte
	Index int
}

func newParsingState(data []byte) *ParsingState {
	return &ParsingState{
		Data:  data,
		Index: 0,
	}
}

func (p *ParsingState) isEof() bool {
	return p.Index >= len(p.Data)
}

func (p *ParsingState) current() byte {
	return p.Data[p.Index]
}

func (p *ParsingState) next() bool {
	p.Index++
	return p.isEof()
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (p *ParsingState) scanNumber() (int, error) {
	if !isDigit(p.current()) {
		return 0, errors.New("Not a digit")
	}

	number := 0
	for isDigit(p.current()) && !p.isEof() {
		number *= 10
		number += int(p.current() - '0')
		p.next()
	}

	return number, nil
}

func (p *ParsingState) scanNewLine() error {
	if p.current() != '\n' {
		return fmt.Errorf("Expected a new line, found %c", p.current())
	}

	p.next()
	return nil
}

func (p *ParsingState) scanDependantSeperator() error {
	if p.current() != DEPENDANT_SEPERATOR {
		return fmt.Errorf("Expected a digit seperator, found %c", p.current())
	}

	p.next()
	return nil
}

func (p *ParsingState) scanSequenceSeperator() error {
	if p.current() != SEQUENCE_SEPERATOR {
		return fmt.Errorf("Expected a sequence seperator, found %c", p.current())
	}

	p.next()
	return nil
}

func (p *ParsingState) scanSequence() (Sequence, error) {
	sequence := make([]int, 0)

	for !p.isEof() {
		num, err := p.scanNumber()
		if err != nil {
			log.Error("Expected a digit", "err", err)
			return Sequence{}, err
		}

		sequence = append(sequence, num)

		if p.isEof() {
			break
		}

		if err = p.scanSequenceSeperator(); err != nil {
			if err = p.scanNewLine(); err == nil {
				break
			} else {
				log.Error("Expected a comma or a new line")
				return Sequence{}, errors.New("Expected a comma or a new line")
			}
		}
	}

	return Sequence{Data: sequence}, nil
}

func parseData(data []byte) (ParsedData, error) {
	state := newParsingState(data)
	ret := ParsedData{
		NumberDependancies: make(map[int][]int),
		Sequences:          make([]Sequence, 0),
	}
	totalDependancies := 0

	// Parse rules
	for !state.isEof() {
		// scan empty line - rules and sequence seperator
		if err := state.scanNewLine(); err == nil {
			break
		}

		dependant, err := state.scanNumber()
		if err != nil {
			log.Error("Expected a digit", "err", err)
			return ret, err
		}

		if err = state.scanDependantSeperator(); err != nil {
			log.Error("Expected a digit seperator", "err", err)
			return ret, err
		}

		num, err := state.scanNumber()
		if err != nil {
			log.Error("Expected a digit", "err", err)
			return ret, err
		}

		if err = state.scanNewLine(); err != nil {
			log.Error("Expected a new line", "err", err)
			return ret, err
		}

		deps, found := ret.NumberDependancies[num]
		if found {
			deps = append(deps, dependant)
			ret.NumberDependancies[num] = deps
		} else {
			list := make([]int, 0)
			list = append(list, dependant)
			ret.NumberDependancies[num] = list
		}

		totalDependancies++
	}

	for !state.isEof() {
		sequence, err := state.scanSequence()

		if err != nil {
			log.Error("Cannot scan sequence", "err", err)
			return ret, err
		}

		ret.Sequences = append(ret.Sequences, sequence)
	}

	log.Info("Parsing meta data",
		"NumberDependancies", len(ret.NumberDependancies),
		"totalDependancies", totalDependancies,
		"Sequences", len(ret.Sequences))

	return ret, nil
}

func testSequence(sequence *Sequence, data *ParsedData) bool {
	seenNumbers := make(map[int]bool)
	allNumbers := make(map[int]bool)

	for _, num := range sequence.Data {
		allNumbers[num] = true
	}

	for _, num := range sequence.Data {
		seenNumbers[num] = true

		deps, hasRules := data.NumberDependancies[num]
		if hasRules {
			for _, dep := range deps {
				_, isInDoc := allNumbers[dep]
				if !isInDoc {
					continue
				}

				_, depFound := seenNumbers[dep]
				if !depFound {
					return false
				}
			}
		}
	}

	return true
}

func part1Calculation(data ParsedData) {
	total := 0
	validSequences := 0

	for _, sequence := range data.Sequences {
		if testSequence(&sequence, &data) {
			total += sequence.Data[len(sequence.Data)/2]
			validSequences++
		}
	}

	log.Info("Complete!", "output", total, "validvalidSequences", validSequences)
}

func contains(arr []int, num int) bool {
	for _, v := range arr {
		if v == num {
			return true
		}
	}

	return false
}

func part2Calculation(data ParsedData) {
	total := 0
	validSequences := 0

	for _, sequence := range data.Sequences {
		if testSequence(&sequence, &data) {
			continue
		}

		slices.SortFunc(sequence.Data, func(a, b int) int {
			// a < b is true when b has a rule that contains a
			rule, hasRule := data.NumberDependancies[b]
			if hasRule {
				if contains(rule, a) {
					return -1
				}
			}

			// Only reorder when a rule is matched
			return 0
		})

		if testSequence(&sequence, &data) {
			total += sequence.Data[len(sequence.Data)/2]
			validSequences++
		} else {
			log.Fatal("The sequence should be valid after sorting", "sequence", sequence)
		}
	}

	log.Info("Complete!", "output", total, "validvalidSequences", validSequences)
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
