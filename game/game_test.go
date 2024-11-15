package game

import "testing"

func Test_loopPoints(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "default", args: args{n: 3}, want: 6},
		{name: "default", args: args{n: 4}, want: 10},
		{name: "default", args: args{n: 5}, want: 15},
		{name: "default", args: args{n: 6}, want: 21},
		{name: "default", args: args{n: 7}, want: 28},
		{name: "default", args: args{n: 8}, want: 36},
		{name: "default", args: args{n: 9}, want: 45},
		{name: "default", args: args{n: 10}, want: 55},
		{name: "default", args: args{n: 11}, want: 66},
		{name: "default", args: args{n: 12}, want: 78},
		{name: "default", args: args{n: 13}, want: 91},
		{name: "default", args: args{n: 14}, want: 105},
		{name: "default", args: args{n: 30}, want: 465},
		{name: "default", args: args{n: 50}, want: 1275},
		{name: "default", args: args{n: 100}, want: 5050},
		{name: "default", args: args{n: 5 * 9 * 3}, want: 9180},
		{name: "default", args: args{n: 5 * 18 * 3}, want: 36585},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loopPointFormula(tt.args.n); got != tt.want {
				t.Errorf("loopPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withComma(t *testing.T) {
	type args struct {
		score string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "no comma", args: args{score: "3"}, want: "3"},
		{name: "4 digits", args: args{score: "3000"}, want: "3,000"},
		{name: "5 digits", args: args{score: "30000"}, want: "30,000"},
		{name: "6 digits", args: args{score: "300000"}, want: "300,000"},
		{name: "7 digits", args: args{score: "3000000"}, want: "3,000,000"},
		{name: "10 digits", args: args{score: "3000000000"}, want: "3,000,000,000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := withCommas(tt.args.score); got != tt.want {
				t.Errorf("withCommas() = %v, want %v", got, tt.want)
			}
		})
	}
}
