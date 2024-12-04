package main

import (
	"errors"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/jessevdk/go-flags"
)

type ParsedData struct {
	letters [][]byte
}

const NOT_SET byte = 255

func parseData(data []byte) (ParsedData, error) {
	ret := ParsedData{letters: make([][]byte, 0)}
	currentLine := make([]byte, 0)

	for _, b := range data {
		if b == '\n' {
			ret.letters = append(ret.letters, currentLine)
			currentLine = make([]byte, 0)

			if len(ret.letters) > 1 {
				if len(ret.letters[len(ret.letters)-1]) != len(ret.letters[len(ret.letters)-2]) {
					return ParsedData{}, errors.New("Inconsistent length of rows")
				}
			}
		} else {
			currentLine = append(currentLine, b)
		}
	}

	return ret, nil
}

var PATTERN []byte = []byte{'X', 'M', 'A', 'S'}

func (p *ParsedData) scanLtR() uint {
	var sum uint

	for _, row := range p.letters {
		matchIndex := 0
		for _, b := range row {
			if b == PATTERN[0] {
				matchIndex = 0
			}

			if b == PATTERN[matchIndex] {
				matchIndex++

				if matchIndex >= len(PATTERN) {
					matchIndex = 0
					sum++
				}
			} else {
				matchIndex = 0
			}
		}
	}

	return sum
}

func (p *ParsedData) scanRtL() uint {
	var sum uint

	for _, row := range p.letters {
		matchIndex := 0
		for i := len(row) - 1; i >= 0; i-- {
			b := row[i]

			if b == PATTERN[0] {
				matchIndex = 0
			}

			if b == PATTERN[matchIndex] {
				matchIndex++

				if matchIndex >= len(PATTERN) {
					sum++
					matchIndex = 0
				}
			} else {
				matchIndex = 0
			}
		}
	}

	return sum
}

func (p *ParsedData) scanTtB() uint {
	var sum uint

	for row := 0; row < len(p.letters[0]); row++ {
		matchIndex := 0
		for col := 0; col < len(p.letters); col++ {
			b := p.letters[col][row]

			if b == PATTERN[0] {
				matchIndex = 0
			}

			if b == PATTERN[matchIndex] {
				matchIndex++

				if matchIndex >= len(PATTERN) {
					sum++
					matchIndex = 0
				}
			} else {
				matchIndex = 0
			}

		}
	}

	return sum
}

func (p *ParsedData) scanBtoT() uint {
	var sum uint

	for row := 0; row < len(p.letters[0]); row++ {
		matchIndex := 0
		for col := len(p.letters) - 1; col >= 0; col-- {
			b := p.letters[col][row]

			if b == PATTERN[0] {
				matchIndex = 0
			}

			if b == PATTERN[matchIndex] {
				matchIndex++
				if matchIndex >= len(PATTERN) {
					sum++
					matchIndex = 0
				}
			} else {
				matchIndex = 0
			}
		}
	}

	return sum
}

func (p *ParsedData) scanPointTLtBR(col, row int) bool {
	// Check bounds for an illegal scan
	if col+len(PATTERN)-1 >= len(p.letters) {
		return false
	}

	if row+len(PATTERN)-1 >= len(p.letters[0]) {
		return false
	}

	for i := 0; i < len(PATTERN); i++ {
		if p.letters[col][row] != PATTERN[i] {
			return false
		}

		col++
		row++
	}

	return true
}

func (p *ParsedData) scanTLtBR() uint {
	var sum uint

	for col := 0; col < len(p.letters); col++ {
		for row := 0; row < len(p.letters[0]); row++ {
			if p.scanPointTLtBR(row, col) {
				sum++
			}
		}
	}

	return sum
}

func (p *ParsedData) scanPointBRtTL(col, row int) bool {
	// Check bounds for an illegal scan
	if col-len(PATTERN)+1 < 0 {
		return false
	}

	if row-len(PATTERN)+1 < 0 {
		return false
	}

	for i := 0; i < len(PATTERN); i++ {
		if p.letters[col][row] != PATTERN[i] {
			return false
		}

		col--
		row--
	}

	return true
}

func (p *ParsedData) scanBRtTL() uint {
	var sum uint

	for col := 0; col < len(p.letters); col++ {
		for row := 0; row < len(p.letters[0]); row++ {
			if p.scanPointBRtTL(row, col) {
				sum++
			}
		}
	}

	return sum
}

func (p *ParsedData) scanPointTRtBL(col, row int) bool {
	// Check bounds for an illegal scan
	if col+len(PATTERN)-1 >= len(p.letters) {
		return false
	}

	if row-len(PATTERN)+1 < 0 {
		return false
	}

	for i := 0; i < len(PATTERN); i++ {
		if p.letters[col][row] != PATTERN[i] {
			return false
		}

		col++
		row--
	}

	return true
}

func (p *ParsedData) scanTRtBL() uint {
	var sum uint

	for col := 0; col < len(p.letters); col++ {
		for row := 0; row < len(p.letters[0]); row++ {
			if p.scanPointTRtBL(row, col) {
				sum++
			}
		}
	}

	return sum
}

func (p *ParsedData) scanPointBLtTR(col, row int) bool {
	// Check bounds for an illegal scan
	if col-len(PATTERN)+1 < 0 {
		return false
	}

	if row+len(PATTERN)-1 >= len(p.letters[0]) {
		return false
	}

	for i := 0; i < len(PATTERN); i++ {
		if p.letters[col][row] != PATTERN[i] {
			return false
		}

		col--
		row++
	}

	return true
}

func (p *ParsedData) scanBLtTR() uint {
	var sum uint

	for col := 0; col < len(p.letters); col++ {
		for row := 0; row < len(p.letters[0]); row++ {
			if p.scanPointBLtTR(row, col) {
				sum++
			}
		}
	}

	return sum
}

func part1Calculation(data ParsedData) {
	sum := 0
	var lock sync.Mutex
	var wg sync.WaitGroup

	addToSum := func(f func() uint, name string) {
		wg.Add(1)

		go func() {
			defer wg.Done()

			localSum := f()
			lock.Lock()
			defer lock.Unlock()

			sum += int(localSum)
			log.Info("Sub-task complete", "sub-task", name, "sum", localSum)
		}()
	}

	addToSum(func() uint { return data.scanLtR() }, "LtR")
	addToSum(func() uint { return data.scanRtL() }, "RtL")
	addToSum(func() uint { return data.scanTtB() }, "TtB")
	addToSum(func() uint { return data.scanBtoT() }, "BtT")
	addToSum(func() uint { return data.scanTLtBR() }, "TLtoBR")
	addToSum(func() uint { return data.scanBRtTL() }, "BRtTL")
	addToSum(func() uint { return data.scanTRtBL() }, "TRtBL")
	addToSum(func() uint { return data.scanBLtTR() }, "BLtTR")

	wg.Wait()

	log.Info("Complete", "output", sum)
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

	log.Info("Dimensions", "cols", len(data.letters), "rows", len(data.letters[0]))

	log.Info("Calculating output")
	if opts.Part2 {
		// part2Calculation(data)
	} else {
		part1Calculation(data)
	}
}
