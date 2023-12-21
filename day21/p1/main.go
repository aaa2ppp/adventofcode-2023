package main

import (
	"adventofcode-2023/lib/queue"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

var stepCount = 64

type Point struct {
	i, j int
}

func (p Point) Add(of Point) Point {
	p.i += of.i
	p.j += of.j
	return p
}

func readPlane(r io.Reader) ([][]byte, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	buf = bytes.TrimSpace(buf)
	plane := bytes.Split(buf, []byte("\n"))

	for i, row := range plane {
		plane[i] = bytes.TrimSpace(row)
	}

	return plane, nil
}

func getStartPoint(plane [][]byte) Point {
	for i, row := range plane {
		for j, c := range row {
			if c == 'S' {
				return Point{i, j}
			}
		}
	}
	return Point{-1, -1}
}

func doSteps(plane [][]byte, start Point, count int) int {
	n := len(plane)
	m := len(plane[0])

	valid := func(p Point) bool {
		return 0 <= p.i && p.i < n && 0 <= p.j && p.j < m
	}

	type item struct {
		Point
		step int
	}
	var frontier queue.Queue[item]

	plane[start.i][start.j] = '0'
	frontier.Push(item{start, 1})
	cellCount := [2]int{1, 0}

	for frontier.Size() > 0 {
		it := frontier.Pop()
		if it.step > count {
			break
		}

		for _, o := range [...]Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			p := it.Point.Add(o)

			if valid(p) && (plane[p.i][p.j] == '.' || plane[p.i][p.j] == 'S') {
				v := it.step % 2
				cellCount[v]++
				plane[p.i][p.j] = byte(v) + '0'
				frontier.Push(item{p, it.step + 1})
			}
		}
	}

	return cellCount[count%2]
}

func countChar(plane [][]byte, c byte) int {
	count := 0
	for _, row := range plane {
		for _, cell := range row {
			if cell == c {
				count++
			}
		}
	}
	return count
}

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	plane, err := readPlane(br)
	if err != nil {
		return err
	}

	// if debugEnable {
	// 	log.Print("readPlane:")
	// 	for _, row := range plane {
	// 		log.Printf("%s", row)
	// 	}
	// }

	start := getStartPoint(plane)
	if debugEnable {
		log.Printf("start: %v", start)
	}

	cellCount := doSteps(plane, start, stepCount)
	if debugEnable {
		log.Printf("doSteps %v %d:", start, stepCount)
		for _, row := range plane {
			log.Printf("%s", row)
		}
	}

	// cellCount := countChar(plane, byte(stepCount%2)+'0')

	fmt.Fprintln(bw, cellCount)
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
