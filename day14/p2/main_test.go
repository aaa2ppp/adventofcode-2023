package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
		debug   bool
	}{
		{
			"1",
			args{strings.NewReader(`O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....`)},
			`64`,
			false,
			true,
		},
		// {
		// 	"2",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			defer func() { debugEnable = false }()
			w := &bytes.Buffer{}
			if err := run(tt.args.r, w); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); strings.TrimSpace(gotW) != strings.TrimSpace(tt.wantW) {
				t.Errorf("run() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_testOfCicles(t *testing.T) {
	type args struct {
		r io.Reader
		n int
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{
			"1",
			args{strings.NewReader(`O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....`),
			3,

			},
			`After 1 cycle:
.....#....
....#...O#
...OO##...
.OO#......
.....OOO#.
.O#...O#.#
....O#....
......OOOO
#...O###..
#..OO#....

After 2 cycles:
.....#....
....#...O#
.....##...
..O#......
.....OOO#.
.O#...O#.#
....O#...O
.......OOO
#..OO###..
#.OOO#...O

After 3 cycles:
.....#....
....#...O#
.....##...
..O#......
.....OOO#.
.O#...O#.#
....O#...O
.......OOO
#...O###.O
#.OOO#...O`,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			testOfCicles(tt.args.r, w, tt.args.n)
			if gotW := w.String(); strings.TrimSpace(gotW) != strings.TrimSpace(tt.wantW) {
				t.Errorf("testOfCicles() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
