package main

import (
	"adventofcode-2023/lib/queue"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	stepCount      = 26501365
	_, debugEnable = os.LookupEnv("DEBUG")
)

type Point struct {
	i, j int
}

func (p Point) Add(of Point) Point {
	p.i += of.i
	p.j += of.j
	return p
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.i, p.j)
}

func (p *Point) Set(s string) error {
	const op = "Point.Set"

	v := strings.Split(s, ",")
	if len(v) != 2 {
		return fmt.Errorf("%s: two comma-separated numbers are required: %s", op, s)
	}

	i, err := strconv.Atoi(strings.TrimPrefix(v[0], "("))
	if err != nil {
		return fmt.Errorf("%s: bad first number: %w", op, err)
	}

	j, err := strconv.Atoi(strings.TrimSuffix(v[1], ")"))
	if err != nil {
		return fmt.Errorf("%s: bad second number %w", op, err)
	}

	p.i = i
	p.j = j

	return nil
}

func readPlan(r io.Reader) ([][]byte, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	buf = bytes.TrimSpace(buf)
	plan := bytes.Split(buf, []byte("\n"))

	for i, row := range plan {
		plan[i] = bytes.TrimSpace(row)
	}

	return plan, nil
}

func getStartPoint(plan [][]byte) Point {
	for i, row := range plan {
		for j, c := range row {
			if c == 'S' {
				return Point{i, j}
			}
		}
	}
	return Point{-1, -1}
}

func copyPlan(plan [][]byte) [][]byte {
	n := len(plan)
	m := len(plan[0])

	buf := make([]byte, n*m)
	plan2 := make([][]byte, n)

	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		row := buf[j : j+m]
		copy(row, plan[i])
		plan2[i] = row
	}

	return plan2
}

func doSteps(plan [][]byte, v0 int, start Point, count int) ([2]int, [][]byte) {
	v0 %= 2
	plan = copyPlan(plan)

	n := len(plan)
	m := len(plan[0])

	valid := func(p Point) bool {
		return 0 <= p.i && p.i < n && 0 <= p.j && p.j < m
	}

	type item struct {
		Point
		step int
	}
	var frontier queue.Queue[item]

	var cellCount [2]int

	cellCount[v0]++
	plan[start.i][start.j] = byte(v0) + '0'
	frontier.Push(item{start, 1})

	for frontier.Size() > 0 {
		it := frontier.Pop()
		if it.step > count {
			break
		}

		for _, o := range [...]Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			p := it.Point.Add(o)

			if valid(p) && (plan[p.i][p.j] == '.' || plan[p.i][p.j] == 'S') {
				v := (it.step + v0) % 2
				cellCount[v]++
				plan[p.i][p.j] = byte(v) + '0'
				frontier.Push(item{p, it.step + 1})
			}
		}
	}

	return cellCount, plan
}

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}

func expandPlan(plan [][]byte, k int) [][]byte {
	n := len(plan)
	m := len(plan[0])

	n2 := n * (k*2 + 1)
	m2 := m * (k*2 + 1)

	plan2 := makeMatrix[byte](n2, m2)

	for i := 0; i < n2; i++ {
		ii := i % n
		for j := 0; j < m2; j += m {
			copy(plan2[i][j:], plan[ii])
		}
	}

	return plan2
}

func solution1(plan [][]byte, k int) int {
	n := len(plan)
	m := len(plan[0])

	if n%2 != 1 || m%2 != 1 {
		panic(fmt.Sprintf("solution1: plan sizes must be odd, got %dx%d", n, m))
	}

	plan2 := expandPlan(plan, k)
	n2 := len(plan2)
	m2 := len(plan2[0])

	s := n2 / 2
	x, plan := doSteps(plan2, 0, Point{n2 / 2, m2 / 2}, s)
	if debugEnable {
		for _, row := range plan {
			log.Printf("%s", row)
		}
	}

	return x[s%2]
}

