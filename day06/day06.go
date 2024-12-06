package main

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/jessevdk/go-flags"
)

type MapTile byte

const (
	MapTileEmpty     MapTile = '.'
	MapTileWall      MapTile = '#'
	StartingLocation         = '^'
)

type Direction int

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

type Vec2 struct {
	X, Y int
}

type ParsedData struct {
	Grid      [][]MapTile
	Location  Vec2
	Direction Direction
}

func (p *ParsedData) clone() *ParsedData {
	ret := &ParsedData{Location: p.Location, Direction: p.Direction, Grid: make([][]MapTile, 0)}

	for _, line := range p.Grid {
		lineCopy := make([]MapTile, len(line))
		copy(lineCopy, line)

		ret.Grid = append(ret.Grid, lineCopy)
	}

	return ret
}

func parseData(data []byte) (ParsedData, error) {
	ret := ParsedData{Grid: make([][]MapTile, 0), Direction: DirectionUp}

	line := make([]MapTile, 0)
	x := 0
	y := 0

	for _, b := range data {
		if b == '\n' {
			y++
			x = 0

			if len(ret.Grid) > 1 {
				if len(line) != len(ret.Grid[0]) {
					return ret, errors.New("There grid is not a rectangle")
				}
			}

			ret.Grid = append(ret.Grid, line)
			line = make([]MapTile, 0)
		} else {
			switch b {
			case StartingLocation:
				ret.Location.X = x
				ret.Location.Y = y
				log.Info("Guard starting at", "x", x, "y", y)

				b = byte(MapTileEmpty)
				fallthrough
			case byte(MapTileWall):
				fallthrough
			case byte(MapTileEmpty):
				line = append(line, MapTile(b))
			default:
				return ret, fmt.Errorf("Unrecognised char %c", b)
			}

			x++
		}
	}

	return ret, nil
}

func (p *ParsedData) rotate() {
	p.Direction++
	p.Direction %= 4
}

func (p *ParsedData) canMoveTo(X, Y int) bool {
	return p.Grid[Y][X] == MapTileEmpty
}

func (p *ParsedData) tryMoveGuard() (bool, error) {
	switch p.Direction {
	case DirectionUp:
		if p.Location.Y-1 < 0 {
			return true, nil
		}

		if p.canMoveTo(p.Location.X, p.Location.Y-1) {
			p.Location.Y--
		} else {
			return false, errors.New("There is a wall in the way of the guard")
		}
	case DirectionDown:
		if p.Location.Y+1 >= len(p.Grid) {
			return true, nil
		}

		if p.canMoveTo(p.Location.X, p.Location.Y+1) {
			p.Location.Y++
		} else {
			return false, errors.New("There is a wall in the way of the guard")
		}
	case DirectionLeft:
		if p.Location.X-1 < 0 {
			return true, nil
		}

		if p.canMoveTo(p.Location.X-1, p.Location.Y) {
			p.Location.X--
		} else {
			return false, errors.New("There is a wall in the way of the guard")
		}
	case DirectionRight:
		if p.Location.X+1 >= len(p.Grid[0]) {
			return true, nil
		}

		if p.canMoveTo(p.Location.X+1, p.Location.Y) {
			p.Location.X++
		} else {
			return false, errors.New("There is a wall in the way of the guard")
		}
	}

	return false, nil
}

func (p *ParsedData) moveGuard() bool {
	for i := 0; true; i++ {
		exitedMap, err := p.tryMoveGuard()
		if exitedMap {
			return true
		}

		if err != nil {
			p.rotate()

			if i == 3 {
				log.Warn("The guard is spinning (what a loser)")
				return false
			}
		} else {
			break
		}
	}

	return false
}

type StateSeen struct {
	Location  Vec2
	Direction Direction
}

func (p *ParsedData) guardPositions() (int, error) {
	states := make(map[StateSeen]bool)
	positions := make(map[Vec2]bool)

  for {
		state := StateSeen{Location: p.Location, Direction: p.Direction}

		_, found := states[state]
		if found {
			log.Warn("The guard is in a loop - he will never leave the lab.",
				"location", p.Location,
				"direction", p.Direction,
				"positions", len(positions),
				"states", len(states))
			return 0, errors.New("The guard is lost")
		}

		states[state] = true
		positions[p.Location] = true

		exited := p.moveGuard()
		if exited {
			log.Debug("The guard has left the map")
			break
		}
	}

	return len(positions), nil
}

func part1Calculation(data ParsedData) {
	count, err := data.guardPositions()
	if err != nil {
		log.Fatal("Guard is unable to escape", "err", err)
	}

	log.Info("Complete", "output", count)
	return
}

func part2Calculation(data ParsedData) {
	positions := make(map[Vec2]bool)

  p := data.clone()
  for {
		positions[p.Location] = true

		exited := p.moveGuard()
		if exited {
			break
		}
	}

  var wg sync.WaitGroup
  var lock sync.Mutex
  total := 0

  log.Info("Inserting obstacles", "obstacles", len(positions))

  for position := range positions {
    wg.Add(1)
    go func() {
      defer wg.Done()

      newGrid := data.clone()
      newGrid.Grid[position.Y][position.X] = MapTileWall

      _, err := newGrid.guardPositions()
      if err != nil {
        lock.Lock()
        defer lock.Unlock()
        total++
      }
    }()
  }

  wg.Wait()
  log.Info("Completed brute force", "output", total)
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
