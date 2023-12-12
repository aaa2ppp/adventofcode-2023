package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	total := 0

	for i := 0; ; i++ {
		line, isPrefix, err := br.ReadLine()
		if isPrefix { // XXX
			return fmt.Errorf("%d: line too long", i+1)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%d: %w", i+1, err)
		}

		templ, groups, err := parseLine(line)
		if err != nil {
			return fmt.Errorf("%d: %w", i+1, err)
		}

		total += calcTempl(templ, groups)
	}

	fmt.Fprintln(bw, total)
	return nil
}

func parseLine(line []byte) (string, []int, error) {
	// `<templ>[ [<num>,...]]`

	p := bytes.IndexByte(line, ' ')
	if p == -1 {
		return string(line), nil, nil
	}
	templ := string(line[:p])

	groups, err := parseGroups(bytes.TrimSpace(line[p:]))
	if err != nil {
		return templ, groups, fmt.Errorf("can't parse groups: %w", err)
	}

	return templ, groups, nil
}

func parseGroups(line []byte) ([]int, error) {
	if len(line) == 0 {
		return nil, nil
	}

	items := bytes.Split(line, []byte(","))
	groups := make([]int, 0, len(items))

	for _, it := range items {
		v, err := strconv.Atoi(string(it)) // TODO: unsafeString
		if err != nil {
			return groups, err
		}
		groups = append(groups, v)
	}

	return groups, nil
}

func permutate(templ string, groups []int, i int, n int, bingo func()) {
	if i == len(groups) {
		if _, ok := checkChar(templ, '.', n); ok {
			bingo()
		}
		return
	}

	n2 := n - (len(groups) - i - 1) // максимальное кол-во пробелов в этой позиции
	k2 := 1
	if i == 0 {
		k2 = 0
		n2++ // XXX чтобы проверить k2==0
	}
	m := 0


	for n2 > 0 {
		if debugEnable {
			log.Printf("%d %s: %d %d %d", i, templ, n2, k2, m)
		}

		if _, ok := checkChar(templ, '.', k2); !ok {
			return
		}

		templ = templ[k2:]
		n -= k2
		n2 -= k2

		if k, ok := checkChar(templ[m:], '#', groups[i]-m); ok {
			permutate(templ[groups[i]:], groups, i+1, n, bingo)
			m = k - 1
			k2 = 1
		} else if k < n2 {
			m = 0
			k2 = k + 1
		} else {
			return
		}

	}
}

func checkChar(templ string, c byte, n int) (int, bool) {
	for i, v := range []byte(templ[:n]) {
		if v != c && v != '?' {
			return i, false
		}
	}
	return n, true
}

func calcTempl(templ string, groups []int) int {
	n := len(templ)
	for _, v := range groups {
		n -= v
	}

	if debugEnable {
		log.Printf("%s %v %d", templ, groups, n)
	}

	count := 0
	permutate(templ, groups, 0, n, func() {
		count++
	})

	return count
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
