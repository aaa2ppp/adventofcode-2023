package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

type Dir int

const (
	LeftToRight Dir = 1 << iota
	RigthToLeft
	TopToBottom
	BottomToTop
)

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	plane, err := readPlane(br)
	if err != nil {
		return err
	}

	count := solution(plane)
	fmt.Fprintln(bw, count)

	return nil
}

func solution(plane [][]byte) int {
	n := len(plane)
	m := len(plane[0])
	visited := makeMatrix[Dir](n, m)

	var dfs func(i, j int, dir Dir)

	doStep := func(i, j int, dir Dir) {
		switch dir {
		case LeftToRight:
			j++
		case RigthToLeft:
			j--
		case TopToBottom:
			i++
		case BottomToTop:
			i--
		}

		if !(0 <= i && i < n && 0 <= j && j < m) {
			return
		}

		if visited[i][j]&dir != 0 {
			return
		}

		dfs(i, j, dir)
	}

	dfs = func(i, j int, dir Dir) {
		visited[i][j] |= dir

		switch plane[i][j] {
		case '/':
			switch dir {
			case LeftToRight:
				dir = BottomToTop
			case RigthToLeft:
				dir = TopToBottom
			case TopToBottom:
				dir = RigthToLeft
			case BottomToTop:
				dir = LeftToRight
			}
		case '\\':
			switch dir {
			case LeftToRight:
				dir = TopToBottom
			case RigthToLeft:
				dir = BottomToTop
			case TopToBottom:
				dir = LeftToRight
			case BottomToTop:
				dir = RigthToLeft
			}
		case '-':
			switch dir {
			case TopToBottom, BottomToTop:
				doStep(i, j, LeftToRight)
				doStep(i, j, RigthToLeft)
				return
			}
		case '|':
			switch dir {
			case LeftToRight, RigthToLeft:
				doStep(i, j, TopToBottom)
				doStep(i, j, BottomToTop)
				return
			}
		}

		doStep(i, j, dir)
	}

	calc := func() int {
		count := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if visited[i][j] != 0 {
					count++
				}
			}
		}
		return count
	}

	// TODO: Cейчас мы пробуем войти из каждой точки периметра.
	// Можно исключить точки входа в которых были отраржения в предыдущих обходах.

	count := 0
	for i := 0; i < n; i++ {
		clearMatrix(visited)
		dfs(i, 0, LeftToRight)
		count = max(count, calc())

		clearMatrix(visited)
		dfs(i, m-1, RigthToLeft)
		count = max(count, calc())
	}

	for j := 0; j < m; j++ {
		clearMatrix(visited)
		dfs(0, j, TopToBottom)
		count = max(count, calc())

		clearMatrix(visited)
		dfs(n-1, j, BottomToTop)
		count = max(count, calc())
	}

	return count
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func readPlane(br io.Reader) ([][]byte, error) {
	buf, err := io.ReadAll(br)
	if err != nil {
		return nil, err
	}

	buf = bytes.TrimSpace(buf)
	plane := bytes.Split(buf, []byte("\n"))

	for i := range plane {
		plane[i] = bytes.TrimSpace(plane[i])
	}

	return plane, nil
}

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}

func clearMatrix[T any](matrix [][]T) {
	for _, row := range matrix {
		for j := range row {
			row[j] = *(new(T)) // zero
		}
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
