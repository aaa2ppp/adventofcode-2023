package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

func _run(br *bufio.Reader, bw *bufio.Writer) error {
	plane, err := readPlane(br)
	if err != nil {
		return err
	}

	if debugEnable {
		for _, row := range plane {
			log.Println(row)
		}
	}

	graph := makeGraph(plane)

	total := searchMinPath(
		graph,
		[]int{0, 1},
		[]int{len(graph) - 1, len(graph) - 2},
	)

	fmt.Fprintln(bw, total)
	return nil
}

func searchMinPath(graph [][]Edge, start, finish []int) int {
	nodes := make([]Item, len(graph))
	for i := range nodes {
		nodes[i] = Item{
			id:    i,
			dist:  math.MaxInt,
			index: -1,
		}
	}

	frontier := &PriorityQueue{}

	for _, id := range start {
		nodes[id].dist = 0
		heap.Push(frontier, &nodes[id])
	}

	for frontier.Len() > 0 {
		node := heap.Pop(frontier).(*Item)
		for _, id := range finish {
			if node.id == id {
				return node.dist
			}	
		}

		for _, edge := range graph[node.id] {
			neig := &nodes[edge.neigID]
			dist := node.dist + edge.dist
			if dist < neig.dist {
				neig.dist = dist
				if neig.index == -1 {
					heap.Push(frontier, neig)
				} else if dist < neig.dist {
					heap.Fix(frontier, neig.index)
				}
			}
		}
	}

	return -1
}

type Edge struct {
	neigID int
	dist   int
}

type Dir byte

const (
	V Dir = iota
	H
)

func getNodeID(m, i, j int, dir Dir) int {
	return (i*m+j)*2 + int(dir)
}

func getVNeigs(plane [][]byte, i, j int) []Edge {
	n := len(plane)
	m := len(plane[0])

	neigs := make([]Edge, 0, 6)

	{
		w := 0
		for k, i2 := 1, i+1; k <= 3 && i2 < n; k, i2 = k+1, i2+1 {
			w += int(plane[i2][j])
			idx := getNodeID(m, i2, j, V)
			neigs = append(neigs, Edge{neigID: idx, dist: w})
		}
	}

	{
		w := 0
		for k, i2 := 1, i-1; k <= 3 && i2 >= 0; k, i2 = k+1, i2-1 {
			w += int(plane[i2][j])
			idx := getNodeID(m, i2, j, V)
			neigs = append(neigs, Edge{neigID: idx, dist: w})
		}
	}

	return neigs
}

func getHNeigs(plane [][]byte, i, j int) []Edge {
	// n := len(plane)
	m := len(plane[0])

	neigs := make([]Edge, 0, 6)

	{
		w := 0
		for k, j2 := 1, j+1; k <= 3 && j2 < m; k, j2 = k+1, j2+1 {
			w += int(plane[i][j2])
			idx := getNodeID(m, i, j2, H)
			neigs = append(neigs, Edge{neigID: idx, dist: w})
		}
	}

	{
		w := 0
		for k, j2 := 1, j-1; k <= 3 && j2 >= 0; k, j2 = k+1, j2-1 {
			w += int(plane[i][j2])
			idx := getNodeID(m, i, j2, H)
			neigs = append(neigs, Edge{neigID: idx, dist: w})
		}
	}

	return neigs
}

func makeGraph(plane [][]byte) [][]Edge {
	n := len(plane)
	m := len(plane[0])

	graph := make([][]Edge, n*m*2)

	var idx int
	for i, row := range plane {
		for j := range row {
			idx = getNodeID(m, i, j, H)
			graph[idx] = getVNeigs(plane, i, j)

			idx = getNodeID(m, i, j, V)
			graph[idx] = getHNeigs(plane, i, j)
		}
	}

	return graph
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

	for _, row := range plane {
		for j := range row {
			row[j] -= '0'
		}
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

// An Item is something we manage in a priority queue.
type Item struct {
	id   int
	dist int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
