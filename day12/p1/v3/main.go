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

func calcTempl(templ string, groups []int) int {
	templ2 := makeTempl2(groups)
	count2 := 0
	for _, c := range []byte(templ2) {
		if c == '.' || c == '#' {
			count2++
		}
	}

	if debugEnable {
		log.Printf("templ2: %s", templ2)
	}

	buf := make([]byte, len(templ))
	count := 0
	permutate(0, 0, templ, templ2, count2, buf, func() {
		if debugEnable {
			log.Printf("found: %s", buf)
		}
		count++
	})

	return count
}

func makeTempl2(groups []int) string {
	n := 0
	for _, v := range groups {
		n += v
	}
	n += len(groups) * 2

	templ2 := make([]byte, n+1)
	templ2[0] = '+'
	for i, j := 0, 1; i < len(groups); i++ {
		for k := 0; k < groups[i]; k++ {
			templ2[j] = '#'
			j++
		}
		templ2[j] = '.'
		j++
		templ2[j] = '+'
		j++
	}
	templ2[n-1] = '+'
	templ2 = templ2[:n]

	return string(templ2)
}

func permutate(i, j int, templ, templ2 string, count2 int, buf []byte, bingo func()) {
	if debugEnable {
		log.Printf("buf: %s|%s", buf[:i], templ[i:])
	}

	if i == len(templ) {
		if count2 == 0 {
			bingo()
		}
		return
	}

	if j == len(templ2) {
		return
	}

	if count2 > len(templ) - i {
		return
	}

	c1 := templ[i]
	c2 := templ2[j]

	switch c2 {
	case '#', '.':
		if count2 > 0 && (c1 == c2 || c1 == '?') {
			buf[i] = c2
			permutate(i+1, j+1, templ, templ2, count2-1, buf, bingo)
		}
	case '+':
		if c1 == '.' || c1 == '?' {
			buf[i] = c2
			permutate(i+1, j, templ, templ2, count2, buf, bingo)
		}
		permutate(i, j+1, templ, templ2, count2, buf, bingo)
	}
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