func checkPlan(plan [][]byte) error {
	n := len(plan)
	m := len(plan[0])

	if n != m {
		panic(fmt.Sprintf("plan must be square, got %dx%d", n, m))
	}

	if n%2 != 1 {
		panic(fmt.Sprintf("plan size must be odd, got %dx%d", n, m))
	}

	for i := 0; i < n; i++ {
		for _, p := range []Point{{i, 0}, {i, n / 2}, {i, n - 1}, {0, i}, {n / 2, i}, {n - 1, i}} {
			switch c := plan[p.i][p.j]; c {
			case '.', 'S':
				// ok
			default:
				return fmt.Errorf("the boundaries and axes of the plan must be clean, got %v = '%c'", p, c)
			}
		}
	}

	return nil
}

func solution2(plan [][]byte, k int) int {
	if err := checkPlan(plan); err != nil {
		panic(fmt.Errorf("solution2: %v", err))
	}

	if k < 1 {
		panic(fmt.Sprintf("solution2: k must be greater 1, got %d", k))
	}

	n := len(plan)
	n2 := n * (k*2 + 1)

	x, _ := doSteps(plan, 0, Point{n / 2, n / 2}, n-1)

	fv0 := (n/2 + (k-1)*n + 1) % 2
	fs := n - 1
	f1, _ := doSteps(plan, fv0, Point{n / 2, 0}, fs)
	f2, _ := doSteps(plan, fv0, Point{n / 2, n - 1}, fs)
	f3, _ := doSteps(plan, fv0, Point{0, n / 2}, fs)
	f4, _ := doSteps(plan, fv0, Point{n - 1, n / 2}, fs)

	if debugEnable {
		log.Printf("fv0=%d fs=%d", fv0, fs)
		log.Printf("f: %v %v %v %v", f1, f2, f3, f4)
	}

	c1v0 := k % 2
	c1s := n - 1 + n/2
	c11, _ := doSteps(plan, c1v0, Point{0, 0}, c1s)
	c12, _ := doSteps(plan, c1v0, Point{0, n - 1}, c1s)
	c13, _ := doSteps(plan, c1v0, Point{n - 1, 0}, c1s)
	c14, _ := doSteps(plan, c1v0, Point{n - 1, n - 1}, c1s)

	if debugEnable {
		log.Printf("c1v0=%d c1s=%d", c1v0, c1s)
		log.Printf("c1: %v %v %v %v", c11, c12, c13, c14)
	}

	c2v0 := (k + 1) % 2
	c2s := n / 2
	c21, _ := doSteps(plan, c2v0, Point{0, 0}, c2s)
	c22, _ := doSteps(plan, c2v0, Point{0, n - 1}, c2s)
	c23, _ := doSteps(plan, c2v0, Point{n - 1, 0}, c2s)
	c24, _ := doSteps(plan, c2v0, Point{n - 1, n - 1}, c2s)

	if debugEnable {
		log.Printf("c2v0=%d c2s=%d", c2v0, c2s)
		log.Printf("c2: %v %v %v %v", c21, c22, c23, c24)
	}

	s := n2 / 2
	p0 := s % 2
	p1 := (p0 + 1) % 2

	total := x[p0]

	for i := 1; i < k; i += 2 {
		total += x[p1] * i * 4
	}

	for i := 2; i < k; i += 2 {
		total += x[p0] * i * 4
	}

	total += f1[p0]
	total += f2[p0]
	total += f3[p0]
	total += f4[p0]

	total += c11[p0] * (k - 1)
	total += c12[p0] * (k - 1)
	total += c13[p0] * (k - 1)
	total += c14[p0] * (k - 1)

	total += c21[p0] * k
	total += c22[p0] * k
	total += c23[p0] * k
	total += c24[p0] * k

	return total
}

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	plan, err := readPlan(br)
	if err != nil {
		return err
	}

	if err := checkPlan(plan); err != nil {
		log.Fatal(fmt.Errorf("bad plan: %v", err))
	}

	n := len(plan)
	k := (stepCount - n/2) / n

	if k*n+n/2 != stepCount {
		log.Fatal(fmt.Errorf("bad stepCount=%d", stepCount))
	}

	cellCount := solution2(plan, k)
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

	flag.Parse()

	if err := run(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
