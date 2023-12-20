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
		debug bool
	}{
		{
			"1",
			args{strings.NewReader(`px{a<2006:qkq,m>2090:A,rfg}
			pv{a>1716:R,A}
			lnx{m>1548:A,A}
			rfg{s<537:gd,x>2440:R,A}
			qs{s>3448:A,lnx}
			qkq{x<1416:A,crn}
			crn{x>2662:A,R}
			in{s<1351:px,qqz}
			qqz{s>2770:qs,m<1801:hdj,R}
			gd{a>3333:R,R}
			hdj{m>838:A,pv}
			
			{x=787,m=2655,a=1222,s=2876}
			{x=1679,m=44,a=2067,s=496}
			{x=2036,m=264,a=79,s=2244}
			{x=2461,m=1339,a=466,s=291}
			{x=2127,m=1623,a=2188,s=1013}`)},
			`19114`,
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
