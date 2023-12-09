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

const modulo = 36

func encode(b []byte) uint16 {
	var v uint16
	for i := 0; i < 3; i++ {
		v *= modulo
		switch c := b[i]; {
		case '0' <= c && c <= '9':
			v += uint16(c - '0' + 10)
		case 'A' <= c && c <= 'Z':
			v += uint16(c - 'A')
		}
	}
	return v
}

func _run(scan *Scanner, bw *bufio.Writer) error {
	scan.Scan()
	instructions := make([]byte, len(scan.Bytes()))

	for i, c := range scan.Bytes() {
		if c == 'R' {
			instructions[i] = 1
		}
	}

	nodes := make([][2]uint16, modulo*modulo*modulo)

	var starts []uint16
	start := encode([]byte("AAA")) % modulo // ??A

	for scan.Scan() {
		node := encode(scan.Bytes())
		if node%modulo == start {
			starts = append(starts, node)
		}

		// skip: =
		scan.Scan()

		scan.Scan()
		left := encode(bytes.Trim(scan.Bytes(), "(),"))

		scan.Scan()
		right := encode(bytes.Trim(scan.Bytes(), "(),"))

		nodes[node] = [2]uint16{left, right}
	}

	if debugEnable {
		log.Println("starts:", starts)
	}

	finish := encode([]byte("ZZZ")) % modulo // ??Z
	count := 0
	i := 0

	isFinish := func() bool {
		for _, start := range starts {
			if start%modulo != finish {
				return false
			}
		}
		return true
	}

	for !isFinish() {
		step := instructions[i]
		if i++; i == len(instructions) {
			i = 0
		}

		for i, start := range starts {
			starts[i] = nodes[start][step]
		}

		count++
	}

	fmt.Fprintln(bw, count)
	return nil
}

func run(r io.Reader, w io.Writer) (err error) {
	sc := NewScanner(r)
	bw := bufio.NewWriter(w)
	defer func() {
		if flushErr := bw.Flush(); flushErr != nil && err == nil {
			err = flushErr
		}
	}()

	return _run(sc, bw)
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
