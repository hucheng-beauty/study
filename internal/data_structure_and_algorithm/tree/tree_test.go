package tree

import (
	"reflect"
	"testing"
)

func TestNode_FindRecursion(t *testing.T) {
	type fields struct {
		Data  int
		Left  *Node
		Right *Node
	}
	type args struct {
		data int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Node
	}{
		{
			name:   "",
			fields: fields{},
			args:   args{data: 10},
			want:   nil,
		},
		{
			name:   "",
			fields: fields{Data: 15},
			args:   args{data: 15},
			want:   &Node{Data: 15},
		},
		{
			name:   "",
			fields: fields{Data: 12},
			args:   args{data: 15},
			want:   nil,
		},
		{
			name: "",
			fields: fields{
				Data:  10,
				Left:  &Node{Data: 6},
				Right: &Node{Data: 11},
			},
			args: args{data: 10},
			want: &Node{
				Data:  10,
				Left:  &Node{Data: 6},
				Right: &Node{Data: 11},
			},
		},
		{
			name: "",
			fields: fields{
				Data:  10,
				Left:  &Node{Data: 6},
				Right: &Node{Data: 11},
			},
			args: args{data: 6},
			want: &Node{Data: 6},
		},
		{
			name: "",
			fields: fields{
				Data:  10,
				Left:  &Node{Data: 6},
				Right: &Node{Data: 11},
			},
			args: args{data: 12},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := &Node{
				Data:  tt.fields.Data,
				Left:  tt.fields.Left,
				Right: tt.fields.Right,
			}
			if got := root.FindRecursion(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindRecursion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMainer(t *testing.T) {
	root := &Node{
		Data: 15,
		Left: &Node{
			Data: 6,
			Left: &Node{
				Data:  3,
				Left:  nil,
				Right: nil,
			},
			Right: &Node{
				Data:  10,
				Left:  nil,
				Right: nil,
			},
		},
		Right: &Node{
			Data: 20,
			Left: &Node{
				Data:  19,
				Left:  nil,
				Right: nil,
			},
			Right: &Node{
				Data:  26,
				Left:  nil,
				Right: nil,
			},
		},
	}
	t.Log("InOrder")
	root.InOrder()
	root.Delete(6)
	t.Logf("after delete data:%d\n", 6)
	root.InOrder()

	t.Log("PreOrder")
	root.PreOrder()
	root.Delete(3)
	t.Logf("after delete data:%d\n", 3)
	root.PreOrder()

	t.Log("PostOrder")
	root.PostOrder()
	root.Delete(26)
	t.Logf("after delete data:%d\n", 26)
	root.PostOrder()
}
