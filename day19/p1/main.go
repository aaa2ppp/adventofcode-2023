package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// x: Extremely cool looking
// m: Musical (it makes a noise when you hit it)
// a: Aerodynamic
// s: Shiny

type Category int

const (
	X Category = iota
	M
	A
	S
)

type Part []int

type Condition func(v, x int) bool

type Rule struct {
	cat      Category
	val      int
	cond     Condition
	workflow string
}

func (r *Rule) Apply(part Part) bool {
	return r.cond(r.val, part[r.cat])
}

func Less(v, x int) bool {
	return x < v
}

func More(v, x int) bool {
	return x > v
}

func Always(v, x int) bool {
	return true
}

func parseCategory(c byte) (Category, error) {
	switch c {
	case 'x':
		return X, nil
	case 'm':
		return M, nil
	case 'a':
		return A, nil
	case 's':
		return S, nil
	}
	return -1, fmt.Errorf("unknown category '%c'", c)
}

func parseRule(s string) (*Rule, error) {
	p := strings.IndexByte(s, ':')
	if p == -1 {
		return &Rule{
			cond:     Always,
			workflow: s,
		}, nil
	}

	cat, err := parseCategory(s[0])
	if err != nil {
		return nil, err
	}

	var cond Condition
	switch s[1] {
	case '>':
		cond = More
	case '<':
		cond = Less
	default:
		return nil, fmt.Errorf("unknown condition '%c", s[1])
	}

	val, err := strconv.Atoi(s[2:p])
	if err != nil {
		return nil, err
	}

	return &Rule{
		cat:      cat,
		cond:     cond,
		val:      val,
		workflow: s[p+1:],
	}, nil
}

type Workflow struct {
	name  string
	rules []*Rule
}

func parseWorkflow(s string) (*Workflow, error) {
	p := strings.IndexByte(s, '{')
	name := s[:p]

	var rules []*Rule
	for _, it := range strings.Split(s[p+1:len(s)-1], ",") {
		r, err := parseRule(it)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}

	return &Workflow{
		name:  name,
		rules: rules,
	}, nil
}

func parsePart(s string) (Part, error) {
	part := make(Part, 4)

	for _, it := range strings.Split(s[1:len(s)-1], ",") {
		cat, err := parseCategory(it[0])
		if err != nil {
			return nil, err
		}
		val, err := strconv.Atoi(it[2:])
		if err != nil {
			return nil, err
		}
		part[cat] = val
	}

	return part, nil
}

func check(p Part, workflows map[string]*Workflow) bool {
	w := workflows["in"]

mainLoop:
	for i := 0; i < len(workflows); i++ {
		for _, r := range w.rules {
			if r.Apply(p) {
				switch r.workflow {
				case "A":
					return true
				case "R":
					return false
				default:
					if debugEnable {
						log.Printf("check: %s", r.workflow)
					}
					w = workflows[r.workflow]
					continue mainLoop
				}
			}
		}
	}

	panic("check: infinite loop detected")
}

func _run(sc *bufio.Scanner, bw *bufio.Writer) error {
	var (
		workflows = map[string]*Workflow{}
		parts     []Part
	)

	i := 0

	for sc.Scan() {
		i++
		s := strings.TrimSpace(sc.Text())
		if s == "" {
			break
		}
		w, err := parseWorkflow(s)
		if err != nil {
			return fmt.Errorf("%d: %w", i, err)
		}
		if debugEnable {
			log.Println("w:", w)
		}
		workflows[w.name] = w
	}

	if debugEnable {
		log.Println("workflows:", workflows)
	}

	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		p, err := parsePart(s)
		if err != nil {
			return fmt.Errorf("%d: %w", i, err)
		}
		parts = append(parts, p)
	}

	count := 0
	for _, p := range parts {
		if check(p, workflows) {
			count += sum(p...)
		}
	}

	fmt.Fprintln(bw, count)

	return nil
}

func run(r io.Reader, w io.Writer) (err error) {
	sc := bufio.NewScanner(r)
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

func sum(a ...int) int {
	total := 0
	for _, v := range a {
		total += v
	}
	return total
}
