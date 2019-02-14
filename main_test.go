package main

import "testing"

func TestBuildEnumFromValue(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"t1", args{s: "haha"}, "Haha"},
		{"t2", args{s: "ha ha"}, "HaHa"},
		{"t3", args{s: "ha-ha"}, "HaHa"},
		{"t4", args{s: "ha- +ha"}, "HaHa"},
		{"t5", args{s: "ha- +2ha"}, "Ha2ha"},
		{"t6", args{s: "ha- +2&ha"}, "Ha2Ha"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildEnumFromValue(tt.args.s); got != tt.want {
				t.Errorf("buildEnumFromValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
