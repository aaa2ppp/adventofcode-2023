package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"unicode"
	"unsafe"
)

type Number struct {
	i      int
	j0, j1 int
	val    int
}

type SimbolSet map[[2]int]struct{}

func parseLine(i int, s []byte, nums *[]Number, simbols SimbolSet) error {
	var numVal int
	var numLen int

	for j, c := range s {
		if unicode.IsDigit(rune(c)) {
			numVal *= 10
			numVal += int(c - '0')
			numLen++
			continue
		}

		if numLen != 0 {
			*nums = append(*nums, Number{i: i, j0: j - numLen, j1: j, val: numVal})
			numVal = 0
			numLen = 0
		}

		if c != '.' {
			simbols[[2]int{i, j}] = struct{}{}
		}
	}

	if numLen != 0 {
		*nums = append(*nums, Number{i: i, j0: len(s) - numLen, j1: len(s), val: numVal})
		numVal = 0
		numLen = 0
	}

	return nil
}

func checkNumber(num *Number, simbols SimbolSet) int {
	for i, j := num.i-1, num.j0-1; j <= num.j1; j++ {
		if _, ok := simbols[[2]int{i, j}]; ok {
			return num.val
		}
	}

	for i, j := num.i+1, num.j0-1; j <= num.j1; j++ {
		if _, ok := simbols[[2]int{i, j}]; ok {
			return num.val
		}
	}

	if _, ok := simbols[[2]int{num.i, num.j0 - 1}]; ok {
		return num.val
	}

	if _, ok := simbols[[2]int{num.i, num.j1}]; ok {
		return num.val
	}

	return 0
}

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	var lineN int
	nums := []Number{}
	simbols := SimbolSet{}

	for {
		lineN++
		s, isPrefix, err := br.ReadLine()
		if isPrefix {
			return fmt.Errorf("%d: line too long", lineN)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%d: %w", lineN, err)
		}

		if err := parseLine(lineN, s, &nums, simbols); err != nil {
			return fmt.Errorf("%d: can't parse line: %w", lineN, err)
		}
	}

	var sum int
	for i := range nums {
		sum += checkNumber(&nums[i], simbols)
	}

	fmt.Fprintln(bw, sum)

	return nil
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
