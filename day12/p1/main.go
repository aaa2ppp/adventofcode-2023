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

func permutate(blanks []int, i int, n int, bingo func()) {

	// TODO: это полный переребор, можно сделать отсечки

	if i == len(blanks) - 1 {
		blanks[i] += n
		bingo()
		blanks[i] -= n
		return
	}

	blanks[i] += n + 1
	for j := n; j >= 0; j-- {
		blanks[i]--
		permutate(blanks, i+1, n-j, bingo)
	}
}

func checkChar(templ string, c byte, n int) bool {
	for _, v := range []byte(templ[:n]) {
		if v != c && v != '?' {
			return false
		}
	}
	return true
}

func check(templ string, groups []int, blanks []int) bool {
	var gi, bi int

	for gi < len(groups) {
		if n := blanks[bi]; !checkChar(templ, '.', n) {
			return false
		} else {
			templ = templ[n:]
			bi++
		}

		if n := groups[gi]; !checkChar(templ, '#', n) {
			return false
		} else {
			templ = templ[n:]
			gi++
		}
	}

	if n := blanks[bi]; !checkChar(templ, '.', n) {
		return false
	}

	return true
}

func calcTempl(templ string, groups []int) int {
	n := len(templ)
	for _, v := range groups {
		n -= v
	}

	blanks := make([]int, len(groups)+1)
	for i := range blanks {
		blanks[i] = 1
		n--
	}

	blanks[0] = 0
	blanks[len(blanks)-1] = 0
	n+=2

	if debugEnable {
		log.Printf("%s %v %v %d", templ, groups, blanks, n)
	}

	count := 0
	permutate(blanks, 0, n, func() {
		ok := check(templ, groups, blanks)
		if ok {
			count++
		}
		if debugEnable {
			log.Printf("%s %v %v - %v", templ, groups, blanks, ok)
		}
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
