package main

import (
	"adventofcode-2023/lib/queue"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type Dir byte

const (
	U Dir = iota
	D
	L
	R
)

type DigPlanItem struct {
	dir   Dir
	len   int
	color int
}

func minMaxIJ(plan []DigPlanItem) (minI, minJ, maxI, maxJ int) {
	var i, j int

	for _, v := range plan {
		switch v.dir {
		case U:
			i -= v.len
			minI = min(minI, i)
		case D:
			i += v.len
			maxI = max(maxI, i)
		case L:
			j -= v.len
			minJ = min(minJ, j)
		case R:
			j += v.len
			maxJ = max(maxJ, j)
		}
	}

	return
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func draw(desk [][]byte, plane []DigPlanItem, i, j int, c byte) {
	desk[i][j] = c
	for _, v := range plane {
		for k := 0; k < v.len; k++ {
			switch v.dir {
			case U:
				i--
			case D:
				i++
			case L:
				j--
			case R:
				j++
			}
			desk[i][j] = c
		}
	}
}

type Point struct {
	i, j int
}

func (p Point) Add(of Point) Point {
	p.i += of.i
	p.j += of.j
	return p
}

func fill(desk [][]byte, i, j int, c byte) {
	n := len(desk)
	m := len(desk[0])

	var frontier queue.Queue[Point]

	desk[i][j] = c
	frontier.Push(Point{i, j})

	for frontier.Size() > 0 {
		p := frontier.Pop()

		for _, of := range []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			p2 := p.Add(of)
			if i, j := p2.i, p2.j; 0 <= i && i < n && 0 <= j && j < m && desk[i][j] == ' ' {
				desk[i][j] = c
				frontier.Push(p2)
			}
		}

	}
}

func count(desk [][]byte, c byte) int {
	n := len(desk)
	m := len(desk[0])

	count := 0

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if desk[i][j] == c {
				count++
			}
		}
	}

	return count
}

const blank = ' '

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	plane, err := readPlane(br)
	if err != nil {
		return err
	}

	i0, j0, i1, j1 := minMaxIJ(plane)

	n := (i1 - i0 + 1) + 2
	m := (j1 - j0 + 1) + 2

	desk := MakeMatrix[byte](n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			desk[i][j] = blank
		}
	}

	draw(desk, plane, -i0+1, -j0+1, '#')
	debugDesk("draw:", desk)

	fill(desk, 0, 0, '.')
	debugDesk("fill:", desk)

	cnt := count(desk, '.')
	if debugEnable {
		log.Println("cnt:", cnt)
	}

	fmt.Fprintln(bw, n*m-cnt)
	return nil
}

func debugDesk(title string, desk [][]byte) {
	if debugEnable {
		log.Println(title)
		for _, row := range desk {
			log.Printf("%c", row)
		}
	}
}

func readPlane(r io.Reader) ([]DigPlanItem, error) {
	var plane []DigPlanItem

	for {
		var (
			it  DigPlanItem
			dir string
			clr string
		)

		_, err := fmt.Fscan(r, &dir, &it.len, &clr)

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// TODO: parseDir()

		switch dir {
		case "U":
			it.dir = U
		case "D":
			it.dir = D
		case "L":
			it.dir = L
		case "R":
			it.dir = R
		default:
			return nil, fmt.Errorf("unknown dir '%s'", dir)
		}

		// TODO: parseColor()

		plane = append(plane, it)
	}

	return plane, nil
}

func run(r io.Reader, w io.Writer) (err error) {
	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)
	defer func() {
		if flushErr := bw.Flush(); flushErr != nil && err == nil {
			err = flushErr
		}
	}()

	return _run(br, bw)
}

func main() {
	_ = debugEnable
	if err := run(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

var _, debugEnable = os.LookupEnv("DEBUG")

func MakeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}
