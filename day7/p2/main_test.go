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
			args{strings.NewReader(`32T3K 765
			T55J5 684
			KK677 28
			KTJJT 220
			QQQJA 483`)},
			`5905`,
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