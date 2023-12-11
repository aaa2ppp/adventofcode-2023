package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"unsafe"
)

var multiplier = int(1e6)

func _run(br *bufio.Reader, bw *bufio.Writer) error {

	// читаем карту (координаты объектов и размеры карты)

	var n, m int
	var points [][2]int

	for {
		line, isPrefix, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%d: %w", n+1, err)
		}
		if isPrefix { // XXX
			return fmt.Errorf("%d: line too long", n+1)
		}

		for j, c := range line {
			if c == '#' {
				points = append(points, [2]int{n, j})
			}
		}

		n++
		m = max(m, len(line))
	}

	// ищем расширившиеся строки и столбцы (в которых нет объектов)

	rows := make([]int, m+1)
	cols := make([]int, n+1)

	for i := 1; i < len(rows); i++ {
		rows[i] = multiplier-1
	}

	for i := 1; i < len(cols); i++ {
		cols[i] = multiplier-1
	}

	for _, p := range points {
		rows[p[0]] = 0
		cols[p[1]] = 0
	}

	// вычисляем увеличение для каждой строки и столбца (префикные суммы)

	for i := 1; i < len(rows); i++ {
		rows[i] += rows[i-1]
	}

	for i := 1; i < len(cols); i++ {
		cols[i] += cols[i-1]
	}

	// корректируем (увеличиваем) координаты объектов

	for i := range points {
		p := &points[i]
		p[0] += rows[p[0]]
		p[1] += cols[p[1]]
	}

	// считаем сумму манхэттенских расстояний между объектами

	total := 0
	for i, p1 := range points {
		for _, p2 := range points[i+1:] {
			dist := abs(p2[0]-p1[0]) + abs(p2[1]-p1[1])
			total += dist
		} 
	}

	// bingo!

	fmt.Fprintln(bw, total)
	return nil
}

func max(a, b int) int {
	if a < b {
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
