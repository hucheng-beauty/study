package main

import "testing"

func Test_findKth(t *testing.T) {
    type args struct {
        nums   []int
        length int
        k      int
    }
    tests := []struct {
        name string
        args args
        want int
    }{
        {
            name: "findKth0",
            args: args{
                nums:   []int{1, 3, 4, 2, 5},
                length: 5,
                k:      1}, want: 5,
        },
        {
            name: "findKth1",
            args: args{
                nums:   []int{1, 3, 4, 2, 5},
                length: 5,
                k:      2}, want: 4,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := findKth(tt.args.nums, tt.args.length, tt.args.k); got != tt.want {
                t.Errorf("findKth() = %v, want %v", got, tt.want)
            }
        })
    }
}
