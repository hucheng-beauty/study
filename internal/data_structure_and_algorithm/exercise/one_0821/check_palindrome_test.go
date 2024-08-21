package one_0821

import "testing"

func Test_checkHuiWenStr(t *testing.T) {
    type args struct {
        str string
    }
    s := "aba"
    s1 := "abca"
    s2 := "abc"
    tests := []struct {
        name string
        args args
        want bool
    }{
        {
            name: "checkPalindrome",
            args: args{s},
            want: true,
        },
        {
            name: "checkPalindrome",
            args: args{s1},
            want: true,
        },
        {
            name: "checkPalindrome",
            args: args{s2},
            want: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := checkPalindrome(tt.args.str); got != tt.want {
                t.Errorf("checkPalindrome() = %v, want %v", got, tt.want)
            }
        })
    }
}
