package second_0822

import (
    "reflect"
    "testing"
)

func TestQuickSort(t *testing.T) {
    type args struct {
        arr []int
    }
    tests := []struct {
        name string
        args args
        want []int
    }{
        {
            name: "QuickSort",
            args: args{
                []int{1, 4, 5, 2, 3},
            },
            want: []int{1, 2, 3, 4, 5},
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := QuickSort(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("QuickSort() = %v, want %v", got, tt.want)
            }
        })
    }
}
