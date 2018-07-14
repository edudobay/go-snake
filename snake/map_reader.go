package snake

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkN(n, actualN int, e error) {
	check(e)
	if n != actualN {
		panic(fmt.Errorf("%d != expected %d", actualN, n))
	}
}

func assert(condition bool, msg string) {
	if !condition {
		panic(fmt.Errorf("assertion failed: %s", msg))
	}
}

func ReadMap(filename string) GameMap {

	f, err := os.Open(filename)
	check(err)

	var width, height int
	n, err := fmt.Fscanf(f, "%d %d", &width, &height)
	check(err)
	assert(n == 2, "invalid size header")

	size := width * height
	assert(size > 0, "invalid width/height")

	cells := make([]MapCellType, size)

	reader := bufio.NewReader(f)

	for i := 0; i < size; {
		b, _, err := reader.ReadRune()
		if err == io.EOF {
			panic("premature end of file")
		}

		if unicode.IsSpace(b) {
			continue
		}

		switch b {
		case '.':
			cells[i] = MapCellFree
		case '#':
			cells[i] = MapCellWall
		case 'x':
			cells[i] = MapCellInvalid
		default:
			panic("invalid char in map")
		}

		i++
	}

	return GameMap{width, height, cells}
}
