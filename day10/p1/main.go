package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"unsafe"
)

func getStart(plane [][]byte) Point {
	for i, line := range plane {
		for j, c := range line {
			if c == 'S' {
				return Point{i, j}
			}
		}
	}
	return Point{-1, -1}
}

type Point struct {
	i, j int
}

func (p Point) Add(dir Point) Point {
	p.i += dir.i
	p.j += dir.j
	return p
}

func (p Point) Valid(n, m int) bool {
	return 0 <= p.i && p.i < n && 0 <= p.j && p.j < m
}

var (
	North = Point{-1, 0}
	South = Point{1, 0}
	East  = Point{0, 1}
	West  = Point{0, -1}
)

func getToDir(c byte, fromDir Point) (Point, bool) {
	switch c {
	case '|': // представляет собой вертикальную трубу, соединяющую север и юг.
		switch fromDir {
		case North:
			return South, true
		case South:
			return North, true
		}
	case '-': // представляет собой горизонтальную трубу, соединяющую восток и запад.
		switch fromDir {
		case East:
			return West, true
		case West:
			return East, true
		}
	case 'L': // представляет собой 90-градусный изгиб, соединяющий север и восток.
		switch fromDir {
		case North:
			return East, true
		case East:
			return North, true
		}
	case 'J': // представляет собой 90-градусный изгиб, соединяющий север и запад.
		switch fromDir {
		case North:
			return West, true
		case West:
			return North, true
		}
	case '7': // представляет собой 90-градусный изгиб, соединяющий юг и запад.
		switch fromDir {
		case West:
			return South, true
		case South:
			return West, true
		}
	case 'F': // представляет собой 90-градусный изгиб, соединяющий юг и восток.
		switch fromDir {
		case East:
			return South, true
		case South:
			return East, true
		}
	case '.': // это земли; нет трубы в плитке.
	case 'S': // это исходное положение животного; на этой плитке изображена труба, но на вашем эскизе не показано, какую форму имеет труба.
	}

	return Point{}, false
}

type Way struct {
	point   Point
	fromDir Point
}

func doStep(plane [][]byte, way Way) (Way, bool) {
	toDir, ok := getToDir(plane[way.point.i][way.point.j], way.fromDir)
	if ok {
		way.point = way.point.Add(toDir)
		way.fromDir = Point{-1 * toDir.i, -1 * toDir.j}
	}
	return way, false
}

func searchForBeast(plane [][]byte) int {
	n := len(plane)
	m := len(plane[0])

	start := getStart(plane)
	if debugEnable {
		log.Println("start:", start)
	}

	way := [2]Way{}

	k := 0
	for _, toDir := range []Point{North, South, West, East} {
		if k == 2 {
			break
		}

		fromDir := Point{-1 * toDir.i, -1 * toDir.j}
		if p := start.Add(toDir); p.Valid(n, m) {
			c := plane[p.i][p.j]
			if _, ok := getToDir(c, fromDir); ok {
				way[k] = Way{fromDir: fromDir, point: p}
				k++
			}
		}
	}

	if k != 2 {
		panic("not found ways from start")
	}

	count := 1
	for way[0].point != way[1].point && count < n*m/2 {
		if debugEnable {
			log.Printf("%d: %v %v", count, way[0], way[1])
		}
		way[0], _ = doStep(plane, way[0])
		way[1], _ = doStep(plane, way[1])
		count++
	}

	return count
}

func makeMatrix(n, m int) [][]int {
	buf := make([]int, n*m)
	matrix := make([][]int, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	plane, err := readPlane(br)
	if err != nil {
		return err
	}

	count := searchForBeast(plane)

	fmt.Fprintln(bw, count)
	return nil
}

func readPlane(r io.Reader) ([][]byte, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	bytes.TrimSpace(buf)
	plane := bytes.Split(buf, []byte("\n"))
	for i := range plane {
		plane[i] = bytes.TrimSpace(plane[i])
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

// Scanner wrapper for the bufio.Scanner with split by words
type Scanner struct {
	bufio.Scanner
}

func NewScanner(r io.Reader) *Scanner {
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanWords)
	return (*Scanner)(unsafe.Pointer(sc))
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//go:noinline
func (sc *Scanner) restoreEOF() error {
	if sc.Err() != nil {
		return sc.Err()
	}
	return io.EOF
}

func (sc *Scanner) Int() (int, error) {
	if sc.Scan() {
		return strconv.Atoi(unsafeString(sc.Bytes()))
	}
	return 0, sc.restoreEOF()
}

func (sc *Scanner) TwoInt() (n1, n2 int, err error) {
	n1, err = sc.Int()
	if err == nil {
		n2, err = sc.Int()
	}
	return
}

func (sc *Scanner) ThreeInt() (n1, n2, n3 int, err error) {
	n1, err = sc.Int()
	if err == nil {
		n2, n3, err = sc.TwoInt()
	}
	return
}

func (sc *Scanner) FourInt() (n1, n2, n3, n4 int, err error) {
	n1, n2, err = sc.TwoInt()
	if err == nil {
		n3, n4, err = sc.TwoInt()
	}
	return
}

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

func ScanToIntSlice[T Int](sc *Scanner, slice []T) (int, error) {
	bitSize := int(unsafe.Sizeof(*(new(T))) * 8)

	for i := 0; i < len(slice); i++ {
		if !sc.Scan() {
			return i, sc.restoreEOF()
		}

		if bitSize <= math.MaxInt {
			v, err := strconv.Atoi(unsafeString(sc.Bytes()))
			if err != nil {
				return i, err
			}
			slice[i] = T(v)
		} else {
			v, err := strconv.ParseInt(unsafeString(sc.Bytes()), 10, bitSize)
			if err != nil {
				return i, err
			}
			slice[i] = T(v)
		}
	}

	return len(slice), nil
}

func WriteIntSlice[T Int](w *bufio.Writer, slice []T, delim string) (int, error) {
	if len(slice) == 0 {
		return 0, nil
	}

	buf := make([]byte, 0, 32) // TODO: how to make it not escape to heap?

	buf = strconv.AppendInt(buf, int64(slice[0]), 10)
	if _, err := w.Write(buf); err != nil {
		return 0, err
	}

	for i := 1; i < len(slice); i++ {
		buf = buf[:0]
		buf = append(buf, delim...)
		buf = strconv.AppendInt(buf, int64(slice[i]), 10)
		if _, err := w.Write(buf); err != nil {
			return i, err
		}
	}

	return len(slice), nil
}
