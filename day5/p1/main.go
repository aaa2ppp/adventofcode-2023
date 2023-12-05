package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

type MapItem struct {
	dst int
	src int
	len int
}

type Map []MapItem

func (m Map) Sort() {
	sort.Slice(m, func(i, j int) bool {
		return m[i].src < m[j].src
	})
}

func (m Map) Get(src int) int {
	i := sort.Search(len(m), func(i int) bool {
		return m[i].src > src
	})
	i--

	if i != -1 && src < m[i].src+m[i].len {
		return m[i].dst + (src - m[i].src)
	}

	return src
}

func ScanMap(scan *Scanner) (Map, error) {
	var err error

	// skip: `map:`
	scan.Scan()

	var m Map
	for {
		var it MapItem
		it.dst, it.src, it.len, err = scan.ThreeInt()
		if err != nil {
			break
		}
		m = append(m, it)
	}

	m.Sort()
	return m, err
}

func _run(scan *Scanner, bw *bufio.Writer) error {
	var err error

	// skip: `seeds:`
	scan.Scan()

	var seeds []int
	for {
		var v int
		v, err = scan.Int()
		if err != nil {
			break
		}
		seeds = append(seeds, v)
	}

	maps := make([]Map, 7)

	for err != io.EOF {
		switch scan.Text() {
		case "seed-to-soil":
			maps[0], err = ScanMap(scan)
		case "soil-to-fertilizer":
			maps[1], err = ScanMap(scan)
		case "fertilizer-to-water":
			maps[2], err = ScanMap(scan)
		case "water-to-light":
			maps[3], err = ScanMap(scan)
		case "light-to-temperature":
			maps[4], err = ScanMap(scan)
		case "temperature-to-humidity":
			maps[5], err = ScanMap(scan)
		case "humidity-to-location":
			maps[6], err = ScanMap(scan)
		default:
			return err
		}
	}

	if debugEnable {
		log.Println(maps)
	}

	minimum := math.MaxInt
	res := make([]int, 0, 8)

	for _, v := range seeds {
		res = res[:0]
		res = append(res, v)
		for _, m := range maps {
			v = m.Get(v)
			res = append(res, v)
		}
		if debugEnable {
			log.Println(res)
		}
		minimum = min(minimum, v)
	}

	fmt.Fprintln(bw, minimum)
	return nil
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
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
