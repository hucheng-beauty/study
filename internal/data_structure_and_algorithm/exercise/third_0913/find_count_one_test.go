package third_0913

import "testing"

func Test_findCountOne(t *testing.T) {
    type args struct {
        arr []int
    }
    tests := []struct {
        name string
        args args
        want int
    }{
        {
            name: "findCountOne1",
            args: args{[]int{1, 2, 3, 2, 1}},
            want: 3,
        },
        {
            name: "findCountOne2",
            args: args{[]int{1, 2, 1}},
            want: 2,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := findCountOne(tt.args.arr); got != tt.want {
                t.Errorf("findCountOne() = %v, want %v", got, tt.want)
            }
        })
    }
}
