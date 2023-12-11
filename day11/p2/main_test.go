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
		name       string
		args       args
		multiplier int
		wantW      string
		wantErr    bool
		debug      bool
	}{
		{
			"2",
			args{strings.NewReader(`...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....`),
			},
			2,
			`374`,
			false,
			true,
		},
		{
			"10",
			args{strings.NewReader(`...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....`),
			},
			10,
			`1030`,
			false,
			true,
		},
		{
			"100",
			args{strings.NewReader(`...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....`),
			},
			100,
			`8410`,
			false,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			multiplier = tt.multiplier
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
