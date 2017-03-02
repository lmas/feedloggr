package feedloggr2

import (
	"testing"
	"time"
)

func Test_min(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"a < b", args{0, 1}, 0},
		{"a > b", args{1, 0}, 0},
		{"a == b", args{1, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inList(t *testing.T) {
	type args struct {
		s string
		l []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"found", args{"a", []string{"a", "b", "c"}}, true},
		{"not found", args{"a", []string{"b", "c", "d"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inList(tt.args.s, tt.args.l); got != tt.want {
				t.Errorf("inList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"unix", args{time.Unix(0, 0).UTC()}, "1970-01-01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := date(tt.args.t); got != tt.want {
				t.Errorf("date() = %v, want %v", got, tt.want)
			}
		})
	}
}
