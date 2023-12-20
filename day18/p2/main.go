package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Dir byte

const (
	U Dir = iota
	D
	L
	R
)

type PlaneItem struct {
	dir Dir
	len int
}

type Point struct {
	i, j int
}

func getNextPoint(p1, p2, p3 Point) Point {
	const op = "getNextPoint"

	v1 := p2.i - p1.i
	h1 := p2.j - p1.j

	v2 := p3.i - p2.i
	h2 := p3.j - p2.j

	if h1 > 0 {
		if v2 > 0 {
			return Point{p2.i, p2.j + 1}
		} else if v2 < 0 {
			return Point{p2.i, p2.j}
		}
	} else if h1 < 0 {
		if v2 > 0 {
			return Point{p2.i + 1, p2.j + 1}
		} else if v2 < 0 {
			return Point{p2.i + 1, p2.j}
		}
	} else if v1 > 0 {
		if h2 > 0 {
			return Point{p2.i, p2.j + 1}
		} else if h2 < 0 {
			return Point{p2.i + 1, p2.j + 1}
		}
	} else if v1 < 0 {
		if h2 > 0 {
			return Point{p2.i, p2.j}
		} else if h2 < 0 {
			return Point{p2.i + 1, p2.j}
		}
	}

	panic(fmt.Sprintf("%s: bad sequention: %v %v %v (%d, %d, %d, %d)", op, p1, p2, p3, v1, h1, v2, h2))
}

func getPath(plane []PlaneItem) []Point {
	path := make([]Point, 0, len(plane))

	var i, j int

	for _, v := range plane {
		switch v.dir {
		case U:
			i -= v.len
		case D:
			i += v.len
		case L:
			j -= v.len
		case R:
			j += v.len
		}
		path = append(path, Point{i, j})
	}

	return path
}

func getPoints(path []Point) []Point {
	points := make([]Point, 0, len(path))

	p1 := path[len(path)-1]
	p2 := path[0]
	for i := 1; i < len(path); i++ {
		p3 := path[i]
		points = append(points, getNextPoint(p1, p2, p3))
		p1 = p2
		p2 = p3
	}
	points = append(points, getNextPoint(p1, p2, path[0]))

	return points
}

func reverse[T any](a []T) []T {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func calcArea(points []Point) int {
	total := 0

	p1 := points[len(points)-1]
	for i := 0; i < len(points); i++ {
		p2 := points[i]
		if p1.j == p2.j {
			total += (p2.i - p1.i) * p1.j
		}
		p1 = p2
	}

	return abs(total)
}

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	plane, err := readPlane(br)
	if err != nil {
		return err
	}

	path := getPath(plane)

	if debugEnable {
		log.Println("path:", path)
	}

	points1 := getPoints(path)

	if debugEnable {
		log.Println("points:", points1)
	}

	area1 := calcArea(points1)

	reverse(path)
	points2 := getPoints(path)

	if debugEnable {
		log.Println("points:", points2)
	}

	area2 := calcArea(points2)

	fmt.Fprintln(bw, max(area1, area2))

	return nil
}

func readPlane(r io.Reader) ([]PlaneItem, error) {
	var plane []PlaneItem

	for {
		var (
			it  PlaneItem
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

		// // TODO: parseDir()

		// switch dir {
		// case "U":
		// 	it.dir = U
		// case "D":
		// 	it.dir = D
		// case "L":
		// 	it.dir = L
		// case "R":
		// 	it.dir = R
		// default:
		// 	return nil, fmt.Errorf("unknown dir '%s'", dir)
		// }

		// TODO: parseColor()

		// 012345678
		// (#XXXXXD)
		switch clr[7] {
		case '0':
			it.dir = R
		case '1':
			it.dir = D
		case '2':
			it.dir = L
		case '3':
			it.dir = U
		}

		len64, err := strconv.ParseInt(clr[2:7], 16, 64)
		if err != nil {
			return nil, err
		}

		it.len = int(len64)
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
