package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_doSteps(t *testing.T) {
	type args struct {
		plan  [][]byte
		v0    int
		start Point
		count int
	}
	tests := []struct {
		name  string
		args  args
		want  [2]int
		want1 [][]byte
	}{
		{
			"3x3 0 0,1 2",
			args{
				[][]byte{
					[]byte("..."),
					[]byte("..."),
					[]byte("..."),
				},
				0,
				Point{0, 1},
				2,
			},
			[2]int{4, 3},
			[][]byte{
				[]byte("101"),
				[]byte("010"),
				[]byte(".0."),
			},
		},
		{
			"3x3 1 0,1 2",
			args{
				[][]byte{
					[]byte("..."),
					[]byte("..."),
					[]byte("..."),
				},
				1,
				Point{0, 1},
				2,
			},
			[2]int{3, 4},
			[][]byte{
				[]byte("010"),
				[]byte("101"),
				[]byte(".1."),
			},
		},
		{
			"3x3 0 1,1 1",
			args{
				[][]byte{
					[]byte("..."),
					[]byte("..."),
					[]byte("..."),
				},
				0,
				Point{1, 1},
				1,
			},
			[2]int{1, 4},
			[][]byte{
				[]byte(".1."),
				[]byte("101"),
				[]byte(".1."),
			},
		},
		{
			"3x3 1 1,1 1",
			args{
				[][]byte{
					[]byte("..."),
					[]byte("..."),
					[]byte("..."),
				},
				1,
				Point{1, 1},
				1,
			},
			[2]int{4, 1},
			[][]byte{
				[]byte(".0."),
				[]byte("010"),
				[]byte(".0."),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := doSteps(tt.args.plan, tt.args.v0, tt.args.start, tt.args.count)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doSteps() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("doSteps() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_solution1(t *testing.T) {
	type args struct {
		plan [][]byte
		k    int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		debug bool
	}{
		{
			"3x3 2",
			args{
				[][]byte{
					[]byte("..."),
					[]byte("..."),
					[]byte("..."),
				},
				2,
			},
			64,
			true,
		},
		{
			"3x3 3",
			args{
				[][]byte{
					[]byte("..."),
					[]byte("..."),
					[]byte("..."),
				},
				3,
			},
			121,
			true,
		},
		{
			"5x5 3",
			args{
				[][]byte{
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
				},
				3,
			},
			324,
			true,
		},
		{
			"5x5 4",
			args{
				[][]byte{
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
				},
				4,
			},
			529,
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			if got := solution1(tt.args.plan, tt.args.k); got != tt.want {
				t.Errorf("solution1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func readPlanFile(inputFile string) [][]byte {
	r, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	plan, err := readPlan(r)
	if err != nil {
		panic(err)
	}

	return plan
}

func Test_solution2(t *testing.T) {
	bigPlan := readPlanFile("../adventofcode.com_2023_day_21_input.txt")

	type args struct {
		plan [][]byte
		k    int
	}
	tests := []struct {
		name string
		args args
		want int
		debug bool
	}{
		{
			"3x3 2",
			args{
				[][]byte{
					[]byte("..."),
					[]byte("..."),
					[]byte("..."),
				},
				2,
			},
			64,
			true,
		},
		{
			"3x3 3",
			args{
				[][]byte{
					[]byte("..."),
					[]byte("..."),
					[]byte("..."),
				},
				3,
			},
			121,
			true,
		},
		{
			"5x5 3",
			args{
				[][]byte{
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
				},
				3,
			},
			324,
			true,
		},
		{
			"5x5 4",
			args{
				[][]byte{
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
					[]byte("....."),
				},
				4,
			},
			529,
			true,
		},
		{
			"5x5 4 x",
			args{
				[][]byte{
					[]byte("....."),
					[]byte(".#..."),
					[]byte("....."),
					[]byte(".#..."),
					[]byte("....."),
				},
				4,
			},
			-1,
			true,
		},
		{
			"7x7 4 x",
			args{
				[][]byte{
					[]byte("......."),
					[]byte(".#...#."),
					[]byte(".##...."),
					[]byte("...S..."),
					[]byte(".#..##."),
					[]byte(".#....."),
					[]byte("......."),
				},
				4,
			},
			-1,
			true,
		},
		{
			"7x7 100 x",
			args{
				[][]byte{
					[]byte("......."),
					[]byte(".#...#."),
					[]byte(".##...."),
					[]byte("...S..."),
					[]byte(".#..##."),
					[]byte(".#....."),
					[]byte("......."),
				},
				4,
			},
			-1,
			true,
		},
		{
			"7x7 101 x",
			args{
				[][]byte{
					[]byte("......."),
					[]byte(".#...#."),
					[]byte(".##...."),
					[]byte("...S..."),
					[]byte(".#..##."),
					[]byte(".#....."),
					[]byte("......."),
				},
				4,
			},
			-1,
			true,
		},
		{
			"bigPlan 3",
			args{
				bigPlan,
				3,
			},
			-1,
			true,
		},
		{
			"bigPlan 4",
			args{
				bigPlan,
				4,
			},
			-1,
			true,
		},
		{
			"bigPlan 10",
			args{
				bigPlan,
				10,
			},
			-1,
			true,
		},
		{
			"bigPlan 11",
			args{
				bigPlan,
				11,
			},
			-1,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = false
			if tt.want == -1 {
				tt.want = solution1(tt.args.plan, tt.args.k)
			}
			debugEnable = tt.debug
			if got := solution2(tt.args.plan, tt.args.k); got != tt.want {
				t.Errorf("solution2() = %v, want %v", got, tt.want)
			}
		})
	}
}

