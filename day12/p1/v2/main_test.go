package main

import (
	"bytes"
	"io"
	"reflect"
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
			args{strings.NewReader(`???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1`)},
			`21`,
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

func Test_parseLine(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []int
		wantErr bool
	}{
		{
			"???.### 1,1,3",
			args{[]byte("???.### 1,1,3")},
			"???.###",
			[]int{1, 1, 3},
			false,
		},
		{
			"???.###",
			args{[]byte("???.###")},
			"???.###",
			nil,
			false,
		},
		{
			"???.### ",
			args{[]byte("???.### ")},
			"???.###",
			nil,
			false,
		},
		{
			"???.### 12",
			args{[]byte("???.### 12")},
			"???.###",
			[]int{12},
			false,
		},
		{
			"???.### 12,a",
			args{[]byte("???.### 12,a")},
			"???.###",
			[]int{12},
			true,
		},
		{
			"???.### aaa,bbb",
			args{[]byte("???.### aaa,bbb")},
			"???.###",
			[]int{},
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseLine() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseLine() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_calcTempl(t *testing.T) {
	type args struct {
		templ  string
		groups []int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		debug bool
	}{
		{
			"???.### 1,1,3",
			args{"???.###", []int{1, 1, 3}},
			1,
			true,
		},
		{
			".??..??...?##. 1,1,3",
			args{".??..??...?##.", []int{1, 1, 3}},
			4,
			true,
		},
		{
			"?#?#?#?#?#?#?#? 1,3,1,6",
			args{"?#?#?#?#?#?#?#?", []int{1, 3, 1, 6}},
			1,
			true,
		},
		{
			"????.#...#... 4,1,1",
			args{"????.#...#...", []int{4, 1, 1}},
			1,
			true,
		},
		{
			"????.######..#####. 1,6,5",
			args{"????.######..#####.", []int{1, 6, 5}},
			4,
			true,
		},
		{
			"?###???????? 3,2,1",
			args{"?###????????", []int{3, 2, 1}},
			10,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			if got := calcTempl(tt.args.templ, tt.args.groups); got != tt.want {
				t.Errorf("calcTempl() = %v, want %v", got, tt.want)
			}
		})
	}
}
