package sort

import (
	"reflect"
	"testing"
)

func TestCountingSort(t *testing.T) {
	type args struct {
		a []int
		n int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "CountingSort",
			args: args{a: []int{1, 4, 5, 3, 2, 0, 1, 3, 4}, n: 9},
			want: []int{0, 1, 1, 2, 3, 3, 4, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountingSort(tt.args.a, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountingSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandInt()
	}
}