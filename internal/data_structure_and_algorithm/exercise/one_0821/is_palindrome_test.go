package one_0821

import (
	"testing"
)

func Test_isHuiWenString(t *testing.T) {
	type args struct {
		str string
	}
	s := "A man, a plan, a canal: Panama"
	s1 := "race a car"
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "isPalindrome",
			args: args{s},
			want: true,
		},
		{
			name: "isPalindrome",
			args: args{s1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPalindrome(tt.args.str); got != tt.want {
				t.Errorf("isPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}
