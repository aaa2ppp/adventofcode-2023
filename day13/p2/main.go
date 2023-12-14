package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/bits"
	"os"
	"strconv"
	"unsafe"
)

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	sum := 0

	for {
		rows, cols, err := readMirror(br)
		if rows == nil {
			if err != nil && err != io.EOF {
				return err
			}
			break
		}

		if debugEnable {
			log.Println("rows:", rows)
			log.Println("cols:", cols)
		}

		if n, ok := searchAxis(cols); ok {
			sum += n
			continue
		}

		if n, ok := searchAxis(rows); ok {
			sum += 100 * n
		}
	}

	fmt.Fprintln(bw, sum)
	return nil
}

func searchAxis(nums []uint) (int, bool) {
	for n := 1; n < len(nums); n++ {
		bingo := 0

		for i, j := n-1, n; i >= 0 && j < len(nums); i, j = i-1, j+1 {
			bingo += bits.OnesCount(nums[i] ^ nums[j])
			if bingo > 1 {
				break
			}
		}

		if bingo == 1 {
			return n, true
		}
	}

	return 0, false
}

func readMirror(br *bufio.Reader) (rows, cols []uint, err error) {
	for i := 0; ; i++ {
		line, isPrefix, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				return rows, cols, err
			}
			return nil, nil, err
		}
		if isPrefix {
			return nil, nil, errors.New("line too long")
		}
		if len(line) == 0 {
			return rows, cols, nil
		}

		if cols == nil {
			if len(line) > 64 {
				return nil, nil, errors.New("line is longer than 64 characters")
			}
			cols = make([]uint, len(line))

		}

		rows = append(rows, calcLine(line, i, cols))
	}
}

func calcLine(line []byte, i int, cols []uint) uint {
	var row uint
	for j, c := range line {
		if c == '#' {
			cols[j] |= 1 << i
			row |= 1 << j
		}
	}
	return row
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
