package exercise

import "testing"

func TestIsPalindrome(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "IsPalindrome",
			args: args{"hello,ll e h"},
			want: true,
		},
		{
			name: "IsPalindrome",
			args: args{"hello,1ll e h"},
			want: false,
		},
		{
			name: "IsPalindrome",
			args: args{"man1 1nam"},
			want: true,
		},
		{
			name: "IsPalindrome",
			args: args{"man1 aam"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPalindrome(tt.args.s); got != tt.want {
				t.Errorf("IsPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}
