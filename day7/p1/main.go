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

type HandType int

//go:generate stringer -type=HandType
const (
	_           HandType = 7 - iota
	FiveOfKind           // where all five cards have the same label: AAAAA
	FourOfKind           // where four cards have the same label and one card has a different label: AA8AA
	FullHouse            // where three cards have the same label, and the remaining two cards share a different label: 23332
	ThreeOfKind          // where three cards have the same label, and the remaining two cards are each different from any other card in the hand: TTT98
	TwoPair              // where two cards share one label, two other cards share a second label, and the remaining card has a third label: 23432
	OnePair              // where two cards share one label, and the other three cards have a different label from the pair and each other: A23A4
	HighCard             // where all cards' labels are distinct: 23456
)

type Hand struct {
	Raw  string
	Val  int
	Type HandType
}

func (h Hand) String() string {
	return fmt.Sprintf("{%s %05x %v}", h.Raw, h.Val, h.Type)
}

func NewHand(b string) *Hand {
	v := 0
	for i := 0; i < 5; i++ {
		shift := (4-i)*4
		// A, K, Q, J, T, 9, 8, 7, 6, 5, 4, 3, or 2.
		switch c := b[i]; {
		case c == 'A':
			v += 14 << shift
		case c == 'K':
			v += 13 << shift
		case c == 'Q':
			v += 12 << shift
		case c == 'J':
			v += 11 << shift
		case c == 'T':
			v += 10 << shift
		case '2' <= c && c <= '9':
			v += int(c-'0') << shift
		}
	}
	return &Hand{
		Raw:  b,
		Val:  v,
		Type: calcHandType(v),
	}
}

func calcHandType(v int) HandType {
	f1 := make([]int, 15)

	for i := 0; i < 5; i++ {
		f1[(v>>(i*4))&15]++
	}

	f2 := make([]int, 6)
	for _, v := range f1 {
		f2[v]++
	}

	switch {
	case f2[5] == 1:
		return FiveOfKind
	case f2[4] == 1:
		return FourOfKind
	case f2[3] == 1 && f2[2] == 1:
		return FullHouse
	case f2[3] == 1:
		return ThreeOfKind
	case f2[2] == 2:
		return TwoPair
	case f2[2] == 1:
		return OnePair
	default:
		return HighCard
	}
}

func (h *Hand) Less(other *Hand) bool {
	return h.Type < other.Type || h.Type == other.Type && h.Val < other.Val
}

type HandBid struct {
	Hand *Hand
	Bid  int
}

// func calcTotalWinning(hands []HandBid) int {

// }

func _run(scan *Scanner, bw *bufio.Writer) error {
	var hands []HandBid

	for scan.Scan() {
		h := NewHand(scan.Text())
		b, err := scan.Int()
		if err != nil {
			return err
		}
		hands = append(hands, HandBid{Hand: h, Bid: b})
	}

	if err := scan.Err(); err != nil {
		return err
	}

	sort.Slice(hands, func(i, j int) bool {
		return hands[i].Hand.Less(hands[j].Hand)
	})

	if debugEnable {
		log.Println("hands:", hands)
	}

	total := 0
	for i, h := range hands {
		total += (i+1)*h.Bid
	}

	fmt.Fprintln(bw, total)
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
