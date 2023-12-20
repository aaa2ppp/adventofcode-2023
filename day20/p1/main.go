package main

import (
	"adventofcode-2023/lib/queue"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	pressButtonCount = 1000
	broadcasterName  = "broadcaster"
)

type PulseType int

const (
	_ PulseType = iota
	LowPulse
	HightPulse
)

func (pt PulseType) String() string {
	switch pt {
	case LowPulse:
		return "-low"
	case HightPulse:
		return "-hight"
	default:
		return fmt.Sprintf("PulseType(%d)", pt)
	}
}

type Pulse struct {
	From string
	To   string
	Val  PulseType
}

func (p Pulse) String() string {
	return fmt.Sprintf("%s %s-> %s", p.From, p.Val, p.To)
}

type Sender interface {
	Send(p Pulse)
}

type Module interface {
	Name() string
	Outputs() []string
	RegInput(name string)
	Process(p Pulse, sender Sender)
}

type Stub struct{}

func (s Stub) Name() string          { return "--STUB--" }
func (s Stub) Outputs() []string     { return nil }
func (s Stub) RegInput(string)       {}
func (s Stub) Process(Pulse, Sender) {}

type Schema struct {
	modules    map[string]Module
	queue      queue.Queue[Pulse]
	LowCount   int
	HightCount int
}

func NewSchema() *Schema {
	return &Schema{
		modules: map[string]Module{},
	}
}

func (s *Schema) AddModule(mod Module) {
	s.modules[mod.Name()] = mod
}

func (s *Schema) Prepare() {
	const op = "Schema.Prepare"

	for senderName, sender := range s.modules {
		for _, receiverName := range sender.Outputs() {
			receiver := s.modules[receiverName]
			if receiver == nil {
				if debugEnable {
					log.Printf("%s: add stub '%s'", op, receiverName)
				}
				s.modules[receiverName] = Stub{}
			} else {
				receiver.RegInput(senderName)
			}
		}
	}
}

func (s *Schema) Send(p Pulse) {
	s.queue.Push(p)

	if debugEnable && pressButtonCount < 10 {
		log.Print(p)
	}

	switch p.Val {
	case LowPulse:
		s.LowCount++
	case HightPulse:
		s.HightCount++
	}
}

func (s *Schema) PressButton() {
	s.Send(Pulse{From: "button", To: broadcasterName, Val: LowPulse})
	for s.queue.Size() > 0 {
		p := s.queue.Pop()
		m := s.modules[p.To]
		m.Process(p, s)
	}
}

type FlipFlop struct {
	name    string
	on      bool
	outputs []string
}

var _ Module = (*FlipFlop)(nil)

func NewFlipFlop(name string, outputs []string) *FlipFlop {
	return &FlipFlop{
		name:    name,
		outputs: outputs,
	}
}

func (ff *FlipFlop) Name() string {
	return ff.name
}

func (ff *FlipFlop) Outputs() []string {
	return ff.outputs
}

func (ff *FlipFlop) RegInput(name string) { /* nothong */ }

func (ff *FlipFlop) Process(p Pulse, s Sender) {
	if p.Val == HightPulse {
		return
	}

	ff.on = !ff.on
	val := LowPulse
	if ff.on {
		val = HightPulse
	}

	for _, output := range ff.outputs {
		s.Send(Pulse{From: ff.name, To: output, Val: val})
	}
}

type Conjunction struct {
	name    string
	inputs  map[string]PulseType
	state   int
	outputs []string
}

var _ Module = (*Conjunction)(nil)

func NewConjuction(name string, outputs []string) *Conjunction {
	return &Conjunction{
		name:    name,
		inputs:  map[string]PulseType{},
		outputs: outputs,
	}
}

func (c *Conjunction) Name() string {
	return c.name
}

func (c *Conjunction) Outputs() []string {
	return c.outputs
}

func (c *Conjunction) RegInput(name string) {
	c.inputs[name] = LowPulse
}

func (c *Conjunction) Process(p Pulse, s Sender) {
	if c.inputs[p.From] != p.Val {
		if p.Val == LowPulse {
			c.state--
		} else {
			c.state++
		}
		c.inputs[p.From] = p.Val
	}

	val := HightPulse
	if c.state == len(c.inputs) {
		val = LowPulse
	}

	for _, output := range c.outputs {
		s.Send(Pulse{From: c.name, To: output, Val: val})
	}
}

type Broadcaster struct {
	outputs []string
}

var _ Module = (*Broadcaster)(nil)

func NewBroadcaster(outputs []string) *Broadcaster {
	return &Broadcaster{
		outputs: outputs,
	}
}

func (b *Broadcaster) Name() string {
	return broadcasterName
}

func (b *Broadcaster) Outputs() []string {
	return b.outputs
}

func (b *Broadcaster) RegInput(name string) { /* nothing */ }

func (b *Broadcaster) Process(p Pulse, s Sender) {
	for _, output := range b.outputs {
		s.Send(Pulse{From: broadcasterName, To: output, Val: p.Val})
	}
}

func parseModule(s string) (Module, error) {
	var (
		name    string
		arrow   string
		outputs []string
	)

	r := strings.NewReader(s)
	_, err := fmt.Fscan(r, &name, &arrow)
	if err != nil {
		return nil, err
	}

	for {
		var output string
		_, err = fmt.Fscan(r, &output)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if n := len(output) - 1; output[n] == ',' {
			output = output[:n]
		}

		outputs = append(outputs, output)
	}

	if len(outputs) == 0 {
		return nil, fmt.Errorf("no any outputs")
	}

	switch name[0] {
	case '%':
		if len(name) < 2 {
			return nil, fmt.Errorf("no module name")
		}
		return NewFlipFlop(name[1:], outputs), nil
	case '&':
		if len(name) < 2 {
			return nil, fmt.Errorf("no module name")
		}
		return NewConjuction(name[1:], outputs), nil
	default: // broadcast
		if name != broadcasterName {
			return nil, fmt.Errorf("unknown module type: %s", name)
		}
		return NewBroadcaster(outputs), nil
	}
}

func _run(sc *bufio.Scanner, bw *bufio.Writer) error {
	schema := NewSchema()

	i := 0
	for sc.Scan() {
		i++
		mod, err := parseModule(sc.Text())
		if err != nil {
			return fmt.Errorf("line %d: %w", i, err)
		}
		schema.AddModule(mod)
	}

	schema.Prepare()

	for i := 0; i < pressButtonCount; i++ {
		schema.PressButton()
	}

	fmt.Fprintln(bw, schema.LowCount*schema.HightCount)

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
