package fourth_0914

import "testing"

func TestCountMoney(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{4},
			want: 10,
		},
		{
			name: "",
			args: args{10},
			want: 37,
		},
		{
			name: "",
			args: args{20},
			want: 96,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountMoney(tt.args.n); got != tt.want {
				t.Errorf("CountMoney() = %v, want %v", got, tt.want)
			}
		})
	}
}
