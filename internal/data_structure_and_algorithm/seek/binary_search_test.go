package seek

import "testing"

func TestBinarySearch(t *testing.T) {
	type args struct {
		arr    []int
		target int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				arr:    []int{1, 2, 4, 5, 7},
				target: 7,
			},
			want: 4,
		},
		{
			name: "",
			args: args{
				arr:    []int{1, 2, 4, 5, 7},
				target: 9,
			},
			want: -1,
		},
		{
			name: "",
			args: args{
				arr:    []int{1, 2, 4, 5, 7},
				target: 0,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BinarySearch(tt.args.arr, tt.args.target); got != tt.want {
				t.Errorf("BinarySearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinarySearchRecursion(t *testing.T) {
	type args struct {
		arr    []int
		low    int
		high   int
		target int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				arr:    []int{1, 2, 3, 4},
				low:    0,
				high:   3,
				target: 3,
			},
			want: 2,
		},
		{
			name: "",
			args: args{
				arr:    []int{1, 2, 3, 4},
				low:    0,
				high:   3,
				target: 5,
			},
			want: -1,
		},
		{
			name: "",
			args: args{
				arr:    []int{1, 2, 3, 4},
				low:    0,
				high:   3,
				target: 0,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BinarySearchRecursion(tt.args.arr, tt.args.low, tt.args.high, tt.args.target); got != tt.want {
				t.Errorf("BinarySearchRecursion() = %v, want %v", got, tt.want)
			}
		})
	}
}
