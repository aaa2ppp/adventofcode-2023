package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	plane, err := readPlane(br)
	if err != nil {
		return err
	}

	n := len(plane)
	m := len(plane[0])

	if debugEnable {
		for i := range plane {
			log.Printf("%c", plane[i])
		}
	}

	vals := make([]int, m)
	for j := range vals {
		vals[j] = n
	}

	total := 0

	for i := range plane {
		for j := 0; j < m; j++ {
			switch plane[i][j] {
			case 'O':
				total += vals[j]
				vals[j]--
			case '#':
				vals[j] = n - i - 1
			}
		}
	}

	fmt.Fprintln(bw, total)
	return nil
}

func readPlane(br *bufio.Reader) ([][]byte, error) {
	buf, err := io.ReadAll(br)
	if err != nil {
		return nil, err
	}

	if buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}

	plane := bytes.Split(buf, []byte("\n"))
	for i := range plane {
		plane[i] = bytes.TrimSpace(plane[i])
	}

	return plane, nil
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
