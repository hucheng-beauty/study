package exercise

import "testing"

func TestReverseWords(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ReverseWords",
			args: args{"hello world"},
			want: "world hello",
		},
		{
			name: "ReverseWords",
			args: args{"xin ye ke ji."},
			want: "ji. ke ye xin",
		},
		{
			name: "ReverseWords",
			args: args{" A  B C    D "},
			want: "D C B A",
		},
		{
			name: "ReverseWords",
			args: args{"  I Love China."},
			want: "China. Love I",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseWords(tt.args.s); got != tt.want {
				t.Errorf("ReverseWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
