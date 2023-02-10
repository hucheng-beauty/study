package validation

import (
    "reflect"
    "testing"
)

func TestUnique(t *testing.T) {
    type args struct {
        src []int
    }
    tests := []struct {
        name    string
        args    args
        wantDes []int
    }{
        {
            name: "Unique",
            args: args{
                []int{1, 2, 3, 3},
            },
            wantDes: []int{1, 2, 3},
        },
        {
            name: "Unique1",
            args: args{
                []int{1, 2, 3},
            },
            wantDes: []int{1, 2, 3},
        },
        {
            name: "Unique2",
            args: args{
                []int{1, 1, 1},
            },
            wantDes: []int{1},
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if gotDes := Unique(tt.args.src); !reflect.DeepEqual(gotDes, tt.wantDes) {
                t.Errorf("Unique() = %v, want %v", gotDes, tt.wantDes)
            }
        })
    }
}

func TestUniqueStr(t *testing.T) {
    type args struct {
        src []string
    }
    tests := []struct {
        name    string
        args    args
        wantDes []string
    }{
        {
            name: "UniqueStr1",
            args: args{
                []string{"c", "c++", "go", "java", "python", "python"},
            },
            wantDes: []string{"c", "c++", "go", "java", "python"},
        },
        {
            name: "UniqueStr",
            args: args{
                []string{"c", "c++", "go", "java", "python"},
            },
            wantDes: []string{"c", "c++", "go", "java", "python"},
        },
        {
            name: "UniqueStr",
            args: args{
                []string{"c", "c", "c", "c", "c"},
            },
            wantDes: []string{"c"},
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if gotDes := UniqueStr(tt.args.src); !reflect.DeepEqual(gotDes, tt.wantDes) {
                t.Errorf("UniqueStr() = %v, want %v", gotDes, tt.wantDes)
            }
        })
    }
}
