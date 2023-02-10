package dag_engine

import (
	"context"
	"reflect"
	"testing"
)

func TestBuildGraphFromRunners(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := BuildGraphFromRunners(D1, D2, D3, D4)
	g.Run(ctx, cancel)
	g.Run(ctx, cancel)

}

func Test_runnerProcess(t *testing.T) {
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
			wantOut: []DependAbleRunner{D1, D2, D3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := runnerProcess(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("runnerProcess() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
