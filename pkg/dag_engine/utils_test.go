package dag_engine

import (
	"reflect"
	"testing"
)

func Test_Deduplicate(t *testing.T) {
	type args struct {
		in []DependAbleRunner
	}
	tests := []struct {
		name    string
		args    args
		wantOut []DependAbleRunner
	}{
		{
			name:    "",
			args:    args{in: []DependAbleRunner{}},
			wantOut: []DependAbleRunner{},
		},
		{
			name:    "",
			args:    args{in: []DependAbleRunner{D3}},
			wantOut: []DependAbleRunner{D3},
		},
		{
			name:    "",
			args:    args{in: []DependAbleRunner{D3, D3, D2, D1}},
			wantOut: []DependAbleRunner{D3, D2, D1},
		},
		{
			name:    "",
			args:    args{in: []DependAbleRunner{D3, D2, D1}},
			wantOut: []DependAbleRunner{D3, D2, D1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := Deduplicate(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("deduplicate() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
