package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

const maxCount = 1000

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	text, plane, err := readPlane(br)
	if err != nil {
		return err
	}

	var values []int
	hashes := make(map[uint64]int)

	var (
		periodBegin int
		periodSize  int
	)

	for i := 0; i < maxCount; i++ {
		toNorth(plane)
		toWest(plane)
		toSouth(plane)
		toEast(plane)

		h := djb2Hash(text)

		if idx, ok := hashes[h]; ok {
			periodBegin = idx
			periodSize = i - idx
			break
		}

		v := calcPlane(plane)

		values = append(values, v)
		hashes[h] = i
	}

	idx := (1_000_000_000 - periodBegin - 1) % periodSize + periodBegin

	if debugEnable {
		log.Printf("%d %d %v", periodBegin, periodSize, values)
	}

	fmt.Fprintln(bw, values[idx])
	return nil
}

func calcPlane(plane [][]byte) int {
	total := 0
	n := len(plane)

	for i := range plane {
		for _, c := range plane[i] {
			if c == 'O' {
				total += n - i
			}
		}
	}

	return total
}

func djb2Hash(text []byte) uint64 {
	hash := uint64(5381)
	for _, c := range text {
		hash = (hash << 5) + hash + uint64(c) // (hash * 33) + c
	}
	return hash
}

func debugPlane(title string, plane [][]byte) {
	if debugEnable {
		log.Printf("%s:", title)
		for i := range plane {
			log.Printf("%c", plane[i])
		}
	}
}

func testOfCicles(r io.Reader, w io.Writer, n int) {
	_, plane, _ := readPlane(r)

	for i := 0; i < n; i++ {
		toNorth(plane)
		toWest(plane)
		toSouth(plane)
		toEast(plane)

		if i == 0 {
			fmt.Fprintf(w, "After 1 cycle:\n")
		} else {
			fmt.Fprintf(w, "\nAfter %d cycles:\n", i+1)
		}

		for _, row := range plane {
			fmt.Fprintf(w, "%s\n", row)
		}
	}
}

func toNorth(plane [][]byte) int {
	n := len(plane)
	m := len(plane[0])

	vals := make([]int, m)
	for j := range vals {
		vals[j] = n
	}

	total := 0

	for i := range plane {
		for j := 0; j < m; j++ {
			switch plane[i][j] {
			case 'O':
				v := vals[j]
				ii := n - v
				plane[i][j], plane[ii][j] = plane[ii][j], plane[i][j]
				total += v
				vals[j]--
			case '#':
				vals[j] = n - i - 1
			}
		}
	}

	return total
}

func toWest(plane [][]byte) int {
	n := len(plane)
	m := len(plane[0])

	vals := make([]int, n)
	for i := range vals {
		vals[i] = m
	}

	total := 0

	for i := range plane {
		for j := 0; j < m; j++ {
			switch plane[i][j] {
			case 'O':
				v := vals[i]
				jj := m - v
				plane[i][j], plane[i][jj] = plane[i][jj], plane[i][j]
				total += v
				vals[i]--
			case '#':
				vals[i] = m - j - 1
			}
		}
	}

	return total
}

func toSouth(plane [][]byte) int {
	n := len(plane)
	m := len(plane[0])

	vals := make([]int, m)
	for j := range vals {
		vals[j] = n
	}

	total := 0

	for i := n - 1; i >= 0; i-- {
		for j := 0; j < m; j++ {
			switch plane[i][j] {
			case 'O':
				v := vals[j]
				ii := v - 1
				plane[i][j], plane[ii][j] = plane[ii][j], plane[i][j]
				total += v
				vals[j]--
			case '#':
				vals[j] = i
			}
		}
	}

	return total
}

func toEast(plane [][]byte) int {
	n := len(plane)
	m := len(plane[0])

	vals := make([]int, n)
	for i := range vals {
		vals[i] = m
	}

	total := 0

	for i := range plane {
		for j := m - 1; j >= 0; j-- {
			switch plane[i][j] {
			case 'O':
				v := vals[i]
				jj := v - 1
				plane[i][j], plane[i][jj] = plane[i][jj], plane[i][j]
				total += v
				vals[i]--
			case '#':
				vals[i] = j
			}
		}
	}

	return total
}

func reverse[T any](a []T) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func transposeCopy[T any](matrix [][]T) [][]T {
	n := len(matrix)
	m := len(matrix[0])

	matrix2 := makeMatrix[T](m, n)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			matrix2[j][i] = matrix[i][j]
		}
	}

	return matrix2
}

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}

func readPlane(r io.Reader) ([]byte, [][]byte, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}

	if buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}

	plane := bytes.Split(buf, []byte("\n"))
	for i := range plane {
		plane[i] = bytes.TrimSpace(plane[i])
	}

	return buf, plane, nil
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
