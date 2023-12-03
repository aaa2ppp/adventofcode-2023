package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"unicode"
	"unsafe"
)

type Color int

const (
	ColorRed Color = iota
	ColorGreen
	ColorBulue
)

func check(data [][3]int) bool {
	for _, set := range data {
		if set[ColorRed] > 12 {
			return false
		}
		if set[ColorGreen] > 13 {
			return false
		}
		if set[ColorBulue] > 14 {
			return false
		}
	}
	return true
}

func parseLine(s []byte) ([][3]int, error) {
	scan := NewScanner(bytes.NewReader(s))

	// skip: Game #:
	scan.Scan()
	scan.Scan()

	var data [][3]int

	mainLoop:
	for {
		var set [3]int

		for {
			n, err := scan.Int()
			if err != nil {
				if err == io.EOF {
					break mainLoop
				}
				return nil, err
			}

			if !scan.Scan() {
				return nil, errors.New("can't read color")
			}

			color := unsafeString(scan.Bytes()) // <color>[,]
			delim := color[len(color)-1]

			if !unicode.IsLetter(rune(delim)) {
				color = color[:len(color)-1]
			} else {
				delim = 0
			}

			switch color {
			case "red":
				set[ColorRed] = n
			case "green":
				set[ColorGreen] = n
			case "blue":
				set[ColorBulue] = n
			default:
				return nil, fmt.Errorf("unknown color '%s'", color)
			}

			if delim != ',' {
				break
			}
		}

		data = append(data, set)
	}

	return data, nil
}

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	var sum int
	var lineNo int

	for {
		lineNo++
		s, isPrefix, err := br.ReadLine()
		if isPrefix {
			return fmt.Errorf("%d: line too long", lineNo)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%d: %w", lineNo, err)
		}

		data, err := parseLine(s)
		if err != nil {
			return fmt.Errorf("%d: can't parse: %w", lineNo, err)
		}
		if check(data) {
			sum += lineNo
		}
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
